package gitlab

import (
	"net/url"

	"github.com/rotisserie/eris"
	glab "github.com/xanzy/go-gitlab"
)

// Client is an implementation of the GitLaber interface that makes real use
// of the GitLab API.
type Client struct {
	client *glab.Client
}

// DefaultClient is the default client instance to use.
var DefaultClient GitLaber

func init() {
	DefaultClient = new(Client)
}

// New GitLab client.
func New(token string, options ...glab.ClientOptionFunc) (GitLaber, error) {
	c, err := glab.NewClient(token, options...)
	if err != nil {
		return nil, eris.Wrap(err, "new client")
	}
	DefaultClient.(*Client).client = c
	return DefaultClient, err
}

func (c *Client) Client() *glab.Client {
	return c.client
}

// GetProject returns project info.
func (c *Client) GetProject(
	id interface{},
	opt *glab.GetProjectOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Project, *glab.Response, error) {
	return c.client.Projects.GetProject(id, opt, options...)
}

// CreateLabel creates a label.
func (c *Client) CreateLabel(
	id interface{},
	opt *glab.CreateLabelOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Label, *glab.Response, error) {
	return c.client.Labels.CreateLabel(id, opt, options...)
}

// ListLabels list all labels.
func (c *Client) ListLabels(
	id interface{},
	opt *glab.ListLabelsOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Label, *glab.Response, error) {
	return c.client.Labels.ListLabels(id, opt, options...)
}

// ListMilestones list all milestones.
func (c *Client) ListMilestones(
	id interface{},
	opt *glab.ListMilestonesOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Milestone, *glab.Response, error) {
	return c.client.Milestones.ListMilestones(id, opt, options...)
}

// CreateMilestone creates a milestone.
func (c *Client) CreateMilestone(
	id interface{},
	opt *glab.CreateMilestoneOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Milestone, *glab.Response, error) {
	return c.client.Milestones.CreateMilestone(id, opt, options...)
}

// UpdateMilestone updates a milestone.
func (c *Client) UpdateMilestone(
	id interface{},
	milestone int,
	opt *glab.UpdateMilestoneOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Milestone, *glab.Response, error) {
	return c.client.Milestones.UpdateMilestone(id, milestone, opt, options...)
}

// ListProjectIssues list all issues.
func (c *Client) ListProjectIssues(
	id interface{},
	opt *glab.ListProjectIssuesOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Issue, *glab.Response, error) {
	return c.client.Issues.ListProjectIssues(id, opt, options...)
}

// DeleteIssue removes an issue.
func (c *Client) DeleteIssue(
	id interface{},
	issue int,
	options ...glab.RequestOptionFunc,
) (*glab.Response, error) {
	return c.client.Issues.DeleteIssue(id, issue, options...)
}

// GetIssue returns an issue.
func (c *Client) GetIssue(
	pid interface{},
	id int,
	options ...glab.RequestOptionFunc,
) (*glab.Issue, *glab.Response, error) {
	return c.client.Issues.GetIssue(pid, id, options...)
}

// CreateIssue creates an issue.
func (c *Client) CreateIssue(
	pid interface{},
	opt *glab.CreateIssueOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Issue, *glab.Response, error) {
	return c.client.Issues.CreateIssue(pid, opt, options...)
}

// ListUsers lists all users.
func (c *Client) ListUsers(
	opt *glab.ListUsersOptions,
	opts ...glab.RequestOptionFunc,
) ([]*glab.User, *glab.Response, error) {
	return c.client.Users.ListUsers(opt, opts...)
}

// ListIssueNotes list issue notes.
func (c *Client) ListIssueNotes(
	pid interface{},
	issue int,
	opt *glab.ListIssueNotesOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Note, *glab.Response, error) {
	return c.client.Notes.ListIssueNotes(pid, issue, opt, options...)
}

// CreateIssueNote creates a note for an issue.
func (c *Client) CreateIssueNote(
	pid interface{},
	issue int,
	opt *glab.CreateIssueNoteOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Note, *glab.Response, error) {
	return c.client.Notes.CreateIssueNote(pid, issue, opt, options...)
}

// UpdateIssue updates an issue.
func (c *Client) UpdateIssue(
	pid interface{},
	issue int,
	opt *glab.UpdateIssueOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Issue, *glab.Response, error) {
	return c.client.Issues.UpdateIssue(pid, issue, opt, options...)
}

// BaseURL returns the base URL used.
func (c *Client) BaseURL() *url.URL {
	return c.client.BaseURL()
}
