package migration

import (
	"errors"
	"strings"
	"testing"

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
				assert.NoError(err)
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
