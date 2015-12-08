package main

import (
	"errors"
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type projectStats struct {
	project                               *gitlab.Project
	nbIssues, nbClosed, nbOpened, nbNotes int
	milestones, labels                    map[string]int
}

func newProjectStats(prj *gitlab.Project) *projectStats {
	p := new(projectStats)
	p.project = prj
	p.milestones = make(map[string]int)
	p.labels = make(map[string]int)
	return p
}

func (p *projectStats) String() string {
	return fmt.Sprintf("%d issues (%d opened, %d closed)", p.nbIssues, p.nbOpened, p.nbClosed)
}

func (p *projectStats) pagination(client *gitlab.Client, f func(*gitlab.Client, *gitlab.ListOptions) (bool, error)) error {
	if client == nil {
		return errors.New("nil client")
	}

	curPage := 1
	opts := &gitlab.ListOptions{PerPage: resultsPerPage, Page: curPage}

	for {
		stop, err := f(client, opts)
		if err != nil {
			return err
		}
		if stop {
			break
		}
		curPage++
		opts.Page = curPage
	}
	return nil
}

func (p *projectStats) computeStats(client *gitlab.Client) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c *gitlab.Client, lo *gitlab.ListOptions) (bool, error) {
		opts := &gitlab.ListProjectIssuesOptions{ListOptions: gitlab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.Issues.ListProjectIssues(*p.project.ID, opts)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			p.nbIssues += len(issues)
			for _, issue := range issues {
				switch issue.State {
				case "opened":
					p.nbOpened++
				case "closed":
					p.nbClosed++
				}
				if issue.Milestone.Title != "" {
					p.milestones[issue.Milestone.Title]++
				}
				if len(issue.Labels) > 0 {
					for _, label := range issue.Labels {
						p.labels[label]++
					}
				}
			}
		} else {
			// Exit
			return true, nil
		}
		return false, nil
	}

	if err := p.pagination(client, action); err != nil {
		return err
	}
	return nil
}

func (p *projectStats) computeIssueNotes(client *gitlab.Client) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c *gitlab.Client, lo *gitlab.ListOptions) (bool, error) {
		opts := &gitlab.ListProjectIssuesOptions{ListOptions: gitlab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.Issues.ListProjectIssues(*p.project.ID, opts)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			for _, issue := range issues {
				notes, _, err := client.Notes.ListIssueNotes(*p.project.ID, issue.ID, nil)
				if err != nil {
					return false, err
				}
				p.nbNotes += len(notes)
			}
		} else {
			// Exit
			return true, nil
		}
		return false, nil
	}

	if err := p.pagination(client, action); err != nil {
		return err
	}
	return nil
}
