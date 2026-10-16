//go:build integration

package main_test

import (
	"importcleaner/integrationtest"
	"importcleaner/run"
	"os"
	"testing"
)

func TestTypescriptIntegration_Grouping(t *testing.T) {
	wd, _ := os.Getwd()
	testScenarioManager := integrationtest.NewTestScenarioManagerForWorkingDir(wd)

	testScenarioManager.TestScenario("grouping", func(testOperationDir string, expectedResultDir string) {
		runArgs := run.RunArgs{
			WorkingDir:    testOperationDir,
			TargetRootDir: ".",
		}
		runError := run.Run(runArgs)

		if runError != nil {
			t.Errorf("Failed to run import cleaner: %v", runError)
		}

		integrationtest.AssertDirectoriesEqual(t, expectedResultDir, testOperationDir, []string{".ts"})
	})
}

func TestTypescriptIntegration_Stability(t *testing.T) {
	wd, _ := os.Getwd()
	testScenarioManager := integrationtest.NewTestScenarioManagerForWorkingDir(wd)

	testScenarioManager.TestScenario("stability", func(testOperationDir string, expectedResultDir string) {
		runArgs := run.RunArgs{
			WorkingDir:    testOperationDir,
			TargetRootDir: ".",
		}
		runError := run.Run(runArgs)

		if runError != nil {
			t.Errorf("Failed to run import cleaner: %v", runError)
		}

		integrationtest.AssertDirectoriesEqual(t, expectedResultDir, testOperationDir, []string{".ts"})
	})
}
