package gitlab

import (
	"net/http"
	"net/url"

	glab "github.com/xanzy/go-gitlab"
)

type Client struct {
	client *glab.Client
}

// New GitLab client.
func New(httpClient *http.Client, token string) GitLaber {
	c := new(Client)
	c.client = glab.NewClient(httpClient, token)
	return c
}

func (c *Client) SetBaseURL(url string) error {
	return c.client.SetBaseURL(url)
}

func (c *Client) GetProject(id interface{}, options ...glab.OptionFunc) (*glab.Project, *glab.Response, error) {
	return c.client.Projects.GetProject(id, options...)
}

func (c *Client) CreateLabel(id interface{}, opt *glab.CreateLabelOptions, options ...glab.OptionFunc) (*glab.Label, *glab.Response, error) {
	return c.client.Labels.CreateLabel(id, opt, options...)
}

func (c *Client) ListLabels(id interface{}, opt *glab.ListLabelsOptions, options ...glab.OptionFunc) ([]*glab.Label, *glab.Response, error) {
	return c.client.Labels.ListLabels(id, opt, options...)
}

func (c *Client) ListMilestones(id interface{}, opt *glab.ListMilestonesOptions, options ...glab.OptionFunc) ([]*glab.Milestone, *glab.Response, error) {
	return c.client.Milestones.ListMilestones(id, opt, options...)
}

func (c *Client) CreateMilestone(id interface{}, opt *glab.CreateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	return c.client.Milestones.CreateMilestone(id, opt, options...)
}

func (c *Client) UpdateMilestone(id interface{}, milestone int, opt *glab.UpdateMilestoneOptions, options ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error) {
	return c.client.Milestones.UpdateMilestone(id, milestone, opt, options...)
}

func (c *Client) ListProjectIssues(id interface{}, opt *glab.ListProjectIssuesOptions, options ...glab.OptionFunc) ([]*glab.Issue, *glab.Response, error) {
	return c.client.Issues.ListProjectIssues(id, opt, options...)
}

func (c *Client) DeleteIssue(id interface{}, issue int, options ...glab.OptionFunc) (*glab.Response, error) {
	return c.client.Issues.DeleteIssue(id, issue, options...)
}

func (c *Client) GetIssue(pid interface{}, id int, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	return c.client.Issues.GetIssue(pid, id, options...)
}
func (c *Client) CreateIssue(pid interface{}, opt *glab.CreateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	return c.client.Issues.CreateIssue(pid, opt, options...)
}

func (c *Client) ListUsers(opt *glab.ListUsersOptions, opts ...glab.OptionFunc) ([]*glab.User, *glab.Response, error) {
	return c.client.Users.ListUsers(opt, opts...)
}

func (c *Client) ListIssueNotes(pid interface{}, issue int, opt *glab.ListIssueNotesOptions, options ...glab.OptionFunc) ([]*glab.Note, *glab.Response, error) {
	return c.client.Notes.ListIssueNotes(pid, issue, opt, options...)
}

func (c *Client) CreateIssueNote(pid interface{}, issue int, opt *glab.CreateIssueNoteOptions, options ...glab.OptionFunc) (*glab.Note, *glab.Response, error) {
	return c.client.Notes.CreateIssueNote(pid, issue, opt, options...)
}

func (c *Client) UpdateIssue(pid interface{}, issue int, opt *glab.UpdateIssueOptions, options ...glab.OptionFunc) (*glab.Issue, *glab.Response, error) {
	return c.client.Issues.UpdateIssue(pid, issue, opt, options...)
}

func (c *Client) BaseURL() *url.URL {
	return c.client.BaseURL()
}
