package migration

import (
	"net/http"
	"net/url"

	"github.com/gotsunami/gitlab-copy/gitlab"
	glab "github.com/xanzy/go-gitlab"
)

type fakeClient struct {
	url    string
	errors struct {
		createIssue, createIssueNote, createLabel, createMilestone   error
		deleteIssue                                                  error
		getIssue, getProject                                         error
		listIssueNotes, listLabels, listMilestones, listProjetIssues error
		listUsers                                                    error
		updateIssue, updateMilestone                                 error
	}
	labels     []*glab.Label
	milestones []*glab.Milestone
}

// New fake GitLab client, for the UT.
func (c *fakeClient) New(httpClient *http.Client, token string) gitlab.GitLaber {
	return c
}

func (c *fakeClient) SetBaseURL(url string) error {
	c.url = url
	return nil
}

func (c *fakeClient) BaseURL() *url.URL {
	return nil
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
	l := &glab.Label{
		Name:        *opt.Name,
		Color:       *opt.Color,
		Description: *opt.Description,
	}
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
	m := new(glab.Milestone)
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
	return nil, nil, nil
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
	return nil, nil, nil
}

func (c *fakeClient) CreateIssue(pid interface{}, opt *glab.CreateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	err := c.errors.createIssue
	if err != nil {
		return nil, nil, err
	}
	return nil, nil, nil
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
