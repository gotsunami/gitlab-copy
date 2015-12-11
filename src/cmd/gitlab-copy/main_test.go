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
