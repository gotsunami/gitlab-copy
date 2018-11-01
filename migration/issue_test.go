package migration

import (
	"strings"
	"testing"

	"github.com/gotsunami/gitlab-copy/config"
	"github.com/gotsunami/gitlab-copy/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cfg = `
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

func init() {
	gitlab.DefaultClient = new(fakeClient)
}

func TestMigrate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m, err := New(nil)
	assert.Error(err)

	conf, err := config.Parse(strings.NewReader(cfg))
	require.NoError(err)

	m, err = New(conf)
	require.NoError(err)

	err = m.Migrate()
	assert.NoError(err)
}
