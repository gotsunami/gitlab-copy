package gitlab

import (
	"net/http"
	"net/url"

	glab "github.com/xanzy/go-gitlab"
)

type FakeClient struct {
	url        string
	err        error
	labels     []*glab.Label
	milestones []*glab.Milestone
}

// New fake GitLab client, for the UT.
func newFake(httpClient *http.Client, token string) GitLaber {
	c := new(FakeClient)
	return c
}

func (c *FakeClient) SetBaseURL(url string) error {
	c.url = url
	return c.err
}

func (c *FakeClient) BaseURL() *url.URL {
	return nil
}

func (c *FakeClient) GetProject(interface{}, ...glab.OptionFunc) (*glab.Project, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	p := new(glab.Project)
	p.Name = "A name"
	return p, nil, nil
}

func (c *FakeClient) CreateLabel(id interface{}, opt *glab.CreateLabelOptions, options ...glab.OptionFunc) (*glab.Label, *glab.Response, error) {
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

func (c *FakeClient) ListLabels(id interface{}, opt *glab.ListLabelsOptions, options ...glab.OptionFunc) ([]*glab.Label, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return c.labels, nil, nil
}

func (c *FakeClient) ListMilestones(id interface{}, opt *glab.ListMilestonesOptions, options ...glab.OptionFunc) ([]*glab.Milestone, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return c.milestones, nil, nil
}

func (c *FakeClient) CreateMilestone(id interface{}, opt *glab.CreateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) UpdateMilestone(id interface{}, milestone int, opt *glab.UpdateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) ListProjectIssues(id interface{}, opt *glab.ListProjectIssuesOptions, options ...glab.OptionFunc) ([]*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) DeleteIssue(id interface{}, issue int, options ...glab.OptionFunc) (*glab.Response, error) {
	if c.err != nil {
		return nil, c.err
	}
	return nil, nil
}

func (c *FakeClient) GetIssue(pid interface{}, id int, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) CreateIssue(pid interface{}, opt *glab.CreateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) ListUsers(opt *glab.ListUsersOptions, opts ...glab.OptionFunc) ([]*glab.User, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) ListIssueNotes(pid interface{}, issue int, opt *glab.ListIssueNotesOptions, options ...glab.OptionFunc) ([]*glab.Note, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) CreateIssueNote(pid interface{}, issue int, opt *glab.CreateIssueNoteOptions, options ...glab.OptionFunc) (*glab.Note, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}

func (c *FakeClient) UpdateIssue(pid interface{}, issue int, opt *glab.UpdateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	return nil, nil, nil
}
