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

func TestMigrate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m, err := New(nil)
	assert.Error(err)

	runs := []struct {
		name    string // Sub-test name
		config  string // YAML config
		setup   func() // Defines any option before calling Migrate()
		asserts func(error, *Migration)
	}{
		{
			"SourceProject returns an error",
			cfg1,
			func() {
				dummyClient.errors.getProject = errors.New("err")
			},
			func(err error, m *Migration) {
				assert.Error(err)
				dummyClient.errors.getProject = nil
			},
		},
		{
			"copy 2 labels only",
			cfg1,
			func() {
				dummyClient.labels = makeLabels("bug", "doc")
			},
			func(err error, m *Migration) {
				require.NoError(err)
				fk := m.Endpoint.DstClient.(*fakeClient)
				assert.Equal(2, len(fk.labels))
				assert.Equal("bug", fk.labels[0].Name)
				assert.Equal("doc", fk.labels[1].Name)
			},
		},
		{
			"copy 1 label and 2 issues",
			cfg2,
			func() {
				dummyClient.labels = makeLabels("P0")
			},
			func(err error, m *Migration) {
				require.NoError(err)
				fk := m.Endpoint.DstClient.(*fakeClient)
				assert.Equal(1, len(fk.labels))
				assert.Equal("P0", fk.labels[0].Name)
			},
		},
		{
			"copy milestones only",
			cfg3,
			func() {
				dummyClient.milestones = makeMilestones("v1", "v2")
			},
			func(err error, m *Migration) {
				require.NoError(err)
				fk := m.Endpoint.DstClient.(*fakeClient)
				assert.Equal(2, len(fk.milestones))
				assert.Equal("v1", fk.milestones[0].Title)
				assert.Equal("v2", fk.milestones[1].Title)
			},
		},
		{
			"copy milestones only, error listing milestones",
			cfg3,
			func() {
				dummyClient.milestones = makeMilestones("v1")
				dummyClient.errors.listMilestones = errors.New("err")
			},
			func(err error, m *Migration) {
				assert.Error(err)
				dummyClient.errors.listMilestones = nil
			},
		},
		{
			"copy milestones only, error creating milestones",
			cfg3,
			func() {
				dummyClient.milestones = makeMilestones("v1")
				dummyClient.errors.createMilestone = errors.New("err")
			},
			func(err error, m *Migration) {
				assert.Error(err)
				dummyClient.errors.createMilestone = nil
			},
		},
		{
			"list labels fails",
			cfg3,
			func() {
				dummyClient.errors.listLabels = errors.New("err")
			},
			func(err error, m *Migration) {
				assert.Error(err)
				dummyClient.errors.listLabels = nil
			},
		},
		{
			"create labels fails",
			cfg3,
			func() {
				dummyClient.errors.createLabel = errors.New("err")
			},
			func(err error, m *Migration) {
				assert.Error(err)
				dummyClient.errors.createLabel = nil
			},
		},
		{
			"copy milestone only state closed",
			cfg3,
			func() {
				dummyClient.milestones = makeMilestones("v1")
				dummyClient.milestones[0].State = "closed"
			},
			func(err error, m *Migration) {
				require.NoError(err)
				fk := m.Endpoint.DstClient.(*fakeClient)
				assert.Equal(1, len(fk.milestones))
				assert.Equal("close", fk.milestones[0].State)
			},
		},
		{
			"copy closed milestone fails",
			cfg3,
			func() {
				dummyClient.milestones = makeMilestones("v1")
				dummyClient.milestones[0].State = "closed"
				dummyClient.errors.updateMilestone = errors.New("err")
			},
			func(err error, m *Migration) {
				assert.Error(err)
				dummyClient.errors.updateMilestone = nil
			},
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			conf, err := config.Parse(strings.NewReader(run.config))
			require.NoError(err)
			m, err = New(conf)
			require.NoError(err)
			run.setup()
			err = m.Migrate()
			run.asserts(err, m)
		})
	}
}

func makeLabels(names ...string) []*glab.Label {
	labels := make([]*glab.Label, len(names))
	for k, n := range names {
		labels[k] = &glab.Label{
			Name: n,
		}
	}
	return labels
}

func makeMilestones(names ...string) []*glab.Milestone {
	ms := make([]*glab.Milestone, len(names))
	for k, n := range names {
		ms[k] = &glab.Milestone{
			Title: n,
		}
	}
	return ms
}
