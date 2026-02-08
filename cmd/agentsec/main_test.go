package main

import "testing"

func TestVersionString(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	t.Cleanup(func() {
		version = origVersion
		commit = origCommit
		date = origDate
	})

	version = "dev"
	commit = "none"
	date = "unknown"
	if got := versionString(); got != "agentsec dev" {
		t.Fatalf("unexpected version string: %q", got)
	}

	version = "v1.2.3"
	commit = "abc1234"
	date = "2026-02-08T18:00:00Z"
	if got := versionString(); got != "agentsec v1.2.3 (commit=abc1234, built=2026-02-08T18:00:00Z)" {
		t.Fatalf("unexpected version string with metadata: %q", got)
	}
}
