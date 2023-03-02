package migration

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"text/template"
	"time"

	"github.com/gotsunami/gitlab-copy/config"
	"github.com/gotsunami/gitlab-copy/gitlab"
	"github.com/rotisserie/eris"
	glab "github.com/xanzy/go-gitlab"
)

var (
	errDuplicateIssue = errors.New("Duplicate Issue")
)

const (
	// ResultsPerPage is the Number of results per page.
	ResultsPerPage = 100
)

// Endpoint refers to the GitLab server endpoints.
type Endpoint struct {
	SrcClient, DstClient gitlab.GitLaber
}

// Migration defines a migration step.
type Migration struct {
	params                 *config.Config
	Endpoint               *Endpoint
	srcProject, dstProject *glab.Project
	toUsers                map[string]gitlab.GitLaber
	skipIssue              bool
}

// New creates a new migration.
func New(c *config.Config) (*Migration, error) {
	if c == nil {
		return nil, errors.New("nil params")
	}
	m := &Migration{params: c}
	m.toUsers = make(map[string]gitlab.GitLaber)

	fromgl, err := gitlab.Service().WithToken(
		c.SrcPrj.Token,
		glab.WithBaseURL(c.SrcPrj.ServerURL),
	)
	if err != nil {
		return nil, eris.Wrap(err, "migration: src token")
	}
	togl, err := gitlab.Service().WithToken(
		c.DstPrj.Token,
		glab.WithBaseURL(c.DstPrj.ServerURL),
	)
	if err != nil {
		return nil, eris.Wrap(err, "migration: dst token")
	}
	for user, token := range c.DstPrj.Users {
		uc, err := gitlab.Service().WithToken(token, glab.WithBaseURL(c.DstPrj.ServerURL))
		if err != nil {
			return nil, eris.Wrap(err, "migration: dst users check")
		}
		m.toUsers[user] = uc
	}
	m.Endpoint = &Endpoint{fromgl, togl}
	return m, nil
}

