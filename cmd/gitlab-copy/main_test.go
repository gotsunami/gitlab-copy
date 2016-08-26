package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	doc := `
from:
  token: srctoken
  project: srcproj
to:
  token: dsttoken
  project: dstproj
`
	t.Log(doc)
}

func TestMap2Human(t *testing.T) {
	cc := []struct {
		m   map[string]int
		exp func(map[string]int) error
	}{
		{map[string]int{"a": 0, "b": 0, "c": 0},
			func(p map[string]int) error {
				// FIXME
				return nil
			},
		},
	}
	for _, v := range cc {
		if err := v.exp(v.m); err != nil {
			t.Errorf("expect '%s', got '%s'", v.exp, map2Human(v.m))
		}
	}
}
