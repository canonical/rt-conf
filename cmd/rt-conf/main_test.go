package main

import (
	"strings"
	"testing"
)

func TestRunNoConfigPath(t *testing.T) {
	err := run([]string{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "path not set") {
		t.Fatalf("expected 'failed to load config file', got: %v", err)
	}
}

func TestRunFlagParseError(t *testing.T) {
	err := run([]string{"-bogus"})
	if err == nil || !strings.Contains(err.Error(), "failed to parse flags") {
		t.Fatalf("expected flag parse error, got %v", err)
	}
}

func TestRunConfigLoadError(t *testing.T) {
	err := run([]string{"-file", "/does/not/exist"})
	if err == nil || !strings.Contains(err.Error(), "failed to find file") {
		t.Fatalf("expected config load error, got %v", err)
	}
}
