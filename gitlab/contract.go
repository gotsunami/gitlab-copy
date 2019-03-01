package gitlab

import (
	"net/http"
	"net/url"

	glab "github.com/xanzy/go-gitlab"
)

// GitLaber defines some methods of glab.Client so it can be mocked easily in
// the unit tests.
type GitLaber interface {
	New(*http.Client, string) GitLaber
	BaseURL() *url.URL
	SetBaseURL(string) error
	// Project
	GetProject(interface{}, ...glab.OptionFunc) (*glab.Project, *glab.Response, error)
	// Labels
	ListLabels(interface{}, *glab.ListLabelsOptions, ...glab.OptionFunc) ([]*glab.Label, *glab.Response, error)
	CreateLabel(interface{}, *glab.CreateLabelOptions, ...glab.OptionFunc) (*glab.Label, *glab.Response, error)
	// Milestones
	ListMilestones(interface{}, *glab.ListMilestonesOptions, ...glab.OptionFunc) ([]*glab.Milestone, *glab.Response, error)
	CreateMilestone(interface{}, *glab.CreateMilestoneOptions, ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error)
	UpdateMilestone(interface{}, int, *glab.UpdateMilestoneOptions, ...glab.OptionFunc) (*glab.Milestone, *glab.Response, error)
	// Issues
	ListProjectIssues(interface{}, *glab.ListProjectIssuesOptions, ...glab.OptionFunc) ([]*glab.Issue, *glab.Response, error)
	GetIssue(interface{}, int, ...glab.OptionFunc) (*glab.Issue, *glab.Response, error)
	CreateIssue(interface{}, *glab.CreateIssueOptions, ...glab.OptionFunc) (*glab.Issue, *glab.Response, error)
	UpdateIssue(interface{}, int, *glab.UpdateIssueOptions, ...glab.OptionFunc) (*glab.Issue, *glab.Response, error)
	DeleteIssue(interface{}, int, ...glab.OptionFunc) (*glab.Response, error)
	// Users
	ListUsers(*glab.ListUsersOptions, ...glab.OptionFunc) ([]*glab.User, *glab.Response, error)
	// Notes
	ListIssueNotes(interface{}, int, *glab.ListIssueNotesOptions, ...glab.OptionFunc) ([]*glab.Note, *glab.Response, error)
	CreateIssueNote(interface{}, int, *glab.CreateIssueNoteOptions, ...glab.OptionFunc) (*glab.Note, *glab.Response, error)
}
