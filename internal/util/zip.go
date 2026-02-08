package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxArchiveEntries         = 10000
	maxArchiveBytes     int64 = 512 << 20 // 512 MiB
	maxEntryBytes       int64 = 64 << 20  // 64 MiB
	maxCompressionRatio       = 200
)

// ZipDir zips a directory into outPath. It preserves relative paths.
// It excludes common junk like .git and OS metadata.
func ZipDir(srcDir, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zw := zip.NewWriter(outFile)
	defer zw.Close()

	srcDir = filepath.Clean(srcDir)
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if rel == "." {
			return nil
		}
		// Excludes
		if strings.HasPrefix(rel, ".git/") || rel == ".git" {
			return nil
		}
		if strings.HasPrefix(rel, ".DS_Store") || strings.Contains(rel, "__MACOSX") {
			return nil
		}

		if info.IsDir() {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to package symlink: %s", rel)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("refusing to package non-regular file: %s", rel)
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		w, err := zw.Create(rel)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		return err
	})
}

func UnzipFile(zipPath, destDir string) error {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zr.Close()
	if len(zr.File) > maxArchiveEntries {
		return fmt.Errorf("archive contains too many entries: %d", len(zr.File))
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	var totalWritten int64
	for _, f := range zr.File {
		if f.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to extract symlink: %s", f.Name)
		}
		if !f.FileInfo().IsDir() {
			if f.UncompressedSize64 > uint64(maxEntryBytes) {
				return fmt.Errorf("archive entry too large: %s", f.Name)
			}
			if f.CompressedSize64 > 0 && f.UncompressedSize64 > f.CompressedSize64*maxCompressionRatio {
				return fmt.Errorf("suspicious compression ratio: %s", f.Name)
			}
		}

		outPath := filepath.Join(destDir, filepath.FromSlash(f.Name))
		// prevent ZipSlip
		if !strings.HasPrefix(filepath.Clean(outPath), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return os.ErrPermission
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		out, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			rc.Close()
			return err
		}

		written, err := io.Copy(out, io.LimitReader(rc, maxEntryBytes+1))
		if err != nil {
			out.Close()
			rc.Close()
			return err
		}
		if written > maxEntryBytes {
			out.Close()
			rc.Close()
			return fmt.Errorf("archive entry exceeded size limit: %s", f.Name)
		}
		if f.UncompressedSize64 > 0 && written != int64(f.UncompressedSize64) {
			out.Close()
			rc.Close()
			return fmt.Errorf("archive entry size mismatch: %s", f.Name)
		}
		totalWritten += written
		if totalWritten > maxArchiveBytes {
			out.Close()
			rc.Close()
			return fmt.Errorf("archive exceeds total uncompressed size limit")
		}

		out.Close()
		rc.Close()
	}
	return nil
}
