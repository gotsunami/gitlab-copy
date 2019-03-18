package migration

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gotsunami/gitlab-copy/config"
	"github.com/gotsunami/gitlab-copy/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	glab "github.com/xanzy/go-gitlab"
)

var dummyClient = new(fakeClient)

func init() {
	gitlab.DefaultClient = dummyClient
}

func source(m *Migration) *fakeClient {
	return m.Endpoint.SrcClient.(*fakeClient)
}

func dest(m *Migration) *fakeClient {
	return m.Endpoint.DstClient.(*fakeClient)
}

func TestMigrate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m, err := New(nil)
	assert.Error(err)

	runs := []struct {
		name    string                     // Sub-test name
		config  string                     // YAML config
		setup   func(src, dst *fakeClient) // Defines any option before calling Migrate()
		asserts func(err error, src, dst *fakeClient)
	}{
		{
			"SourceProject returns an error",
			cfg1,
			func(src, dst *fakeClient) {
				src.errors.getProject = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
				src.errors.getProject = nil
			},
		},
		{
			"copy 2 labels only",
			cfg1,
			func(src, dst *fakeClient) {
				src.labels = makeLabels("bug", "doc")
			},
			func(err error, src, dst *fakeClient) {
				require.NoError(err)
				if assert.Equal(2, len(dst.labels)) {
					assert.Equal("bug", dst.labels[0].Name)
					assert.Equal("doc", dst.labels[1].Name)
				}
			},
		},
		{
			"copy multiple pages labels",
			cfg1,
			func(src, dst *fakeClient) {
				src.labels = makeLabels("bug", "doc")
				src.labelsPage2 = makeLabels("foo", "bar")
			},
			func(err error, src, dst *fakeClient) {
				require.NoError(err)
				if assert.Equal(4, len(dst.labels)) {
					assert.Equal("bug", dst.labels[0].Name)
					assert.Equal("doc", dst.labels[1].Name)
					assert.Equal("foo", dst.labels[2].Name)
					assert.Equal("bar", dst.labels[3].Name)
				}
			},
		},
		{
			"copy 1 label and 2 issues",
			cfg2,
			func(src, dst *fakeClient) {
				src.labels = makeLabels("P0")
			},
			func(err error, src, dst *fakeClient) {
				require.NoError(err)
				if assert.Equal(1, len(dst.labels)) {
					assert.Equal("P0", dst.labels[0].Name)
				}
			},
		},
		{
			"copy milestones only",
			cfg3,
			func(src, dst *fakeClient) {
				src.clearMilestones()
				src.milestones = makeMilestones("v1", "v2")
			},
			func(err error, src, dst *fakeClient) {
				require.NoError(err)
				if assert.Equal(2, len(dst.milestones)) {
					assert.Equal("v1", dst.milestones[0].Title)
					assert.Equal("v2", dst.milestones[1].Title)
				}
			},
		},
		{
			"copy milestones only, error listing milestones",
			cfg3,
			func(src, dst *fakeClient) {
				src.clearMilestones()
				src.milestones = makeMilestones("v1")
				src.errors.listMilestones = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
				src.errors.listMilestones = nil
			},
		},
		{
			"copy milestones only, error creating milestones",
			cfg3,
			func(src, dst *fakeClient) {
				src.clearMilestones()
				src.milestones = makeMilestones("v1")
				dst.errors.createMilestone = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
				dst.errors.createMilestone = nil
			},
		},
		{
			"list labels fails",
			cfg3,
			func(src, dst *fakeClient) {
				src.errors.listLabels = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
				src.errors.listLabels = nil
			},
		},
		{
			"create labels fails",
			cfg3,
			func(src, dst *fakeClient) {
				src.clearLabels()
				src.labels = makeLabels("P0")
				dst.errors.createLabel = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				require.Error(err)
				dst.errors.createLabel = nil
			},
		},
		{
			"copy milestone only state closed",
			cfg3,
			func(src, dst *fakeClient) {
				src.clearMilestones()
				src.milestones = makeMilestones("v1")
				src.milestones[0].State = "closed"
			},
			func(err error, src, dst *fakeClient) {
				require.NoError(err)
				if assert.Equal(1, len(dst.milestones)) {
					assert.Equal("close", dst.milestones[0].State)
				}
			},
		},
		{
			"copy closed milestone fails",
			cfg3,
			func(src, dst *fakeClient) {
				src.clearMilestones()
				src.milestones = makeMilestones("v1")
				src.milestones[0].State = "closed"
				dst.errors.updateMilestone = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				require.Error(err)
				dst.errors.updateMilestone = nil
			},
		},
		{
			"copy 1 issue",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
			},
			func(err error, src, dst *fakeClient) {
				require.NoError(err)
				if assert.Len(dst.issues, 1) {
					assert.Equal("issue1", dst.issues[0].Title)
				}
			},
		},
		{
			"copy 1 issue with error",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.errors.listProjetIssues = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				require.Error(err)
			},
		},
		{
			"No fatal error if duplicate issue",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				dst.issues = makeIssues("issue1")
			},
			func(err error, src, dst *fakeClient) {
				// No error since we don't want the program to exit.
				assert.NoError(err)
			},
		},
		{
			"Move issue",
			cfg4,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
			},
		},
		{
			"No fatal error if delete issue fails",
			cfg4,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.errors.deleteIssue = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				// No error since we don't want the program to exit.
				assert.NoError(err)
			},
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			conf, err := config.Parse(strings.NewReader(run.config))
			require.NoError(err)
			// Load the conf.
			m, err = New(conf)
			require.NoError(err)
			// Setup.
			run.setup(source(m), dest(m))
			// Run the migration.
			err = m.Migrate()
			// Asserts and tear down.
			run.asserts(err, source(m), dest(m))
		})
	}
}

