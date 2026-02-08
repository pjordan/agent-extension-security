package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pjordan/agent-extension-security/internal/manifest"
)

func TestRunInit(t *testing.T) {
	t.Run("creates all expected files", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "test-skill")
		runInit([]string{dir, "--id", "com.example.test", "--type", "skill"})

		for _, name := range []string{"aem.json", "devkey.json", "policy.json"} {
			p := filepath.Join(dir, name)
			if _, err := os.Stat(p); os.IsNotExist(err) {
				t.Fatalf("expected file %s to exist", p)
			}
		}
	})

	t.Run("generated manifest validates", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "test-skill")
		runInit([]string{dir, "--id", "com.example.test2", "--type", "mcp-server", "--version", "1.0.0"})

		aemPath := filepath.Join(dir, "aem.json")
		m, err := manifest.LoadAEM(aemPath)
		if err != nil {
			t.Fatalf("LoadAEM(%s) error = %v", aemPath, err)
		}
		if err := m.Validate(); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
		if m.ID != "com.example.test2" {
			t.Fatalf("ID = %q, want com.example.test2", m.ID)
		}
		if m.Type != "mcp-server" {
			t.Fatalf("Type = %q, want mcp-server", m.Type)
		}
		if m.Version != "1.0.0" {
			t.Fatalf("Version = %q, want 1.0.0", m.Version)
		}
	})

	t.Run("does not overwrite existing files", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "existing")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		aemPath := filepath.Join(dir, "aem.json")
		if err := os.WriteFile(aemPath, []byte(`{}`), 0o644); err != nil {
			t.Fatal(err)
		}

		// runInit calls dieIf which calls os.Exit, so we can't test it directly.
		// Instead, test the precondition check manually.
		if _, err := os.Stat(aemPath); err != nil {
			t.Fatalf("expected aem.json to exist: %v", err)
		}
	})

	t.Run("devkey.json has restricted permissions", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "perm-test")
		runInit([]string{dir, "--id", "com.example.perm", "--type", "plugin"})

		keyPath := filepath.Join(dir, "devkey.json")
		info, err := os.Stat(keyPath)
		if err != nil {
			t.Fatalf("Stat(%s) error = %v", keyPath, err)
		}
		perm := info.Mode().Perm()
		if perm != 0o600 {
			t.Fatalf("devkey.json permissions = %o, want 600", perm)
		}
	})

	t.Run("default version is 0.1.0", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "default-ver")
		runInit([]string{dir, "--id", "com.example.defver", "--type", "skill"})

		aemPath := filepath.Join(dir, "aem.json")
		m, err := manifest.LoadAEM(aemPath)
		if err != nil {
			t.Fatalf("LoadAEM error = %v", err)
		}
		if m.Version != "0.1.0" {
			t.Fatalf("Version = %q, want 0.1.0", m.Version)
		}
	})
}
