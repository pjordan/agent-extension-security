package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pjordan/agent-extension-security/internal/manifest"
)

func boolPtr(v bool) *bool {
	return &v
}

func TestEvaluate(t *testing.T) {
	a := &manifest.AEM{
		Permissions: manifest.Permissions{
			Network: manifest.NetPerms{
				AllowIPLiterals: true,
			},
			Process: manifest.ProcessPerm{
				AllowShell: true,
			},
		},
	}

	t.Run("findings when denied values match", func(t *testing.T) {
		var p Policy
		p.Permissions.Deny.Network.AllowIPLiterals = boolPtr(true)
		p.Permissions.Deny.Process.AllowShell = boolPtr(true)

		findings := Evaluate(&p, a)
		if len(findings) != 2 {
			t.Fatalf("len(findings) = %d, want 2", len(findings))
		}
	})

	t.Run("no findings when deny values do not match", func(t *testing.T) {
		var p Policy
		p.Permissions.Deny.Network.AllowIPLiterals = boolPtr(false)
		p.Permissions.Deny.Process.AllowShell = boolPtr(false)

		findings := Evaluate(&p, a)
		if len(findings) != 0 {
			t.Fatalf("len(findings) = %d, want 0", len(findings))
		}
	})
}

func TestLoadPolicy(t *testing.T) {
	t.Run("defaults mode to enforce", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(p, []byte(`{"permissions":{}}`), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		got, err := Load(p)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if got.Mode != "enforce" {
			t.Fatalf("Mode = %q, want enforce", got.Mode)
		}
	})

	t.Run("respects explicit mode", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(p, []byte(`{"mode":"warn"}`), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		got, err := Load(p)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if got.Mode != "warn" {
			t.Fatalf("Mode = %q, want warn", got.Mode)
		}
	})

	t.Run("rejects invalid json", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(p, []byte("{nope"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		if _, err := Load(p); err == nil {
			t.Fatal("Load() expected error, got nil")
		}
	})

	t.Run("rejects invalid mode", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(p, []byte(`{"mode":"monitor"}`), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		if _, err := Load(p); err == nil {
			t.Fatal("Load() expected invalid mode error, got nil")
		}
	})

	t.Run("rejects unknown fields", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(p, []byte(`{"mode":"warn","unknown":true}`), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		if _, err := Load(p); err == nil {
			t.Fatal("Load() expected unknown field error, got nil")
		}
	})

	t.Run("rejects trailing json values", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(p, []byte(`{"mode":"warn"}{"extra":true}`), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		if _, err := Load(p); err == nil {
			t.Fatal("Load() expected trailing-json error, got nil")
		}
	})
}
