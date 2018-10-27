package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	require := require.New(t)
	_, err := Parse("wrongfile")
	require.NotNil(err)
}

func TestParseIssues(t *testing.T) {
	assert := assert.New(t)

	issues := []struct {
		name       string
		ranges     []string
		shouldFail bool
		expect     []issueRange
	}{
		{"Standalone char in range sequence", []string{"3", "5-10", "Z"}, true, nil},
		{"Char in range sequence", []string{"1", "a-10"}, true, nil},
		{"Malformed range", []string{"1-3-5"}, true, nil},
		{"Valid range, distinct numbers", []string{"4-8"}, false, []issueRange{{4, 8}}},
		{"Valid range, same number", []string{"5-5"}, false, []issueRange{{5, 5}}},
		{"Valid ranges, single number", []string{"2", "5-15", "3"}, false, []issueRange{{2, 2}, {5, 15}, {3, 3}}},
		{"Invalid range bounds", []string{"4-3"}, true, nil},
	}

	p := new(project)
	for _, k := range issues {
		p.Issues = k.ranges
		err := p.parseIssues()
		if err == nil && k.shouldFail {
			t.Errorf("expects an error for %q, got nil", p.Issues)
		}
		if err != nil && !k.shouldFail {
			t.Errorf("expects no error for %q, got one: %s", p.Issues, err.Error())
		}
		if k.expect != nil {
			for j, r := range k.expect {
				t.Run(k.name, func(t *testing.T) {
					assert.Equal(r.from, p.issues[j].from)
					assert.Equal(r.to, p.issues[j].to)
				})
			}
		}
	}
}

func TestMatches(t *testing.T) {
	p := new(project)

	set := []struct {
		ranges []issueRange
		val    int
		match  bool
	}{
		{[]issueRange{{4, 8}}, 5, true},
		{[]issueRange{}, 1, true},
		{[]issueRange{{4, 8}}, 9, false},
		{[]issueRange{{2, 2}}, 2, true},
	}
	for _, r := range set {
		p.issues = r.ranges
		if m := p.Matches(r.val); m != r.match {
			t.Errorf("expected: %v, got a match: %v", r.match, m)
		}
	}
}
