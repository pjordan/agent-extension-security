package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pjordan/agent-extension-security/internal/crypto"
	"github.com/pjordan/agent-extension-security/internal/manifest"
	"github.com/pjordan/agent-extension-security/internal/util"
)

func TestVerifyArtifactRequiresTrustedKey(t *testing.T) {
	key, err := crypto.GenerateDevKeypair()
	if err != nil {
		t.Fatalf("GenerateDevKeypair() error = %v", err)
	}

	keyPath := filepath.Join(t.TempDir(), "devkey.json")
	if err := crypto.SaveKey(keyPath, key); err != nil {
		t.Fatalf("SaveKey() error = %v", err)
	}

	art := filepath.Join(t.TempDir(), "artifact.aext")
	if err := os.WriteFile(art, []byte("artifact"), 0o644); err != nil {
		t.Fatalf("WriteFile(artifact) error = %v", err)
	}

	sha, err := util.Sha256File(art)
	if err != nil {
		t.Fatalf("Sha256File() error = %v", err)
	}
	pub, err := key.PublicKey()
	if err != nil {
		t.Fatalf("PublicKey() error = %v", err)
	}
	priv, err := key.PrivateKey()
	if err != nil {
		t.Fatalf("PrivateKey() error = %v", err)
	}
	sig, err := crypto.SignDigest("sha256:"+sha, priv, pub)
	if err != nil {
		t.Fatalf("SignDigest() error = %v", err)
	}
	sigPath := filepath.Join(t.TempDir(), "sig.json")
	if err := crypto.SaveSignature(sigPath, sig); err != nil {
		t.Fatalf("SaveSignature() error = %v", err)
	}

	if err := verifyArtifact(art, sigPath, "", false); err == nil || !strings.Contains(err.Error(), "no trusted key provided") {
		t.Fatalf("verifyArtifact() error = %v, want no trusted key provided", err)
	}
	if err := verifyArtifact(art, sigPath, keyPath, false); err != nil {
		t.Fatalf("verifyArtifact() with trusted pub error = %v", err)
	}
	if err := verifyArtifact(art, sigPath, "", true); err != nil {
		t.Fatalf("verifyArtifact() with allow-embedded-key error = %v", err)
	}
}

func TestEvaluateInstallPolicy(t *testing.T) {
	aem := &manifest.AEM{
		Schema:  "aessf.dev/aem/v0",
		ID:      "com.example.poc",
		Type:    "skill",
		Version: "1.0.0",
		Permissions: manifest.Permissions{
			Process: manifest.ProcessPerm{
				AllowShell: true,
			},
		},
	}
	aemBytes, err := aem.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}
	aemPath := filepath.Join(t.TempDir(), "aem.json")
	if err := os.WriteFile(aemPath, aemBytes, 0o644); err != nil {
		t.Fatalf("WriteFile(aem) error = %v", err)
	}

	t.Run("enforce blocks denied permissions", func(t *testing.T) {
		policyPath := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(policyPath, []byte(`{
  "mode": "enforce",
  "permissions": {
    "deny": {
      "process": {
        "allow_shell": true
      }
    }
  }
}`), 0o644); err != nil {
			t.Fatalf("WriteFile(policy) error = %v", err)
		}

		if _, _, err := evaluateInstallPolicy(aemPath, policyPath, false); err == nil || !strings.Contains(err.Error(), "policy denied install") {
			t.Fatalf("evaluateInstallPolicy() error = %v, want policy denied install", err)
		}
	})

	t.Run("warn returns findings but allows install", func(t *testing.T) {
		policyPath := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(policyPath, []byte(`{
  "mode": "warn",
  "permissions": {
    "deny": {
      "process": {
        "allow_shell": true
      }
    }
  }
}`), 0o644); err != nil {
			t.Fatalf("WriteFile(policy) error = %v", err)
		}

		p, findings, err := evaluateInstallPolicy(aemPath, policyPath, false)
		if err != nil {
			t.Fatalf("evaluateInstallPolicy() error = %v", err)
		}
		if p.Mode != "warn" {
			t.Fatalf("Mode = %q, want warn", p.Mode)
		}
		if len(findings) != 1 {
			t.Fatalf("len(findings) = %d, want 1", len(findings))
		}
	})
}

func TestInstallDevMode(t *testing.T) {
	aem := &manifest.AEM{
		Schema:  "aessf.dev/aem/v0",
		ID:      "com.example.devtest",
		Type:    "skill",
		Version: "1.0.0",
		Permissions: manifest.Permissions{
			Process: manifest.ProcessPerm{
				AllowShell: true,
			},
		},
	}
	aemBytes, err := aem.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}
	aemPath := filepath.Join(t.TempDir(), "aem.json")
	if err := os.WriteFile(aemPath, aemBytes, 0o644); err != nil {
		t.Fatalf("WriteFile(aem) error = %v", err)
	}

	t.Run("dev mode succeeds without policy file", func(t *testing.T) {
		p, findings, err := evaluateInstallPolicy(aemPath, "", true)
		if err != nil {
			t.Fatalf("evaluateInstallPolicy() dev mode error = %v", err)
		}
		if p.Mode != "warn" {
			t.Fatalf("Mode = %q, want warn", p.Mode)
		}
		// permissive policy has no deny rules, so no findings
		if len(findings) != 0 {
			t.Fatalf("len(findings) = %d, want 0", len(findings))
		}
	})

	t.Run("dev mode with explicit policy uses provided policy", func(t *testing.T) {
		policyPath := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(policyPath, []byte(`{
  "mode": "enforce",
  "permissions": {
    "deny": {
      "process": {
        "allow_shell": true
      }
    }
  }
}`), 0o644); err != nil {
			t.Fatalf("WriteFile(policy) error = %v", err)
		}

		// Even in dev mode, if --policy is provided it should be used
		_, _, err := evaluateInstallPolicy(aemPath, policyPath, true)
		if err == nil || !strings.Contains(err.Error(), "policy denied install") {
			t.Fatalf("evaluateInstallPolicy() dev+policy error = %v, want policy denied install", err)
		}
	})

	t.Run("non-dev mode fails without policy", func(t *testing.T) {
		_, _, err := evaluateInstallPolicy(aemPath, "", false)
		if err == nil || !strings.Contains(err.Error(), "--policy is required") {
			t.Fatalf("evaluateInstallPolicy() no-dev error = %v, want --policy is required", err)
		}
	})
}
