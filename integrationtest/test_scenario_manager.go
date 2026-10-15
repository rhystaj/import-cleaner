//go:build integration

package integrationtest

import (
	"os"
	"path/filepath"
)

const TEST_OPERATION_DIR = "./_testop"

type testOperation = func(testOperationDir string, expectedResultDir string)

type TestScenarioManager interface {
	TestScenario(scenarioName string, testOperation testOperation) error
}

type testScenarioManagerImpl struct {
	filesSourceDir   string
	scenariosDir     string
	baseOperationDir string
}

func (m testScenarioManagerImpl) TestScenario(scenarioName string, testOperation testOperation) error {
	sourceConfigFile := filepath.Join(m.scenariosDir, scenarioName, "imports-config.yml")
	operationDir := filepath.Join(m.baseOperationDir, scenarioName)

	if err := copyDir(m.filesSourceDir, operationDir); err != nil {
		return err
	}
	if err := copyFile(sourceConfigFile, filepath.Join(operationDir, "imports-config.yml"), 0o644); err != nil {
		return err
	}

	expectedResultDir := filepath.Join(m.scenariosDir, scenarioName, "expected_result")
	testOperation(operationDir, expectedResultDir)

	if err := os.RemoveAll(operationDir); err != nil {
		return err
	}

	return nil
}

func NewTestScenarioManager(workingDir string) TestScenarioManager {
	filesSourceDir := filepath.Join(workingDir, "test_files")
	scenariosDir := filepath.Join(workingDir, "scenarios")
	baseOperationDir := filepath.Join(workingDir, TEST_OPERATION_DIR)

	return testScenarioManagerImpl{filesSourceDir, scenariosDir, baseOperationDir}
}
