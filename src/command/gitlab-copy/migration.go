package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/xanzy/go-gitlab"
)

const (
	resultsPerPage = 100
)

// GitLab server endpoints
type endpoint struct {
	from, to *gitlab.Client
}

type migration struct {
	params                 *config
	endpoint               *endpoint
	srcProject, dstProject *gitlab.Project
}

func NewMigration(c *config) (*migration, error) {
	if c == nil {
		return nil, errors.New("nil params")
	}
	fromgl := gitlab.NewClient(nil, c.From.Token)
	if err := fromgl.SetBaseURL(c.From.ServerURL); err != nil {
		return nil, err
	}
	togl := gitlab.NewClient(nil, c.To.Token)
	if err := togl.SetBaseURL(c.To.ServerURL); err != nil {
		return nil, err
	}
	m := &migration{params: c, endpoint: &endpoint{fromgl, togl}}
	return m, nil
}

// Returns project by name.
func (m *migration) project(endpoint *gitlab.Client, name, which string) (*gitlab.Project, error) {
	proj, resp, err := endpoint.Projects.GetProject(name)
	if resp == nil {
		return nil, errors.New("network error: " + err.Error())
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s project '%s' not found", which, name)
	}
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (m *migration) sourceProject(name string) (*gitlab.Project, error) {
	p, err := m.project(m.endpoint.from, name, "source")
	if err != nil {
		return nil, err
	}
	m.srcProject = p
	return p, nil
}

func (m *migration) destProject(name string) (*gitlab.Project, error) {
	p, err := m.project(m.endpoint.to, name, "target")
	if err != nil {
		return nil, err
	}
	m.dstProject = p
	return p, nil
}

func (m *migration) migrateIssue(issueID int) error {
	source := m.endpoint.from
	target := m.endpoint.to

	srcProjectID := *m.srcProject.ID
	tarProjectID := *m.dstProject.ID

	issue, _, err := source.Issues.GetIssue(srcProjectID, issueID)
	if err != nil {
		return fmt.Errorf("target: can't fetch issue: %s", err.Error())
	}
	tis, _, err := target.Issues.ListProjectIssues(tarProjectID, nil)
	if err != nil {
		return fmt.Errorf("target: can't fetch issue: %s", err.Error())
	}
	skipIssue := false
	for _, t := range tis {
		if issue.Title == t.Title {
			// Target issue already exists, let's skip this one
			skipIssue = true
			log.Printf("target: issue '%s' already exists, skipping...", issue.Title)
			break
		}
	}
	if skipIssue {
		return nil
	}
	iopts := &gitlab.CreateIssueOptions{
		Title:       issue.Title,
		Description: issue.Description,
		Labels:      make([]string, 0),
	}
	if issue.Assignee.Username != "" {
		// Assigned, does target user exist?
		// User may have a different ID on target
		users, _, err := target.Users.ListUsers(nil)
		if err == nil {
			for _, u := range users {
				if u.Username == issue.Assignee.Username {
					iopts.AssigneeID = u.ID
					break
				}
			}
		} else {
			return fmt.Errorf("target: error fetching users: %s", err.Error())
		}
	}
	if issue.Milestone.Title != "" {
		miles, _, err := target.Milestones.ListMilestones(tarProjectID, nil)
		if err == nil {
			found := false
			for _, mi := range miles {
				found = false
				if mi.Title == issue.Milestone.Title {
					found = true
					iopts.MilestoneID = mi.ID
					break
				}
			}
			if !found {
				// Create target milestone
				cmopts := &gitlab.CreateMilestoneOptions{
					Title:       issue.Milestone.Title,
					Description: issue.Milestone.Description,
					DueDate:     issue.Milestone.DueDate,
				}
				mi, _, err := target.Milestones.CreateMilestone(tarProjectID, cmopts)
				if err == nil {
					iopts.MilestoneID = mi.ID
				} else {
					return fmt.Errorf("target: error creating milestone '%s': %s", issue.Milestone.Title, err.Error())
				}
			}
		}
	}
	for _, label := range issue.Labels {
		iopts.Labels = append(iopts.Labels, label)
	}
	// Create target issue if not existing (same name)
	ni, resp, err := target.Issues.CreateIssue(tarProjectID, iopts)
	if err != nil {
		if resp.StatusCode == http.StatusRequestURITooLong {
			fmt.Printf("target: catched a \"%s\" error, shortening issue's decription length ...\n", http.StatusText(resp.StatusCode))
			iopts.Description = iopts.Description[:1024]
			ni, _, err = target.Issues.CreateIssue(tarProjectID, iopts)
			if err != nil {
				return fmt.Errorf("target: error creating empty issue: %s", err.Error())
			}
		} else {
			return fmt.Errorf("target: error creating issue: %s", err.Error())
		}
	}

	// Copy related notes (comments)
	notes, _, err := source.Notes.ListIssueNotes(srcProjectID, issue.ID, nil)
	if err != nil {
		return fmt.Errorf("source: can't get issue #%d notes: %s", issue.ID, err.Error())
	}
	opts := &gitlab.CreateIssueNoteOptions{}
	for _, n := range notes {
		head := fmt.Sprintf("%s @%s wrote on %s :", n.Author.Name, n.Author.Username, n.CreatedAt.Format(time.RFC1123))
		opts.Body = fmt.Sprintf("%s\n\n%s", head, n.Body)
		_, resp, err := target.Notes.CreateIssueNote(tarProjectID, ni.ID, opts)
		if err != nil {
			if resp.StatusCode == http.StatusRequestURITooLong {
				fmt.Printf("target: note's body too long, shortening it ...\n")
				opts.Body = fmt.Sprintf("%s\n\n%s", head, n.Body[:1024])
				_, _, err := target.Notes.CreateIssueNote(tarProjectID, ni.ID, opts)
				if err != nil {
					return fmt.Errorf("target: error creating note (with shorter body) for issue #%d: %s", ni.IID, err.Error())
				}
			} else {
				return fmt.Errorf("target: error creating note for issue #%d: %s", ni.IID, err.Error())
			}
		}
	}

	if issue.State == "closed" {
		_, _, err := target.Issues.UpdateIssue(tarProjectID, ni.ID, &gitlab.UpdateIssueOptions{StateEvent: "close", Labels: issue.Labels})
		if err != nil {
			return fmt.Errorf("target: error closing issue #%d: %s", ni.IID, err.Error())
		}
	}
	fmt.Printf("target: created issue #%d: %s [%s]\n", ni.IID, ni.Title, issue.State)

	return nil
}

type issueId struct {
	IID, ID int
}

type byIID []issueId

func (a byIID) Len() int           { return len(a) }
func (a byIID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIID) Less(i, j int) bool { return a[i].IID < a[j].IID }

// Performs the issues migration.
func (m *migration) migrate() error {
	if m.srcProject == nil || m.dstProject == nil {
		return errors.New("nil project.")
	}
	fmt.Println("Copying labels ...")

	source := m.endpoint.from
	target := m.endpoint.to

	srcProjectID := *m.srcProject.ID
	tarProjectID := *m.dstProject.ID

	curPage := 1
	opts := &gitlab.ListProjectIssuesOptions{Sort: "asc", ListOptions: gitlab.ListOptions{PerPage: resultsPerPage, Page: curPage}}

	s := make([]issueId, 0)

	// Copy all source labels on target
	labels, _, err := source.Labels.ListLabels(srcProjectID)
	if err != nil {
		return fmt.Errorf("source: can't fetch labels: %s", err.Error())
	}
	for _, label := range labels {
		clopts := &gitlab.CreateLabelOptions{Name: label.Name, Color: label.Color}
		_, resp, err := target.Labels.CreateLabel(tarProjectID, clopts)
		if err != nil {
			// GitLab returns a 409 code if label already exists
			if resp.StatusCode != http.StatusConflict {
				return fmt.Errorf("target: error creating label '%s': %s", label, err.Error())
			}
		}
	}

	if m.params.From.LabelsOnly {
		// We're done here
		return nil
	}

	fmt.Println("Copying issues ...")

	// First, count issues
	for {
		issues, _, err := source.Issues.ListProjectIssues(srcProjectID, opts)
		if err != nil {
			return err
		}
		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {
			s = append(s, issueId{IID: issue.IID, ID: issue.ID})
		}
		curPage++
		opts.Page = curPage
	}

	// Then sort
	sort.Sort(byIID(s))

	for _, issue := range s {
		if m.params.From.matches(issue.IID) {
			if err := m.migrateIssue(issue.ID); err != nil {
				log.Printf(err.Error())
			}
			if m.params.From.MoveIssues {
				// Delete issue from source project
				_, err := source.Issues.DeleteIssue(srcProjectID, issue.ID)
				if err != nil {
					log.Printf("could not delete the issue %d: %s", issue.ID, err.Error())
				}
			}
		}
	}

	return nil
}
