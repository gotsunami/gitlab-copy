package main

import (
	"errors"
	"fmt"

	"github.com/gotsunami/gitlab-copy/gitlab"
	"github.com/gotsunami/gitlab-copy/migration"
	glab "github.com/xanzy/go-gitlab"
)

type projectStats struct {
	project                               *glab.Project
	nbIssues, nbClosed, nbOpened, nbNotes int
	milestones, labels                    map[string]int
}

func newProjectStats(prj *glab.Project) *projectStats {
	p := new(projectStats)
	p.project = prj
	p.milestones = make(map[string]int)
	p.labels = make(map[string]int)
	return p
}

func (p *projectStats) String() string {
	return fmt.Sprintf("%d issues (%d opened, %d closed)", p.nbIssues, p.nbOpened, p.nbClosed)
}

func (p *projectStats) pagination(client gitlab.GitLaber, f func(gitlab.GitLaber, *glab.ListOptions) (bool, error)) error {
	if client == nil {
		return errors.New("nil client")
	}

	curPage := 1
	opts := &glab.ListOptions{PerPage: migration.ResultsPerPage, Page: curPage}

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

func (p *projectStats) computeStats(client gitlab.GitLaber) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c gitlab.GitLaber, lo *glab.ListOptions) (bool, error) {
		opts := &glab.ListProjectIssuesOptions{ListOptions: glab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.ListProjectIssues(p.project.ID, opts)
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
				if issue.Milestone != nil && issue.Milestone.Title != "" {
					p.milestones[issue.Milestone.Title]++
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

	labels, _, err := client.ListLabels(p.project.ID, nil)
	if err != nil {
		return fmt.Errorf("source: can't fetch labels: %s", err.Error())
	}
	for _, label := range labels {
		p.labels[label.Name]++
	}
	return nil
}

func (p *projectStats) computeIssueNotes(client gitlab.GitLaber) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c gitlab.GitLaber, lo *glab.ListOptions) (bool, error) {
		opts := &glab.ListProjectIssuesOptions{ListOptions: glab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.ListProjectIssues(p.project.ID, opts)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			for _, issue := range issues {
				notes, _, err := client.ListIssueNotes(p.project.ID, issue.IID, nil)
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
