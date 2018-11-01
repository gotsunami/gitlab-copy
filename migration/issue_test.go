package migration

import (
	"strings"
	"testing"

	"github.com/gotsunami/gitlab-copy/config"
	"github.com/gotsunami/gitlab-copy/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	glab "github.com/xanzy/go-gitlab"
)

const cfg1 = `
from:
    url: https://gitlab.mydomain.com
    token: sourcetoken
    project: source/project
#    issues:
#    - 5
#    - 8-10
    labelsOnly: true
    # moveIssues: true
to:
    url: https://gitlab.mydomain.com
    token: desttoken
    project: dest/project
`

const cfg2 = `
from:
    url: https://gitlab.mydomain.com
    token: sourcetoken
    project: source/project
to:
    url: https://gitlab.mydomain.com
    token: desttoken
    project: dest/project
`

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
		asserts func(*Migration)
	}{
		{
			"copy 2 labels only",
			cfg1,
			func() {
				dummyClient.labels = makeLabels("bug", "doc")
			},
			func(m *Migration) {
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
			func(m *Migration) {
				fk := m.Endpoint.DstClient.(*fakeClient)
				assert.Equal(1, len(fk.labels))
				assert.Equal("P0", fk.labels[0].Name)
			},
		},
	}

	for _, run := range runs {
		conf, err := config.Parse(strings.NewReader(run.config))
		require.NoError(err)
		m, err = New(conf)
		require.NoError(err)
		t.Run(run.name, func(t *testing.T) {
			run.setup()
			err = m.Migrate()
			assert.NoError(err)
			run.asserts(m)
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
