package util

import (
    "archive/zip"
    "io"
    "os"
    "path/filepath"
    "strings"
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

    if err := os.MkdirAll(destDir, 0o755); err != nil {
        return err
    }

    for _, f := range zr.File {
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
        defer rc.Close()

        out, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
        if err != nil {
            rc.Close()
            return err
        }
        if _, err := io.Copy(out, rc); err != nil {
            out.Close()
            rc.Close()
            return err
        }
        out.Close()
        rc.Close()
    }
    return nil
}
