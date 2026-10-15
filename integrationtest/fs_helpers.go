//go:build integration

package integrationtest

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		return copyFile(path, target, info.Mode().Perm())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	contents, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, contents, mode)
}

type treeEntry struct {
	relPath string
	isDir   bool
}

func collectTree(root string) ([]treeEntry, error) {
	var entries []treeEntry

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		entries = append(entries, treeEntry{relPath: rel, isDir: entry.IsDir()})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func treeSignature(entries []treeEntry) []string {
	signature := make([]string, len(entries))
	for i, entry := range entries {
		kind := "file"
		if entry.isDir {
			kind = "dir"
		}
		signature[i] = kind + " " + entry.relPath
	}

	sort.Strings(signature)
	return signature
}

func filterBySuffixes(entries []treeEntry, suffixes []string) []treeEntry {
	filtered := make([]treeEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.isDir || !hasAnySuffix(entry.relPath, suffixes) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func hasAnySuffix(path string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}
