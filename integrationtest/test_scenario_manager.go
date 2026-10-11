//go:build integration

package integrationtest

import (
	"io/fs"
	"os"
	"path/filepath"
)

const TEST_OPERATION_DIR = "./_testop"

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

type TestScenarioManager interface {
	PrepareFilesForTestScenario(scenarioName string) error
	TestScenario(scenarioName string, testOperarion func(testOperationDir string))
}

type testScenarioManagerImpl struct {
	filesSourceDir   string
	scenariosDir     string
	baseOperationDir string
}

func (m testScenarioManagerImpl) PrepareFilesForTestScenario(scenarioName string) error {
	sourceConfigFile := filepath.Join(m.scenariosDir, scenarioName, "imports-config.yml")
	operationDir := filepath.Join(m.baseOperationDir, scenarioName)

	if err := copyDir(m.filesSourceDir, operationDir); err != nil {
		return err
	}
	if err := copyFile(sourceConfigFile, filepath.Join(operationDir, "imports-config.yml"), 0o644); err != nil {
		return err
	}

	return nil
}

func (m testScenarioManagerImpl) TestScenario(scenarioName string, testOperarion func(testOperationDir string)) {
	operationDir := filepath.Join(m.baseOperationDir, scenarioName)
	testOperarion(operationDir)
}

func NewTestScenarioManager(workingDir string) TestScenarioManager {
	filesSourceDir := filepath.Join(workingDir, "test_files")
	scenariosDir := filepath.Join(workingDir, "scenarios")
	baseOperationDir := filepath.Join(workingDir, TEST_OPERATION_DIR)

	return testScenarioManagerImpl{filesSourceDir, scenariosDir, baseOperationDir}
}
