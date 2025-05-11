package config_test

import (
	"importcleaner/internal/config"
	"importcleaner/internal/organisation"
	"reflect"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {

	testIgnoredPaths := []string{}

	testConfig := config.RawConfig{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []config.RawGroupingRule{
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

	result, _ := config.ProcessAndValidateConfig(testConfig)

	if !reflect.DeepEqual(expectedOutput, result) {
		t.Errorf("Expected %+v, but recieved %+v.", expectedOutput, result)
	}
}
