package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

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
func (m *migration) project(endpoint *gitlab.Client, name string) (*gitlab.Project, error) {
	proj, resp, err := endpoint.Projects.GetProject(name)
	if resp == nil {
		return nil, errors.New("network error: " + err.Error())
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("source project '%s' not found", name)
	}
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (m *migration) sourceProject(name string) (*gitlab.Project, error) {
	p, err := m.project(m.endpoint.from, name)
	if err != nil {
		return nil, err
	}
	m.srcProject = p
	return p, nil
}

func (m *migration) destProject(name string) (*gitlab.Project, error) {
	p, err := m.project(m.endpoint.to, name)
	if err != nil {
		return nil, err
	}
	m.dstProject = p
	return p, nil
}

// Performs the issues migration.
func (m *migration) migrate() error {
	if m.srcProject == nil || m.dstProject == nil {
		return errors.New("nil project.")
	}
	fmt.Println("Migrating ...")

	source := m.endpoint.from
	target := m.endpoint.to

	srcProjectID := *m.srcProject.ID
	tarProjectID := *m.dstProject.ID

	curPage := 1
	opts := &gitlab.ListProjectIssuesOptions{ListOptions: gitlab.ListOptions{PerPage: resultsPerPage, Page: curPage}}
	issues, _, err := source.Issues.ListProjectIssues(srcProjectID, opts)
	if err != nil {
		return err
	}
	if len(issues) > 0 {
		skipIssue := false
		for _, issue := range issues {
			tis, _, err := target.Issues.ListProjectIssues(tarProjectID, nil)
			if err != nil {
				log.Printf("target: can't fetch issues, skipping...")
				continue
			}
			skipIssue = false
			for _, t := range tis {
				if issue.Title == t.Title {
					// Target issue already exists, let's skip this one
					skipIssue = true
					log.Printf("target: issue '%s' already exists, skipping...", issue.Title)
					break
				}
			}
			if skipIssue {
				continue
			}
			iopts := &gitlab.CreateIssueOptions{
				Title:       issue.Title,
				Description: issue.Description,
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
					log.Printf("target: error fetching users: %s", err.Error())
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
							log.Printf("target: error creating milestone '%s': %s", issue.Milestone.Title, err.Error())
						}
					}
				}
			}
			if len(issue.Labels) > 0 {
				lbls, _, err := target.Labels.ListLabels(tarProjectID)
				targetLabels := make([]string, 0)
				if err == nil {
					found := false
					for _, label := range issue.Labels {
						found = false
						for _, l := range lbls {
							if l.Name == label {
								found = true
								break
							}
						}
						if !found {
							// Create target label
							// FIXME: label color
							clopts := &gitlab.CreateLabelOptions{Name: label, Color: "#329557"}
							_, _, err := target.Labels.CreateLabel(tarProjectID, clopts)
							if err == nil {
								targetLabels = append(targetLabels, label)
							} else {
								log.Printf("target: error creating label '%s': %s", label, err.Error())
							}
						} else {
							targetLabels = append(targetLabels, label)
						}
					}
				} else {
					log.Printf("target: error fetching labels: %s", err.Error())
				}
				iopts.Labels = targetLabels
			}
			// Create target issue if not existing (same name)
			ni, _, err := target.Issues.CreateIssue(tarProjectID, iopts)
			if err != nil {
				log.Printf("target: error creating issue: %s", err.Error())
			}
			if issue.State == "closed" {
				_, _, err := target.Issues.UpdateIssue(tarProjectID, ni.ID, &gitlab.UpdateIssueOptions{StateEvent: "close"})
				if err != nil {
					log.Printf("target: error closing issue #%d: %s", ni.IID, err.Error())
				}
			}
			fmt.Printf("target: created issue #%d: %s [%s]\n", ni.IID, ni.Title, issue.State)

			// Copy related notes (comments)
			notes, _, err := source.Notes.ListIssueNotes(srcProjectID, issue.ID, nil)
			if err != nil {
				log.Printf("source: can't get issue #%d notes: %s", issue.ID, err.Error())
			}
			opts := &gitlab.CreateIssueNoteOptions{}
			ns, _, err := target.Notes.ListIssueNotes(tarProjectID, ni.ID, nil)
			if err != nil {
				log.Printf("target: can't get issue #%d notes: %s", ni.ID, err.Error())
			}
			for _, n := range notes {
				if len(ns) > 0 {
					for _, m := range ns {
						if m.Body != n.Body {
							opts.Body = m.Body
							_, _, err := target.Notes.CreateIssueNote(tarProjectID, ni.ID, opts)
							if err != nil {
								log.Printf("target: error creating note for issue #%d: %s", ni.IID, err.Error())
							}
						}
					}
				} else {
					opts.Body = n.Body
					_, _, err := target.Notes.CreateIssueNote(tarProjectID, ni.ID, opts)
					if err != nil {
						log.Printf("target: error creating note for issue #%d: %s", ni.IID, err.Error())
					}
				}
			}
		}
	}
	return nil
}
