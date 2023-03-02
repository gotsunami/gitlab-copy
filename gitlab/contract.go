package gitlab

import (
	"net/url"

	glab "github.com/xanzy/go-gitlab"
)

// GitLaber defines some methods of glab.Client so it can be mocked easily in
// the unit tests.
type GitLaber interface {
	BaseURL() *url.URL
	Client() *glab.Client
	// Project
	GetProject(interface{}, *glab.GetProjectOptions, ...glab.RequestOptionFunc) (*glab.Project, *glab.Response, error)
	// Labels
	ListLabels(interface{}, *glab.ListLabelsOptions, ...glab.RequestOptionFunc) ([]*glab.Label, *glab.Response, error)
	CreateLabel(interface{}, *glab.CreateLabelOptions, ...glab.RequestOptionFunc) (*glab.Label, *glab.Response, error)
	// Milestones
	ListMilestones(interface{}, *glab.ListMilestonesOptions, ...glab.RequestOptionFunc) ([]*glab.Milestone, *glab.Response, error)
	CreateMilestone(interface{}, *glab.CreateMilestoneOptions, ...glab.RequestOptionFunc) (*glab.Milestone, *glab.Response, error)
	UpdateMilestone(interface{}, int, *glab.UpdateMilestoneOptions, ...glab.RequestOptionFunc) (*glab.Milestone, *glab.Response, error)
	// Issues
	ListProjectIssues(interface{}, *glab.ListProjectIssuesOptions, ...glab.RequestOptionFunc) ([]*glab.Issue, *glab.Response, error)
	GetIssue(interface{}, int, ...glab.RequestOptionFunc) (*glab.Issue, *glab.Response, error)
	CreateIssue(interface{}, *glab.CreateIssueOptions, ...glab.RequestOptionFunc) (*glab.Issue, *glab.Response, error)
	UpdateIssue(interface{}, int, *glab.UpdateIssueOptions, ...glab.RequestOptionFunc) (*glab.Issue, *glab.Response, error)
	DeleteIssue(interface{}, int, ...glab.RequestOptionFunc) (*glab.Response, error)
	// Users
	ListUsers(*glab.ListUsersOptions, ...glab.RequestOptionFunc) ([]*glab.User, *glab.Response, error)
	// Notes
	ListIssueNotes(interface{}, int, *glab.ListIssueNotesOptions, ...glab.RequestOptionFunc) ([]*glab.Note, *glab.Response, error)
	CreateIssueNote(interface{}, int, *glab.CreateIssueNoteOptions, ...glab.RequestOptionFunc) (*glab.Note, *glab.Response, error)
}
