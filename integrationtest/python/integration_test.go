//go:build integration

package main_test

import (
	"fmt"
	"importcleaner/integrationtest"
	"importcleaner/run"
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

func TestPythonIntegration(t *testing.T) {
	wd, _ := os.Getwd()
	testScenarioManager := integrationtest.NewTestScenarioManager(wd)

	testScenarioManager.PrepareFilesForTestScenario("grouping")
	testScenarioManager.TestScenario("grouping", func(testOperationDir string, expectedResultDir string) {
		runArgs := run.RunArgs{
			WorkingDir:    testOperationDir,
			TargetRootDir: ".",
		}
		runError := run.Run(runArgs)

		if runError != nil {
			t.Fatalf("Failed to run import cleaner: %v", runError)
		}

		integrationtest.AssertDirectoriesEqual(t, expectedResultDir, testOperationDir, []string{".py"})
	})
}
