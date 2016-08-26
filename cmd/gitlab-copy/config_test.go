package main

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	if _, err := parseConfig("wrongfile"); err == nil {
		t.Fatal("expect an error mon missing config file, got nil")
	}
}

func TestParseIssues(t *testing.T) {
	issues := []struct {
		ranges     []string
		shouldFail bool
		expect     []issueRange
	}{
		{[]string{"3", "5-10", "Z"}, true, nil},
		{[]string{"1", "a-10"}, true, nil},
		{[]string{"1-3-5"}, true, nil},
		{[]string{"4-8"}, false, []issueRange{{4, 8}}},
		{[]string{"5-5"}, false, []issueRange{{5, 5}}},
		{[]string{"2", "5-15", "3"}, false, []issueRange{{2, 2}, {5, 15}, {3, 3}}},
		{[]string{"4-3"}, true, nil},
	}

	p := new(project)
	for _, k := range issues {
		p.Issues = k.ranges
		err := p.parseIssues()
		if err == nil && k.shouldFail {
			t.Errorf("expects an error for '%s', got nil", p.Issues)
		}
		if err != nil && !k.shouldFail {
			t.Errorf("expects no error for '%s', got one: %s", p.Issues, err.Error())
		}
		if k.expect != nil {
			for j, r := range k.expect {
				if r.from != p.issues[j].from || r.to != p.issues[j].to {
					t.Errorf("range mismatch: expects %v, got %v", r, p.issues[j])
				}
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
		if m := p.matches(r.val); m != r.match {
			t.Errorf("expected: %v, got a match: %v", r.match, m)
		}
	}
}