func TestMigrateIssue(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	runs := []struct {
		name    string                     // Sub-test name
		config  string                     // YAML config
		setup   func(src, dst *fakeClient) // Defines any option before calling Migrate()
		asserts func(err error, src, dst *fakeClient)
	}{
		{
			"Duplicate issue fails",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				dst.issues = makeIssues("issue1")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Get issue fails",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.errors.getIssue = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"List project issues fails",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				dst.errors.listProjetIssues = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Assigned user, but no target user match",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issues[0].Assignee.Username = "mat"
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
				assert.Empty(dst.issues[0].Assignee.Username)
			},
		},
		{
			"Assigned user, fetching users fails",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issues[0].Assignee.Username = "mat"
				dst.errors.listUsers = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Issue has assigned user, target user match",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issues[0].Assignee.Username = "mat"
				dst.users = makeUsers("mat")
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
				if assert.Len(dst.issues, 1) {
					assert.Equal("mat", dst.issues[0].Assignee.Username)
				}
			},
		},
		{
			"Issue has a milestone, not target match",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				m := &glab.Milestone{
					Title: "v1.0",
				}
				src.issues[0].Milestone = m
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
				if assert.Len(dst.milestones, 1) {
					assert.Equal("v1.0", dst.milestones[0].Title)
				}
			},
		},
		{
			"Issue has a milestone, create target milestone error",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				m := &glab.Milestone{
					Title: "v1.0",
				}
				src.issues[0].Milestone = m
				dst.errors.createMilestone = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Issue has a milestone, list target milestones error",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				m := &glab.Milestone{
					Title: "v1.0",
				}
				src.issues[0].Milestone = m
				dst.errors.listMilestones = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Issue has a milestone, found on target",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				m := &glab.Milestone{
					Title: "v1.0",
				}
				src.issues[0].Milestone = m
				dst.milestones = makeMilestones("v1.0")
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
				assert.Len(dst.milestones, 1)
			},
		},
		{
			"Copy existing labels",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issues[0].Labels = []string{"P1", "P2"}
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
			},
		},
		{
			"Failing creating target issue",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				dst.errors.createIssue = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Failing creating target issue, URI too long HTTP error, empty issue description",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				dst.errors.createIssue = errors.New("err")
				dst.httpErrorRaiseURITooLong = true
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Failing creating target issue, URI too long HTTP error, with issue description",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				buf := make([]byte, 1128)
				desc := bytes.NewBuffer(buf)
				desc.WriteString("Some desc")
				src.issues[0].Description = desc.String()
				dst.errors.createIssue = errors.New("err")
				dst.httpErrorRaiseURITooLong = true
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"List issue notes fails",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.errors.listIssueNotes = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Issue has notes",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issueNotes = makeNotes("n1", "n2")
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
			},
		},
		{
			"Issue with notes, create issue note error",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issueNotes = makeNotes("n1", "n2")
				dst.errors.createIssueNote = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Issue with notes, but description too long",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issueNotes = makeNotes("n1", "n2")
				// Large data buffer to raise an URITooLong error.
				buf := make([]byte, 1128)
				desc := bytes.NewBuffer(buf)
				desc.WriteString("Some desc")
				src.issueNotes[1].Body = desc.String()
				dst.httpErrorRaiseURITooLong = true
				dst.errors.createIssueNote = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
		{
			"Closed issue",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issues[0].State = "closed"
			},
			func(err error, src, dst *fakeClient) {
				assert.NoError(err)
			},
		},
		{
			"Closed issue, with error when updating target issue",
			cfg2,
			func(src, dst *fakeClient) {
				src.issues = makeIssues("issue1")
				src.issues[0].State = "closed"
				dst.errors.updateIssue = errors.New("err")
			},
			func(err error, src, dst *fakeClient) {
				assert.Error(err)
			},
		},
	}
	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			conf, err := config.Parse(strings.NewReader(run.config))
			require.NoError(err)
			// Load the conf.
			m, err := New(conf)
			require.NoError(err)
			_, err = m.SourceProject(m.params.SrcPrj.Name)
			require.NoError(err)
			_, err = m.DestProject(m.params.DstPrj.Name)
			require.NoError(err)
			// Setup.
			run.setup(source(m), dest(m))
			// Run the migration.
			err = m.migrateIssue(0)
			// Asserts and tear down.
			run.asserts(err, source(m), dest(m))
		})
	}
}

func makeLabels(names ...string) []*glab.Label {
	labels := make([]*glab.Label, len(names))
	for k, n := range names {
		labels[k] = &glab.Label{
			ID:   k,
			Name: n,
		}
	}
	return labels
}

func makeMilestones(names ...string) []*glab.Milestone {
	ms := make([]*glab.Milestone, len(names))
	for k, n := range names {
		ms[k] = &glab.Milestone{
			ID:    k,
			Title: n,
		}
	}
	return ms
}

func makeIssues(names ...string) []*glab.Issue {
	issues := make([]*glab.Issue, len(names))
	for k, n := range names {
		issues[k] = &glab.Issue{
			ID:    k,
			Title: n,
		}
	}
	return issues
}

func makeUsers(names ...string) []*glab.User {
	users := make([]*glab.User, len(names))
	for k, n := range names {
		users[k] = &glab.User{
			ID:       k,
			Username: n,
		}
	}
	return users
}

func makeNotes(names ...string) []*glab.Note {
	notes := make([]*glab.Note, len(names))
	now := time.Now()
	for k, n := range names {
		notes[k] = &glab.Note{
			ID:        k,
			Title:     n,
			CreatedAt: &now,
		}
		notes[k].Author.Name = "me"
		notes[k].Author.Username = "me"
	}
	return notes
}
