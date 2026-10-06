package config_test

import (
	"importcleaner/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessAndValidationConfig_ValidConfig(t *testing.T) {

	testIgnoredPaths := []string{}

	testConfig := config.RawConfig{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []config.RawGroupingRuleDefinition{
			{},
			{
				DependencySourceType: "internal",
				ModulePaths: []string{
					"somecoolmodule",
					"helpers.database",
				},
			},
			{
				DependencySourceType: "external",
			},
			{
				ModulePaths: []string{
					"pydantic.types",
					"argparse",
				},
			},
		},
	}

	expectedOutput := config.Config{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []config.ImportStatementGroupRuleDefinition{
			{
				DependencySourceType: config.DependencySourceTypeNone,
				ModulePaths:          [][]string{},
			},
			{
				DependencySourceType: config.DependencySourceTypeInternal,
				ModulePaths: [][]string{
					{"somecoolmodule"},
					{"helpers", "database"},
				},
			},
			{
				DependencySourceType: config.DependencySourceTypeExternal,
				ModulePaths:          [][]string{},
			},
			{
				DependencySourceType: config.DependencySourceTypeNone,
				ModulePaths: [][]string{
					{"pydantic", "types"},
					{"argparse"},
				},
			},
		},
	}

	result, err := config.ProcessAndValidateConfig(testConfig)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, result)
}

func TestProcessAndValidationConfig_InvalidConfig(t *testing.T) {

	testIgnoredPaths := []string{}

	testConfig := config.RawConfig{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []config.RawGroupingRuleDefinition{
			{},
			{
				DependencySourceType: "notvalid",
			},
			{
				DependencySourceType: "external",
			},
			{
				DependencySourceType: "notvalideither",
			},
		},
	}

	expectedError := config.ConfigValidationError{
		InvalidPaths: []config.InvalidPath{
			{
				Path:             "groupingRules.dependencySourceType",
				ErrorDescription: "Must be 'internal', or 'external', but was 'notvalid'",
			},
			{
				Path:             "groupingRules.dependencySourceType",
				ErrorDescription: "Must be 'internal', or 'external', but was 'notvalideither'",
			},
		},
	}

	_, err := config.ProcessAndValidateConfig(testConfig)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, *err)
}
