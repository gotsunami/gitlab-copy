package migration

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gotsunami/gitlab-copy/gitlab"
	glab "github.com/xanzy/go-gitlab"
)

type fakeClient struct {
	baseUrl *url.URL
	errors  struct {
		createIssue, createIssueNote, createLabel, createMilestone   error
		deleteIssue                                                  error
		getIssue, getProject                                         error
		listIssueNotes, listLabels, listMilestones, listProjetIssues error
		listUsers                                                    error
		updateIssue, updateMilestone                                 error
		baseURL                                                      error
	}
	labels         []*glab.Label
	milestones     []*glab.Milestone
	users          []*glab.User
	issues         []*glab.Issue
	exitPagination bool
}

// New fake GitLab client, for the UT.
func (c *fakeClient) New(httpClient *http.Client, token string) gitlab.GitLaber {
	return new(fakeClient)
}

func (c *fakeClient) SetBaseURL(u string) error {
	err := c.errors.baseURL
	if err != nil {
		return err
	}
	uu, err := url.Parse(u)
	c.baseUrl = uu
	return err
}

func (c *fakeClient) BaseURL() *url.URL {
	return c.baseUrl
}

func (c *fakeClient) GetProject(interface{}, ...glab.OptionFunc) (*glab.Project, *glab.Response, error) {
	err := c.errors.getProject
	if err != nil {
		return nil, nil, err
	}
	p := new(glab.Project)
	p.Name = "A name"
	r := &glab.Response{
		Response: new(http.Response),
	}
	r.StatusCode = http.StatusOK
	return p, r, nil
}

func (c *fakeClient) CreateLabel(id interface{}, opt *glab.CreateLabelOptions, options ...glab.OptionFunc) (*glab.Label, *glab.Response, error) {
	r := &glab.Response{
		Response: new(http.Response),
	}
	err := c.errors.createLabel
	if err != nil {
		r.StatusCode = http.StatusBadRequest
		return nil, r, err
	}
	r.StatusCode = http.StatusOK
	for _, l := range c.labels {
		if l.Name == *opt.Name {
			return nil, nil, fmt.Errorf("label %q already exists", l.Name)
		}
	}
	l := &glab.Label{
		Name:        *opt.Name,
		Color:       *opt.Color,
		Description: *opt.Description,
	}
	c.labels = append(c.labels, l)
	return l, r, nil
}

func (c *fakeClient) ListLabels(id interface{}, opt *glab.ListLabelsOptions, options ...glab.OptionFunc) ([]*glab.Label, *glab.Response, error) {
	err := c.errors.listLabels
	if err != nil {
		return nil, nil, err
	}
	return c.labels, nil, nil
}

func (c *fakeClient) ListMilestones(id interface{}, opt *glab.ListMilestonesOptions, options ...glab.OptionFunc) ([]*glab.Milestone, *glab.Response, error) {
	err := c.errors.listMilestones
	if err != nil {
		return nil, nil, err
	}
	return c.milestones, nil, nil
}

func (c *fakeClient) CreateMilestone(id interface{}, opt *glab.CreateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	err := c.errors.createMilestone
	if err != nil {
		return nil, nil, err
	}
	m := &glab.Milestone{
		ID:    len(c.milestones),
		Title: *opt.Title,
	}
	for _, p := range c.milestones {
		if p.Title == m.Title {
			return nil, nil, fmt.Errorf("milestone %q already exists", p.Title)
		}
	}
	c.milestones = append(c.milestones, m)
	return m, nil, nil
}

func (c *fakeClient) UpdateMilestone(id interface{}, milestone int, opt *glab.UpdateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	err := c.errors.updateMilestone
	if err != nil {
		return nil, nil, err
	}
	m := c.milestones[id.(int)]
	m.State = *opt.StateEvent
	return m, nil, nil
}

func (c *fakeClient) ListProjectIssues(id interface{}, opt *glab.ListProjectIssuesOptions, options ...glab.OptionFunc) ([]*glab.Issue, *glab.Response, error) {
	err := c.errors.listProjetIssues
	if err != nil {
		return nil, nil, err
	}
	if opt != nil && opt.ListOptions.Page > 1 {
		// No more pages. End of pagination.
		return nil, nil, nil
	}
	return c.issues, nil, nil
}

func (c *fakeClient) DeleteIssue(id interface{}, issue int, options ...glab.OptionFunc) (*glab.Response, error) {
	err := c.errors.deleteIssue
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *fakeClient) GetIssue(pid interface{}, id int, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	err := c.errors.getIssue
	if err != nil {
		return nil, nil, err
	}
	return c.issues[id], nil, nil
}

func (c *fakeClient) CreateIssue(pid interface{}, opt *glab.CreateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	err := c.errors.createIssue
	if err != nil {
		return nil, nil, err
	}
	i := &glab.Issue{
		ID:    len(c.issues),
		Title: *opt.Title,
	}
	for _, p := range c.issues {
		if p.Title == i.Title {
			return nil, nil, fmt.Errorf("milestone %q already exists", p.Title)
		}
	}
	c.issues = append(c.issues, i)
	return i, nil, nil
}

func (c *fakeClient) ListUsers(opt *glab.ListUsersOptions, opts ...glab.OptionFunc) ([]*glab.User, *glab.Response, error) {
	err := c.errors.listUsers
	if err != nil {
		return nil, nil, err
	}
	return nil, nil, nil
}

func (c *fakeClient) ListIssueNotes(pid interface{}, issue int, opt *glab.ListIssueNotesOptions, options ...glab.OptionFunc) ([]*glab.Note, *glab.Response, error) {
	err := c.errors.listIssueNotes
	if err != nil {
		return nil, nil, err
	}
	return nil, nil, nil
}

func (c *fakeClient) CreateIssueNote(pid interface{}, issue int, opt *glab.CreateIssueNoteOptions, options ...glab.OptionFunc) (*glab.Note, *glab.Response, error) {
	err := c.errors.createIssueNote
	if err != nil {
		return nil, nil, err
	}
	return nil, nil, nil
}

func (c *fakeClient) UpdateIssue(pid interface{}, issue int, opt *glab.UpdateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	err := c.errors.updateIssue
	if err != nil {
		return nil, nil, err
	}
	return nil, nil, nil
}

func (c *fakeClient) clearMilestones() {
	c.milestones = nil
	c.milestones = make([]*glab.Milestone, 0)
}

func (c *fakeClient) clearLabels() {
	c.labels = nil
	c.labels = make([]*glab.Label, 0)
}
