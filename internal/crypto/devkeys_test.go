package crypto

import (
	"crypto/ed25519"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateDevKeypairAndDecode(t *testing.T) {
	k, err := GenerateDevKeypair()
	if err != nil {
		t.Fatalf("GenerateDevKeypair() error = %v", err)
	}
	if k.Type != "ed25519" {
		t.Fatalf("key type = %q, want ed25519", k.Type)
	}

	pub, err := k.PublicKey()
	if err != nil {
		t.Fatalf("PublicKey() error = %v", err)
	}
	if len(pub) != ed25519.PublicKeySize {
		t.Fatalf("public key len = %d, want %d", len(pub), ed25519.PublicKeySize)
	}

	priv, err := k.PrivateKey()
	if err != nil {
		t.Fatalf("PrivateKey() error = %v", err)
	}
	if len(priv) != ed25519.PrivateKeySize {
		t.Fatalf("private key len = %d, want %d", len(priv), ed25519.PrivateKeySize)
	}
}

func TestLoadDevKeyValidation(t *testing.T) {
	t.Run("unsupported type", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "key.json")
		content := `{"type":"rsa","public":"abc"}`
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := LoadDevKey(p); err == nil {
			t.Fatal("LoadDevKey() expected error, got nil")
		}
	})

	t.Run("missing public", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "key.json")
		content := `{"type":"ed25519"}`
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := LoadDevKey(p); err == nil {
			t.Fatal("LoadDevKey() expected error, got nil")
		}
	})
}

func TestSaveAndLoadDevKeyRoundTrip(t *testing.T) {
	k, err := GenerateDevKeypair()
	if err != nil {
		t.Fatalf("GenerateDevKeypair() error = %v", err)
	}

	p := filepath.Join(t.TempDir(), "devkey.json")
	if err := SaveKey(p, k); err != nil {
		t.Fatalf("SaveKey() error = %v", err)
	}

	st, err := os.Stat(p)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if got := st.Mode().Perm(); got != 0o600 {
		t.Fatalf("key mode = %o, want 600", got)
	}

	loaded, err := LoadDevKey(p)
	if err != nil {
		t.Fatalf("LoadDevKey() error = %v", err)
	}
	if loaded.Type != "ed25519" {
		t.Fatalf("loaded type = %q, want ed25519", loaded.Type)
	}
	if loaded.Public != k.Public {
		t.Fatalf("loaded public mismatch")
	}
	if loaded.Private != k.Private {
		t.Fatalf("loaded private mismatch")
	}
}
