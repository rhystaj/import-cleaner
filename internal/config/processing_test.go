package config_test

import (
	"importcleaner/internal/config"
	"importcleaner/internal/organisation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessAndValidationConfig_ValidConfig(t *testing.T) {

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

	result, err := config.ProcessAndValidateConfig(testConfig)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, result)
}

func TestProcessAndValidationConfig_InvalidConfig(t *testing.T) {

	testIgnoredPaths := []string{}

	testConfig := config.RawConfig{
		IgnoredPaths: testIgnoredPaths,
		GroupingRules: []config.RawGroupingRule{
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
