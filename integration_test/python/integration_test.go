//go:build integration

package integration_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const TEST_OPERATION_DIR = "./_testop"

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()

	code := m.Run()

	if err := os.RemoveAll(filepath.Join(wd, TEST_OPERATION_DIR)); err != nil {
		fmt.Println("Failed to remove test operation directory:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

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

func prepareFilesForTestScenario(scenarioName string) error {
	wd, _ := os.Getwd()

	sourceDir := filepath.Join(wd, "test_files")
	sourceConfigFile := filepath.Join(wd, "scenarios", scenarioName, "imports-config.yml")
	operationDir := filepath.Join(wd, TEST_OPERATION_DIR, scenarioName)

	if err := copyDir(sourceDir, operationDir); err != nil {
		return err
	}
	if err := copyFile(sourceConfigFile, filepath.Join(operationDir, "imports-config.yml"), 0o644); err != nil {
		return err
	}

	return nil
}

func TestTest(t *testing.T) {
	if err := prepareFilesForTestScenario("grouping"); err != nil {
		t.Fatalf("Failed to prepare files for test scenario: %v", err)
	}
}
