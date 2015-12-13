package main

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	if _, err := parseConfig("wrongfile"); err == nil {
		t.Fatal("expect an error mon missing config file, got nil")
	}
}
