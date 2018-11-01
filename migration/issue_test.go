package migration

import (
	"strings"
	"testing"

	"github.com/gotsunami/gitlab-copy/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cfg = `
from:
    url: https://gitlab.keeneyetechnologies.com
    token: CKZtsLqaHVVryZxebVmt
    project: core/opencv-contrib
#    issues:
#    - 5
#    - 8-10
    labelsOnly: true
    # moveIssues: true
to:
    url: https://gitlab.keeneyetechnologies.com
    token: CKZtsLqaHVVryZxebVmt
    project: core/io
`

func TestParseConfig(t *testing.T) {
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
