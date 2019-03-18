package stats

import (
	"errors"
	"fmt"

	"github.com/gotsunami/gitlab-copy/gitlab"
	"github.com/gotsunami/gitlab-copy/migration"
	glab "github.com/xanzy/go-gitlab"
)

const (
	LabelsPerPage = 100
)

type ProjectStats struct {
	Project                               *glab.Project
	NbIssues, NbClosed, NbOpened, NbNotes int
	Milestones, Labels                    map[string]int
}

func NewProject(prj *glab.Project) *ProjectStats {
	p := new(ProjectStats)
	p.Project = prj
	p.Milestones = make(map[string]int)
	p.Labels = make(map[string]int)
	return p
}

func (p *ProjectStats) String() string {
	return fmt.Sprintf("%d issues (%d opened, %d closed)", p.NbIssues, p.NbOpened, p.NbClosed)
}

func (p *ProjectStats) pagination(client gitlab.GitLaber, f func(gitlab.GitLaber, *glab.ListOptions) (bool, error)) error {
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

func (p *ProjectStats) ComputeStats(client gitlab.GitLaber) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c gitlab.GitLaber, lo *glab.ListOptions) (bool, error) {
		opts := &glab.ListProjectIssuesOptions{ListOptions: glab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.ListProjectIssues(p.Project.ID, opts)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			p.NbIssues += len(issues)
			for _, issue := range issues {
				switch issue.State {
				case "opened":
					p.NbOpened++
				case "closed":
					p.NbClosed++
				}
				if issue.Milestone != nil && issue.Milestone.Title != "" {
					p.Milestones[issue.Milestone.Title]++
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

	currentLabelPage := 1
	for {
		paginatedLabels, _, err := client.ListLabels(p.Project.ID, &glab.ListLabelsOptions{PerPage: LabelsPerPage, Page: currentLabelPage})
		if err != nil {
			return fmt.Errorf("source: can't fetch labels: %s", err.Error())
		}
		if len(paginatedLabels) == 0 {
			break
		}

		for _, label := range paginatedLabels {
			p.Labels[label.Name]++
		}
		currentLabelPage++
	}

	return nil
}

func (p *ProjectStats) ComputeIssueNotes(client gitlab.GitLaber) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c gitlab.GitLaber, lo *glab.ListOptions) (bool, error) {
		opts := &glab.ListProjectIssuesOptions{ListOptions: glab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.ListProjectIssues(p.Project.ID, opts)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			for _, issue := range issues {
				notes, _, err := client.ListIssueNotes(p.Project.ID, issue.IID, nil)
				if err != nil {
					return false, err
				}
				p.NbNotes += len(notes)
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
