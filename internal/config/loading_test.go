package config_test

import (
	"importcleaner/internal/config"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/organisation"
	"reflect"
	"testing"

	"github.com/go-yaml/yaml"
)

func TestLoadValidConfig(t *testing.T) {

	testIgnoredPaths := []string{}

	testYaml := config.ConfigYaml{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []config.GroupingRuleYaml{
			{},
			{
				DependencySourceType: "internal",
			},
			{
				DependencySourceType: "external",
			},
		},
	}

	expectedOutput := config.Config{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []organisation.ImportStatementGroupRule{
			{
				DependencySourceType: organisation.DependencySourceTypeNone,
			},
			{
				DependencySourceType: organisation.DependencySourceTypeInternal,
			},
			{
				DependencySourceType: organisation.DependencySourceTypeExternal,
			},
		},
	}

	yamlBytes, _ := yaml.Marshal(testYaml)
	yamlString := string(yamlBytes)

	testFileManager := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile("testFile.yml", yamlString),
	})

	result, _ := config.LoadAndValidateConfigFromFile(testFileManager, "testFile.yml")

	if !reflect.DeepEqual(expectedOutput, result) {
		t.Errorf("Expected %+v, but recieved %+v.", expectedOutput, result)
	}
}
