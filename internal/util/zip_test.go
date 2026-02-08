package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestZipDirAndUnzipFileRoundTrip(t *testing.T) {
	src := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "root.txt"), []byte("root"), 0o644); err != nil {
		t.Fatalf("WriteFile(root) error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(src, "sub"), 0o755); err != nil {
		t.Fatalf("MkdirAll(sub) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "sub", "file.txt"), []byte("nested"), 0o644); err != nil {
		t.Fatalf("WriteFile(sub/file.txt) error = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "bundle.zip")
	if err := ZipDir(src, zipPath); err != nil {
		t.Fatalf("ZipDir() error = %v", err)
	}

	dest := filepath.Join(t.TempDir(), "out")
	if err := UnzipFile(zipPath, dest); err != nil {
		t.Fatalf("UnzipFile() error = %v", err)
	}

	b, err := os.ReadFile(filepath.Join(dest, "root.txt"))
	if err != nil {
		t.Fatalf("ReadFile(root.txt) error = %v", err)
	}
	if string(b) != "root" {
		t.Fatalf("root.txt = %q, want %q", string(b), "root")
	}
	b, err = os.ReadFile(filepath.Join(dest, "sub", "file.txt"))
	if err != nil {
		t.Fatalf("ReadFile(sub/file.txt) error = %v", err)
	}
	if string(b) != "nested" {
		t.Fatalf("sub/file.txt = %q, want %q", string(b), "nested")
	}
}

func TestZipDirExcludesJunk(t *testing.T) {
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, ".git"), 0o755); err != nil {
		t.Fatalf("MkdirAll(.git) error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(src, "__MACOSX"), 0o755); err != nil {
		t.Fatalf("MkdirAll(__MACOSX) error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(src, ".git", "config"), []byte("secret"), 0o644); err != nil {
		t.Fatalf("WriteFile(.git/config) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, ".DS_Store"), []byte("junk"), 0o644); err != nil {
		t.Fatalf("WriteFile(.DS_Store) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "__MACOSX", "meta"), []byte("junk"), 0o644); err != nil {
		t.Fatalf("WriteFile(__MACOSX/meta) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "good.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatalf("WriteFile(good.txt) error = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "bundle.zip")
	if err := ZipDir(src, zipPath); err != nil {
		t.Fatalf("ZipDir() error = %v", err)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("zip.OpenReader() error = %v", err)
	}
	defer r.Close()

	seen := map[string]bool{}
	for _, f := range r.File {
		seen[f.Name] = true
	}

	if !seen["good.txt"] {
		t.Fatal("expected good.txt in archive")
	}
	if seen[".git/config"] {
		t.Fatal("did not expect .git/config in archive")
	}
	if seen[".DS_Store"] {
		t.Fatal("did not expect .DS_Store in archive")
	}
	if seen["__MACOSX/meta"] {
		t.Fatal("did not expect __MACOSX/meta in archive")
	}
}

func TestUnzipFileBlocksZipSlip(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "slip.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("os.Create() error = %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("../evil.txt")
	if err != nil {
		t.Fatalf("zw.Create() error = %v", err)
	}
	if _, err := w.Write([]byte("owned")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zw.Close() error = %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("f.Close() error = %v", err)
	}

	dest := filepath.Join(t.TempDir(), "out")
	if err := UnzipFile(zipPath, dest); err == nil {
		t.Fatal("UnzipFile() expected ZipSlip protection error, got nil")
	}
}

func TestZipDirRejectsSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior differs on windows")
	}

	src := t.TempDir()
	target := filepath.Join(src, "target.txt")
	if err := os.WriteFile(target, []byte("target"), 0o644); err != nil {
		t.Fatalf("WriteFile(target) error = %v", err)
	}
	if err := os.Symlink(target, filepath.Join(src, "link.txt")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "bundle.zip")
	if err := ZipDir(src, zipPath); err == nil || !strings.Contains(err.Error(), "refusing to package symlink") {
		t.Fatalf("ZipDir() error = %v, want symlink rejection", err)
	}
}

func TestUnzipFileBlocksSuspiciousCompressionRatio(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "ratio.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("os.Create() error = %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("huge.txt")
	if err != nil {
		t.Fatalf("zw.Create() error = %v", err)
	}
	if _, err := w.Write(bytes.Repeat([]byte("0"), 2<<20)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zw.Close() error = %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("f.Close() error = %v", err)
	}

	dest := filepath.Join(t.TempDir(), "out")
	if err := UnzipFile(zipPath, dest); err == nil || !strings.Contains(err.Error(), "suspicious compression ratio") {
		t.Fatalf("UnzipFile() error = %v, want suspicious compression ratio", err)
	}
}

func TestUnzipFileRejectsSymlinkEntries(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "symlink.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("os.Create() error = %v", err)
	}
	zw := zip.NewWriter(f)
	h := &zip.FileHeader{Name: "link"}
	h.SetMode(os.ModeSymlink | 0o777)
	w, err := zw.CreateHeader(h)
	if err != nil {
		t.Fatalf("CreateHeader() error = %v", err)
	}
	if _, err := w.Write([]byte("target.txt")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zw.Close() error = %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("f.Close() error = %v", err)
	}

	dest := filepath.Join(t.TempDir(), "out")
	if err := UnzipFile(zipPath, dest); err == nil || !strings.Contains(err.Error(), "refusing to extract symlink") {
		t.Fatalf("UnzipFile() error = %v, want symlink rejection", err)
	}
}

func TestUnzipFileRejectsTooManyEntries(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "many.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("os.Create() error = %v", err)
	}
	zw := zip.NewWriter(f)
	for i := 0; i < maxArchiveEntries+1; i++ {
		name := fmt.Sprintf("f-%05d.txt", i)
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("zw.Create(%q) error = %v", name, err)
		}
		if _, err := w.Write([]byte("x")); err != nil {
			t.Fatalf("Write(%q) error = %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zw.Close() error = %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("f.Close() error = %v", err)
	}

	dest := filepath.Join(t.TempDir(), "out")
	err = UnzipFile(zipPath, dest)
	if err == nil || !strings.Contains(err.Error(), "too many entries") {
		t.Fatalf("UnzipFile() error = %v, want too many entries", err)
	}
}
