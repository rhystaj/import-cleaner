package typescript

import (
	types "importcleaner/internal"
	"importcleaner/internal/config"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/organisation"
)

type TypescriptRules struct {
	IsTypeImport types.OptionalBool
}

type TypescriptImportStatementGroupRule struct {
	base       organisation.ImportStatementGroupRule[lp.ImportStatement]
	typescript TypescriptRules
}

func (r TypescriptImportStatementGroupRule) AppliesTo(statement TypescriptImportStatement) bool {
	if !r.base.AppliesTo(statement) {
		return false
	}

	if r.typescript.IsTypeImport != types.OBNil {
		return r.typescript.IsTypeImport.EqualsBool(statement.IsTypeImport)
	}

	return true
}

func NewTypeScriptImportStatementGroupRule(
	base organisation.ImportStatementGroupRule[lp.ImportStatement],
	typescript TypescriptRules) organisation.ImportStatementGroupRule[TypescriptImportStatement] {
	return TypescriptImportStatementGroupRule{
		base:       base,
		typescript: typescript,
	}
}

func NewTypescriptImportStatementFromDefinition(definition config.ImportStatementGroupRuleDefinition) organisation.ImportStatementGroupRule[TypescriptImportStatement] {
	return TypescriptImportStatementGroupRule{
		base: organisation.NewGenericGroupingRuleFromDefinition(definition),
		typescript: TypescriptRules{
			IsTypeImport: definition.IsTypeImport,
		},
	}
}
