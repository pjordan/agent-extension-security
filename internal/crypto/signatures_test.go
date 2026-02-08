package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSignAndVerifyDigest(t *testing.T) {
	key, err := GenerateDevKeypair()
	if err != nil {
		t.Fatalf("GenerateDevKeypair() error = %v", err)
	}
	pub, err := key.PublicKey()
	if err != nil {
		t.Fatalf("PublicKey() error = %v", err)
	}
	priv, err := key.PrivateKey()
	if err != nil {
		t.Fatalf("PrivateKey() error = %v", err)
	}

	const digest = "sha256:abcdef"
	sig, err := SignDigest(digest, priv, pub)
	if err != nil {
		t.Fatalf("SignDigest() error = %v", err)
	}

	if sig.Alg != "ed25519" {
		t.Fatalf("sig alg = %q, want ed25519", sig.Alg)
	}
	if sig.Digest != digest {
		t.Fatalf("sig digest = %q, want %q", sig.Digest, digest)
	}
	if sig.CreatedAt == "" {
		t.Fatal("sig created_at empty")
	}

	if err := Verify(sig, digest, pub); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestVerifyRejectsTampering(t *testing.T) {
	key, err := GenerateDevKeypair()
	if err != nil {
		t.Fatalf("GenerateDevKeypair() error = %v", err)
	}
	pub, err := key.PublicKey()
	if err != nil {
		t.Fatalf("PublicKey() error = %v", err)
	}
	priv, err := key.PrivateKey()
	if err != nil {
		t.Fatalf("PrivateKey() error = %v", err)
	}

	sig, err := SignDigest("sha256:abcdef", priv, pub)
	if err != nil {
		t.Fatalf("SignDigest() error = %v", err)
	}

	t.Run("digest mismatch", func(t *testing.T) {
		if err := Verify(sig, "sha256:zzz", pub); err == nil {
			t.Fatal("Verify() expected digest mismatch error, got nil")
		}
	})

	t.Run("unsupported algorithm", func(t *testing.T) {
		mod := *sig
		mod.Alg = "rsa"
		if err := Verify(&mod, mod.Digest, pub); err == nil {
			t.Fatal("Verify() expected algorithm error, got nil")
		}
	})

	t.Run("invalid signature payload", func(t *testing.T) {
		mod := *sig
		mod.Sig = "!!!!"
		if err := Verify(&mod, mod.Digest, pub); err == nil {
			t.Fatal("Verify() expected signature decode error, got nil")
		}
	})
}

func TestSignatureSaveLoadRoundTrip(t *testing.T) {
	key, err := GenerateDevKeypair()
	if err != nil {
		t.Fatalf("GenerateDevKeypair() error = %v", err)
	}
	pub, err := key.PublicKey()
	if err != nil {
		t.Fatalf("PublicKey() error = %v", err)
	}
	priv, err := key.PrivateKey()
	if err != nil {
		t.Fatalf("PrivateKey() error = %v", err)
	}

	sig, err := SignDigest("sha256:abcdef", priv, pub)
	if err != nil {
		t.Fatalf("SignDigest() error = %v", err)
	}

	p := filepath.Join(t.TempDir(), "sig.json")
	if err := SaveSignature(p, sig); err != nil {
		t.Fatalf("SaveSignature() error = %v", err)
	}
	loaded, err := LoadSignature(p)
	if err != nil {
		t.Fatalf("LoadSignature() error = %v", err)
	}

	if loaded.Alg != sig.Alg || loaded.Digest != sig.Digest || loaded.Sig != sig.Sig || loaded.PublicKey != sig.PublicKey {
		t.Fatal("loaded signature does not match saved signature")
	}
}

func TestLoadSignatureRejectsInvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sig.json")
	if err := os.WriteFile(p, []byte("{not-json"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if _, err := LoadSignature(p); err == nil {
		t.Fatal("LoadSignature() expected JSON parse error, got nil")
	}
}
