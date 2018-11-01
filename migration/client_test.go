package migration

import (
	"net/http"
	"net/url"

	"github.com/gotsunami/gitlab-copy/gitlab"
	glab "github.com/xanzy/go-gitlab"
)

type fakeClient struct {
	url        string
	err        error
	labels     []*glab.Label
	milestones []*glab.Milestone
}

// New fake GitLab client, for the UT.
func (c *fakeClient) New(httpClient *http.Client, token string) gitlab.GitLaber {
	return c
}

func (c *fakeClient) SetBaseURL(url string) error {
	c.url = url
	return c.err
}

func (c *fakeClient) BaseURL() *url.URL {
	return nil
}

func (c *fakeClient) GetProject(interface{}, ...glab.OptionFunc) (*glab.Project, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
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
	if c.err != nil {
		return nil, nil, c.err
	}
	l := &glab.Label{
		Name:        *opt.Name,
		Color:       *opt.Color,
		Description: *opt.Description,
	}
	return l, nil, nil
}

func (c *fakeClient) ListLabels(id interface{}, opt *glab.ListLabelsOptions, options ...glab.OptionFunc) ([]*glab.Label, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return c.labels, nil, nil
}

func (c *fakeClient) ListMilestones(id interface{}, opt *glab.ListMilestonesOptions, options ...glab.OptionFunc) ([]*glab.Milestone, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return c.milestones, nil, nil
}

func (c *fakeClient) CreateMilestone(id interface{}, opt *glab.CreateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) UpdateMilestone(id interface{}, milestone int, opt *glab.UpdateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) ListProjectIssues(id interface{}, opt *glab.ListProjectIssuesOptions, options ...glab.OptionFunc) ([]*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) DeleteIssue(id interface{}, issue int, options ...glab.OptionFunc) (*glab.Response, error) {
	if c.err != nil {
		return nil, c.err
	}
	return nil, nil
}

func (c *fakeClient) GetIssue(pid interface{}, id int, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) CreateIssue(pid interface{}, opt *glab.CreateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) ListUsers(opt *glab.ListUsersOptions, opts ...glab.OptionFunc) ([]*glab.User, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) ListIssueNotes(pid interface{}, issue int, opt *glab.ListIssueNotesOptions, options ...glab.OptionFunc) ([]*glab.Note, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) CreateIssueNote(pid interface{}, issue int, opt *glab.CreateIssueNoteOptions, options ...glab.OptionFunc) (*glab.Note, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *fakeClient) UpdateIssue(pid interface{}, issue int, opt *glab.UpdateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}
