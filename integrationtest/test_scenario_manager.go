//go:build integration

package integrationtest

import (
	"os"
	"path/filepath"
)

const TEST_OPERATION_DIR = "_testop"
const SETUP_DIR = "setup"
const EXPECTED_RESULT_DIR = "expected_result"

type testOperation = func(testOperationDir string, expectedResultDir string)

type TestScenarioManager interface {
	TestScenario(scenarioName string, testOperation testOperation) error
}

type testScenarioManagerImpl struct {
	scenariosDir     string
	baseOperationDir string
}

func (m testScenarioManagerImpl) TestScenario(scenarioName string, testOperation testOperation) error {
	setupDir := filepath.Join(m.scenariosDir, scenarioName, SETUP_DIR)
	operationDir := filepath.Join(m.baseOperationDir, scenarioName)

	if err := copyDir(setupDir, operationDir); err != nil {
		return err
	}

	expectedResultDir := filepath.Join(m.scenariosDir, scenarioName, "expected_result")
	testOperation(operationDir, expectedResultDir)

	if err := os.RemoveAll(operationDir); err != nil {
		return err
	}

	return nil
}

func NewTestScenarioManagerForWorkingDir(workingDir string) TestScenarioManager {
	scenariosDir := filepath.Join(workingDir, "scenarios")
	baseOperationDir := filepath.Join(workingDir, TEST_OPERATION_DIR)

	return testScenarioManagerImpl{scenariosDir, baseOperationDir}
}
