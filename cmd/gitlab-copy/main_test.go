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
		exp string
	}{
		{map[string]int{"a": 0, "b": 0, "c": 0}, "a, b, c"},
	}
	for _, v := range cc {
		if map2Human(v.m) != v.exp {
			t.Errorf("expect '%s', got '%s'", v.exp, map2Human(v.m))
		}
	}
}
