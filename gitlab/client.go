package gitlab

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rotisserie/eris"
	glab "github.com/xanzy/go-gitlab"
)

// client is an implementation of the GitLaber interface that makes real use
// of the GitLab API.
type client struct {
	c *glab.Client
}

// NewClient returns a new client.
func NewClient() GitLaber {
	return new(client)
}

var skipTLSVerification bool

// SkipTLSVerificationProcess skips the TLS verification process by using a custom HTTP transport.
func SkipTLSVerificationProcess() {
	skipTLSVerification = true
}

// WithToken sets the token to use, along with any client options.
func (c *client) WithToken(token string, options ...glab.ClientOptionFunc) (GitLaber, error) {
	f := new(client)

	if skipTLSVerification {
		// Setup a custom HTTP client to ignore TLS issues.
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc := retryablehttp.NewClient()
		hc.HTTPClient.Transport = tr
		options = append(options, glab.WithHTTPClient(hc.StandardClient()))
	}

	p, err := glab.NewClient(token, options...)
	if err != nil {
		return nil, eris.Wrap(err, "with token")
	}
	f.c = p
	return f, nil
}

// GitLab returns the GitLab client.
func (c *client) GitLab() *glab.Client {
	return c.c
}

// GetProject returns project info.
func (c *client) GetProject(
	id interface{},
	opt *glab.GetProjectOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Project, *glab.Response, error) {
	return c.c.Projects.GetProject(id, opt, options...)
}

// CreateLabel creates a label.
func (c *client) CreateLabel(
	id interface{},
	opt *glab.CreateLabelOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Label, *glab.Response, error) {
	return c.c.Labels.CreateLabel(id, opt, options...)
}

// ListLabels list all labels.
func (c *client) ListLabels(
	id interface{},
	opt *glab.ListLabelsOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Label, *glab.Response, error) {
	return c.c.Labels.ListLabels(id, opt, options...)
}

// ListMilestones list all milestones.
func (c *client) ListMilestones(
	id interface{},
	opt *glab.ListMilestonesOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Milestone, *glab.Response, error) {
	return c.c.Milestones.ListMilestones(id, opt, options...)
}

// CreateMilestone creates a milestone.
func (c *client) CreateMilestone(
	id interface{},
	opt *glab.CreateMilestoneOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Milestone, *glab.Response, error) {
	return c.c.Milestones.CreateMilestone(id, opt, options...)
}

// UpdateMilestone updates a milestone.
func (c *client) UpdateMilestone(
	id interface{},
	milestone int,
	opt *glab.UpdateMilestoneOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Milestone, *glab.Response, error) {
	return c.c.Milestones.UpdateMilestone(id, milestone, opt, options...)
}

// ListProjectIssues list all issues.
func (c *client) ListProjectIssues(
	id interface{},
	opt *glab.ListProjectIssuesOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Issue, *glab.Response, error) {
	return c.c.Issues.ListProjectIssues(id, opt, options...)
}

// DeleteIssue removes an issue.
func (c *client) DeleteIssue(
	id interface{},
	issue int,
	options ...glab.RequestOptionFunc,
) (*glab.Response, error) {
	return c.c.Issues.DeleteIssue(id, issue, options...)
}

// GetIssue returns an issue.
func (c *client) GetIssue(
	pid interface{},
	id int,
	options ...glab.RequestOptionFunc,
) (*glab.Issue, *glab.Response, error) {
	return c.c.Issues.GetIssue(pid, id, options...)
}

// CreateIssue creates an issue.
func (c *client) CreateIssue(
	pid interface{},
	opt *glab.CreateIssueOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Issue, *glab.Response, error) {
	return c.c.Issues.CreateIssue(pid, opt, options...)
}

// ListUsers lists all users.
func (c *client) ListUsers(
	opt *glab.ListUsersOptions,
	opts ...glab.RequestOptionFunc,
) ([]*glab.User, *glab.Response, error) {
	return c.c.Users.ListUsers(opt, opts...)
}

// ListIssueNotes list issue notes.
func (c *client) ListIssueNotes(
	pid interface{},
	issue int,
	opt *glab.ListIssueNotesOptions,
	options ...glab.RequestOptionFunc,
) ([]*glab.Note, *glab.Response, error) {
	return c.c.Notes.ListIssueNotes(pid, issue, opt, options...)
}

// CreateIssueNote creates a note for an issue.
func (c *client) CreateIssueNote(
	pid interface{},
	issue int,
	opt *glab.CreateIssueNoteOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Note, *glab.Response, error) {
	return c.c.Notes.CreateIssueNote(pid, issue, opt, options...)
}

// UpdateIssue updates an issue.
func (c *client) UpdateIssue(
	pid interface{},
	issue int,
	opt *glab.UpdateIssueOptions,
	options ...glab.RequestOptionFunc,
) (*glab.Issue, *glab.Response, error) {
	return c.c.Issues.UpdateIssue(pid, issue, opt, options...)
}

// BaseURL returns the base URL used.
func (c *client) BaseURL() *url.URL {
	return c.c.BaseURL()
}
