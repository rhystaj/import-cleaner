package config_test

import (
	types "importcleaner/internal"
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
				IsTypeImport:         "true",
			},
			{
				ModulePaths: []string{
					"pydantic.types",
					"argparse",
				},
				IsTypeImport: "false",
			},
		},
	}

	expectedOutput := config.Config{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []config.ImportStatementGroupRuleDefinition{
			{
				DependencySourceType: config.DependencySourceTypeNone,
				ModulePaths:          [][]string{},
				IsTypeImport:         types.OBNil,
			},
			{
				DependencySourceType: config.DependencySourceTypeInternal,
				ModulePaths: [][]string{
					{"somecoolmodule"},
					{"helpers", "database"},
				},
				IsTypeImport: types.OBNil,
			},
			{
				DependencySourceType: config.DependencySourceTypeExternal,
				ModulePaths:          [][]string{},
				IsTypeImport:         types.OBTrue,
			},
			{
				DependencySourceType: config.DependencySourceTypeNone,
				ModulePaths: [][]string{
					{"pydantic", "types"},
					{"argparse"},
				},
				IsTypeImport: types.OBFalse,
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
				IsTypeImport:         "notvalid",
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
				Path:             "groupingRules.isTypeImport",
				ErrorDescription: "Must be 'true', or 'false', but was 'notvalid'",
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