// Returns project by name.
func (m *Migration) project(endpoint gitlab.GitLaber, name, which string) (*glab.Project, error) {
	proj, resp, err := endpoint.GetProject(name, nil)
	if resp == nil {
		return nil, errors.New("network error while fetching project info: nil response")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s project '%s' not found", which, name)
	}
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (m *Migration) SourceProject(name string) (*glab.Project, error) {
	p, err := m.project(m.Endpoint.SrcClient, name, "source")
	if err != nil {
		return nil, err
	}
	m.srcProject = p
	return p, nil
}

func (m *Migration) DestProject(name string) (*glab.Project, error) {
	p, err := m.project(m.Endpoint.SrcClient, name, "target")
	if err != nil {
		return nil, err
	}
	m.dstProject = p
	return p, nil
}

func (m *Migration) migrateIssue(issueID int) error {
	source := m.Endpoint.SrcClient
	target := m.Endpoint.DstClient

	srcProjectID := m.srcProject.ID
	tarProjectID := m.dstProject.ID

	issue, _, err := source.GetIssue(srcProjectID, issueID)
	if err != nil {
		return fmt.Errorf("target: can't fetch issue: %s", err.Error())
	}
	tis, _, err := target.ListProjectIssues(tarProjectID, nil)
	if err != nil {
		return fmt.Errorf("target: can't fetch issue: %s", err.Error())
	}
	for _, t := range tis {
		if issue.Title == t.Title {
			// Target issue already exists, let's skip this one.
			return errDuplicateIssue
		}
	}
	labels := make(glab.Labels, 0)
	iopts := &glab.CreateIssueOptions{
		Title:       &issue.Title,
		Description: &issue.Description,
		Labels:      &labels,
	}
	if issue.Assignee.Username != "" {
		// Assigned, does target user exist?
		// User may have a different ID on target
		users, _, err := target.ListUsers(nil)
		if err == nil {
			for _, u := range users {
				if u.Username == issue.Assignee.Username {
					iopts.AssigneeIDs = &[]int{u.ID}
					break
				}
			}
		} else {
			return fmt.Errorf("target: error fetching users: %v", err)
		}
	}
	if issue.Milestone != nil && issue.Milestone.Title != "" {
		miles, _, err := target.ListMilestones(tarProjectID, nil)
		if err == nil {
			found := false
			for _, mi := range miles {
				found = false
				if mi.Title == issue.Milestone.Title {
					found = true
					iopts.MilestoneID = &mi.ID
					break
				}
			}
			if !found {
				// Create target milestone
				cmopts := &glab.CreateMilestoneOptions{
					Title:       &issue.Milestone.Title,
					Description: &issue.Milestone.Description,
					DueDate:     issue.Milestone.DueDate,
				}
				mi, _, err := target.CreateMilestone(tarProjectID, cmopts)
				if err == nil {
					iopts.MilestoneID = &mi.ID
				} else {
					return fmt.Errorf("target: error creating milestone '%s': %s", issue.Milestone.Title, err.Error())
				}
			}
		} else {
			return fmt.Errorf("target: error listing milestones: %s", err.Error())
		}
	}
	// Copy existing labels.
	for _, label := range issue.Labels {
		*iopts.Labels = append(*iopts.Labels, label)
	}
	// Create target issue if not existing (same name).
	ni, resp, err := target.CreateIssue(tarProjectID, iopts)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusRequestURITooLong {
			fmt.Printf("target: caught a %q error, shortening issue's decription length ...\n", http.StatusText(resp.StatusCode))
			if len(*iopts.Description) == 0 {
				return fmt.Errorf("target: error creating issue: no description but %q error", http.StatusText(resp.StatusCode))
			}
			smalld := (*iopts.Description)[:1024]
			iopts.Description = &smalld
			ni, _, err = target.CreateIssue(tarProjectID, iopts)
			if err != nil {
				return fmt.Errorf("target: error creating empty issue: %s", err.Error())
			}
		} else {
			return fmt.Errorf("target: error creating issue: %s", err.Error())
		}
	}

	// Copy related notes (comments)
	notes, _, err := source.ListIssueNotes(srcProjectID, issue.IID, nil)
	if err != nil {
		return fmt.Errorf("source: can't get issue #%d notes: %s", issue.IID, err.Error())
	}
	opts := &glab.CreateIssueNoteOptions{}
	// Notes on target will be added in reverse order.
	for j := len(notes) - 1; j >= 0; j-- {
		n := notes[j]
		target = m.Endpoint.DstClient
		// Can we write the comment with user ownership?
		if _, ok := m.toUsers[n.Author.Username]; ok {
			target = m.toUsers[n.Author.Username]
			opts.Body = &n.Body
		} else {
			// Nope. Let's add a header note instead.
			head := fmt.Sprintf("%s @%s wrote on %s :", n.Author.Name, n.Author.Username, n.CreatedAt.Format(time.RFC1123))
			bd := fmt.Sprintf("%s\n\n%s", head, n.Body)
			opts.Body = &bd
		}
		_, resp, err := target.CreateIssueNote(tarProjectID, ni.IID, opts)
		if err != nil {
			if resp.StatusCode == http.StatusRequestURITooLong {
				fmt.Printf("target: note's body too long, shortening it ...\n")
				if len(*opts.Body) > 1024 {
					smallb := (*opts.Body)[:1024]
					opts.Body = &smallb
				}
				_, _, err := target.CreateIssueNote(tarProjectID, ni.ID, opts)
				if err != nil {
					return fmt.Errorf("target: error creating note (with shorter body) for issue #%d: %s", ni.IID, err.Error())
				}
			} else {
				return fmt.Errorf("target: error creating note for issue #%d: %s", ni.IID, err.Error())
			}
		}
	}
	target = m.Endpoint.DstClient

	if issue.State == "closed" {
		event := "close"
		_, _, err := target.UpdateIssue(tarProjectID, ni.IID,
			&glab.UpdateIssueOptions{StateEvent: &event, Labels: &issue.Labels})
		if err != nil {
			return fmt.Errorf("target: error closing issue #%d: %s", ni.IID, err.Error())
		}
	}
	// Add a link to target issue if needed
	if m.params.SrcPrj.LinkToTargetIssue {
		var dstProjectURL string
		// Strip URL if moving on the same GitLab  installation.
		if m.Endpoint.SrcClient.BaseURL().Host == m.Endpoint.DstClient.BaseURL().Host {
			dstProjectURL = m.dstProject.PathWithNamespace
		} else {
			dstProjectURL = m.dstProject.WebURL
		}
		tmpl, err := template.New("link").Parse(m.params.SrcPrj.LinkToTargetIssueText)
		if err != nil {
			return fmt.Errorf("link to target issue: error parsing linkToTargetIssueText parameter: %s", err.Error())
		}
		noteLink := fmt.Sprintf("%s#%d", dstProjectURL, ni.IID)
		type link struct {
			Link string
		}
		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, &link{noteLink}); err != nil {
			return fmt.Errorf("link to target issue: %s", err.Error())
		}
		nopt := buf.String()
		opts := &glab.CreateIssueNoteOptions{
			Body: &nopt,
		}
		_, _, err = target.CreateIssueNote(srcProjectID, issue.IID, opts)
		if err != nil {
			return fmt.Errorf("source: error adding closing note for issue #%d: %s", issue.IID, err.Error())
		}
	}
	// Auto close source issue if needed
	if m.params.SrcPrj.AutoCloseIssues {
		event := "close"
		_, _, err := source.UpdateIssue(srcProjectID, issue.ID,
			&glab.UpdateIssueOptions{StateEvent: &event, Labels: &issue.Labels})
		if err != nil {
			return fmt.Errorf("source: error closing issue #%d: %s", issue.IID, err.Error())
		}
	}

	fmt.Printf("target: created issue #%d: %s [%s]\n", ni.IID, ni.Title, issue.State)
	return nil
}

