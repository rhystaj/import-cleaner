package typescript_test

import (
	types "importcleaner/internal"
	"importcleaner/internal/config"
	"importcleaner/internal/languageprocessing/typescript"
	"testing"

	"github.com/stretchr/testify/assert"
)

type appliesToCase struct {
	description    string
	statement      typescript.TypescriptImportStatement
	ruleDefinition config.ImportStatementGroupRuleDefinition
	applies        bool
}

func TestAppliesTo(t *testing.T) {
	cases := []appliesToCase{
		{
			description: "Type import explicitly expected, statement is a type import",
			statement: typescript.TypescriptImportStatement{
				IsTypeImport: true,
			},
			ruleDefinition: config.ImportStatementGroupRuleDefinition{
				IsTypeImport: types.OBTrue,
			},
			applies: true,
		},
		{
			description: "Type import explicitly not expected, statement is not a type import",
			statement: typescript.TypescriptImportStatement{
				IsTypeImport: false,
			},
			ruleDefinition: config.ImportStatementGroupRuleDefinition{
				IsTypeImport: types.OBFalse,
			},
			applies: true,
		},
		{
			description: "Type import expectation not specified, statement is a type import",
			statement: typescript.TypescriptImportStatement{
				IsTypeImport: true,
			},
			ruleDefinition: config.ImportStatementGroupRuleDefinition{
				IsTypeImport: types.OBNil,
			},
			applies: true,
		},
		{
			description: "Type import expectation not specified, statement is not a type import",
			statement: typescript.TypescriptImportStatement{
				IsTypeImport: false,
			},
			ruleDefinition: config.ImportStatementGroupRuleDefinition{
				IsTypeImport: types.OBNil,
			},
			applies: true,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			rule := typescript.NewTypescriptImportStatementFromDefinition(c.ruleDefinition)

			applies := rule.AppliesTo(c.statement)
			assert.Equal(t, c.applies, applies)
		})
	}
}
