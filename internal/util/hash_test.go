package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSha256File(t *testing.T) {
	p := filepath.Join(t.TempDir(), "data.txt")
	if err := os.WriteFile(p, []byte("hello world\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := Sha256File(p)
	if err != nil {
		t.Fatalf("Sha256File() error = %v", err)
	}

	const want = "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447"
	if got != want {
		t.Fatalf("Sha256File() = %s, want %s", got, want)
	}
}

func TestSha256FileMissing(t *testing.T) {
	if _, err := Sha256File(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("Sha256File() expected error, got nil")
	}
}