type issueID struct {
	IID, ID int
}

type byIID []issueID

func (a byIID) Len() int           { return len(a) }
func (a byIID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIID) Less(i, j int) bool { return a[i].IID < a[j].IID }

// Migrate performs the issues migration.
func (m *Migration) Migrate() error {
	_, err := m.SourceProject(m.params.SrcPrj.Name)
	if err != nil {
		return eris.Wrap(err, "migrate")
	}
	_, err = m.DestProject(m.params.DstPrj.Name)
	if err != nil {
		return eris.Wrap(err, "migrate")
	}

	source := m.Endpoint.SrcClient
	target := m.Endpoint.DstClient

	srcProjectID := m.srcProject.ID
	tarProjectID := m.dstProject.ID

	curPage := 1
	optSort := "asc"
	opts := &glab.ListProjectIssuesOptions{Sort: &optSort, ListOptions: glab.ListOptions{PerPage: ResultsPerPage, Page: curPage}}

	s := make([]issueID, 0)

	// Copy all source labels on target
	labels, _, err := source.ListLabels(srcProjectID, nil)
	if err != nil {
		return fmt.Errorf("source: can't fetch labels: %s", err.Error())
	}
	fmt.Printf("Found %d labels ...\n", len(labels))
	for _, label := range labels {
		clopts := &glab.CreateLabelOptions{Name: &label.Name, Color: &label.Color, Description: &label.Description}
		_, resp, err := target.CreateLabel(tarProjectID, clopts)
		if err != nil {
			// GitLab returns a 409 code if label already exists
			if resp.StatusCode != http.StatusConflict {
				return fmt.Errorf("target: error creating label '%s': %s", label, err.Error())
			}
		}
	}

	if m.params.SrcPrj.LabelsOnly {
		// We're done here
		return nil
	}

	if m.params.SrcPrj.MilestonesOnly {
		fmt.Println("Copying milestones ...")
		miles, _, err := source.ListMilestones(srcProjectID, nil)
		if err != nil {
			return fmt.Errorf("error getting the milestones from source project: %s", err.Error())
		}
		fmt.Printf("Found %d milestones\n", len(miles))
		for _, mi := range miles {
			// Create target milestone
			cmopts := &glab.CreateMilestoneOptions{
				Title:       &mi.Title,
				Description: &mi.Description,
				DueDate:     mi.DueDate,
			}
			tmi, _, err := target.CreateMilestone(tarProjectID, cmopts)
			if err != nil {
				return fmt.Errorf("target: error creating milestone '%s': %s", mi.Title, err.Error())
			}
			if mi.State == "closed" {
				event := "close"
				umopts := &glab.UpdateMilestoneOptions{
					StateEvent: &event,
				}
				_, _, err := target.UpdateMilestone(tarProjectID, tmi.ID, umopts)
				if err != nil {
					return fmt.Errorf("target: error closing milestone '%s': %s", mi.Title, err.Error())
				}
			}
		}
		// We're done here
		return nil
	}

	fmt.Println("Copying issues ...")

	// First, count issues
	for {
		issues, _, err := source.ListProjectIssues(srcProjectID, opts)
		if err != nil {
			return err
		}
		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {
			s = append(s, issueID{IID: issue.IID, ID: issue.ID})
		}
		curPage++
		opts.Page = curPage
	}

	// Then sort
	sort.Sort(byIID(s))

	for _, issue := range s {
		if m.params.SrcPrj.Matches(issue.IID) {
			if err := m.migrateIssue(issue.IID); err != nil {
				if err == errDuplicateIssue {
					fmt.Printf("target: issue %d already exists, skipping...", issue.IID)
					continue
				}
				return err
			}
			if m.params.SrcPrj.MoveIssues {
				// Delete issue from source project
				_, err := source.DeleteIssue(srcProjectID, issue.ID)
				if err != nil {
					log.Printf("could not delete the issue %d: %s", issue.ID, err.Error())
				}
			}
		}
	}

	return nil
}
