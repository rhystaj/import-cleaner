package config_test

import (
	types "importcleaner/internal"
	"importcleaner/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessAndValidationConfig_ValidConfig(t *testing.T) {

	testIgnoredPaths := []string{}

	testFormattingCases := []struct {
		description             string
		rawFormattingRules      config.RawFormattingRules
		expectedFormattingRules config.FormattingRules
	}{
		{
			description: "Default space indent",
			rawFormattingRules: config.RawFormattingRules{
				ConsolidateImports:  true,
				LineBreakAtNumChars: 66,
			},
			expectedFormattingRules: config.FormattingRules{
				ConsolidateImports:  true,
				LineBreakAtNumChars: 66,
				Indent:              "    ",
			},
		},
		{
			description: "Space indent",
			rawFormattingRules: config.RawFormattingRules{
				ConsolidateImports:  true,
				LineBreakAtNumChars: 66,
				Indent: config.RawIndentSettings{
					Type: "space",
					Size: 3,
				},
			},
			expectedFormattingRules: config.FormattingRules{
				ConsolidateImports:  true,
				LineBreakAtNumChars: 66,
				Indent:              "   ",
			},
		},
		{
			description: "Tab indent",
			rawFormattingRules: config.RawFormattingRules{
				ConsolidateImports:  true,
				LineBreakAtNumChars: 66,
				Indent: config.RawIndentSettings{
					Type: "tab",
					Size: 1,
				},
			},
			expectedFormattingRules: config.FormattingRules{
				ConsolidateImports:  true,
				LineBreakAtNumChars: 66,
				Indent:              "\t",
			},
		},
	}

	for _, testFormattingCase := range testFormattingCases {
		t.Run("Formatting: "+testFormattingCase.description, func(t *testing.T) {
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
				FormattingRules: testFormattingCase.rawFormattingRules,
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
				FormattingRules: testFormattingCase.expectedFormattingRules,
			}

			result, err := config.ProcessAndValidateConfig(testConfig)

			assert.Nil(t, err)
			assert.Equal(t, expectedOutput, result)
		})
	}
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
		FormattingRules: config.RawFormattingRules{
			ConsolidateImports:  true,
			LineBreakAtNumChars: 66,
			Indent: config.RawIndentSettings{
				Type: "notvalid",
				Size: -3,
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
			{
				Path:             "formattingRules.indent.type",
				ErrorDescription: "Must be 'space', or 'tab', but was 'notvalid'",
			},
			{
				Path:             "formattingRules.indent.size",
				ErrorDescription: "Must be greater than 0",
			},
		},
	}

	_, err := config.ProcessAndValidateConfig(testConfig)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, *err)
}
