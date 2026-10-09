package typescript_test

import (
	types "importcleaner/internal"
	"importcleaner/internal/config"
	"importcleaner/internal/languageprocessing/typescript"
	"importcleaner/internal/organisation"
	"testing"

	"github.com/stretchr/testify/assert"
)

type typscriptImportStatementDetails struct {
	IsTypeImport bool
}

type typescriptRulesCase struct {
	description            string
	importStatementDetails typscriptImportStatementDetails
	typescriptRules        typescript.TypescriptRules
}

func buildTrueTypescriptRulesCases() []typescriptRulesCase {
	return []typescriptRulesCase{
		{
			description: "Type import explicitly expected, statement is a type import",
			importStatementDetails: typscriptImportStatementDetails{
				IsTypeImport: true,
			},
			typescriptRules: typescript.TypescriptRules{
				IsTypeImport: types.OBTrue,
			},
		},
		{
			description: "Type import explicitly not expected, statement is not a type import",
			importStatementDetails: typscriptImportStatementDetails{
				IsTypeImport: false,
			},
			typescriptRules: typescript.TypescriptRules{
				IsTypeImport: types.OBFalse,
			},
		},
		{
			description: "Type import expectation not specified, statement is a type import",
			importStatementDetails: typscriptImportStatementDetails{
				IsTypeImport: true,
			},
			typescriptRules: typescript.TypescriptRules{
				IsTypeImport: types.OBNil,
			},
		},
		{
			description: "Type import expectation not specified, statement is not a type import",
			importStatementDetails: typscriptImportStatementDetails{
				IsTypeImport: false,
			},
			typescriptRules: typescript.TypescriptRules{
				IsTypeImport: types.OBNil,
			},
		},
	}
}

func buildFalseTypescriptRulesCases() []typescriptRulesCase {
	return []typescriptRulesCase{
		{
			description: "Type import explicitly expected, statement is not a type import",
			importStatementDetails: typscriptImportStatementDetails{
				IsTypeImport: true,
			},
			typescriptRules: typescript.TypescriptRules{
				IsTypeImport: types.OBFalse,
			},
		},
		{
			description: "Type import explicitly not expected, statement is a type import",
			importStatementDetails: typscriptImportStatementDetails{
				IsTypeImport: false,
			},
			typescriptRules: typescript.TypescriptRules{
				IsTypeImport: types.OBTrue,
			},
		},
	}
}

func Test_AppliesTo_BaseMatches_TypescriptRulesMatch(t *testing.T) {
	baseRule := organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{})

	for _, c := range buildTrueTypescriptRulesCases() {
		testRule := typescript.NewTypeScriptImportStatementGroupRule(baseRule, c.typescriptRules)
		statement := typescript.TypescriptImportStatement{
			IsTypeImport: c.importStatementDetails.IsTypeImport,
		}

		t.Run(c.description, func(t *testing.T) {
			assert.True(t, testRule.AppliesTo(statement))
		})
	}
}

func Test_AppliesTo_BaseMatches_TypescriptRulesDoNotMatch(t *testing.T) {
	baseRule := organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{})

	for _, c := range buildFalseTypescriptRulesCases() {
		testRule := typescript.NewTypeScriptImportStatementGroupRule(baseRule, c.typescriptRules)
		statement := typescript.TypescriptImportStatement{
			IsTypeImport: c.importStatementDetails.IsTypeImport,
		}

		t.Run(c.description, func(t *testing.T) {
			assert.False(t, testRule.AppliesTo(statement))
		})
	}
}

func Test_AppliesTo_BaseDoesNotMatch(t *testing.T) {
	baseRuleDefinitions := []config.ImportStatementGroupRuleDefinition{
		{
			DependencySourceType: config.DependencySourceTypeInternal,
			ModulePaths:          [][]string{{"someothermodule"}},
		},
		{
			DependencySourceType: config.DependencySourceTypeExternal,
			ModulePaths:          [][]string{{"target"}},
		},
	}

	for _, c := range buildTrueTypescriptRulesCases() {
		t.Run(c.description, func(t *testing.T) {
			for _, baseRuleDefinition := range baseRuleDefinitions {
				baseRule := organisation.NewGenericGroupingRuleFromDefinition(baseRuleDefinition)
				testRule := typescript.NewTypeScriptImportStatementGroupRule(baseRule, c.typescriptRules)

				statement := typescript.TypescriptImportStatement{
					ModulePathParts: []string{"target"},
					External:        false,
					IsTypeImport:    c.importStatementDetails.IsTypeImport,
				}

				assert.False(t, testRule.AppliesTo(statement))
			}
		})
	}
}
