package typescript

import (
	types "importcleaner/internal"
	"importcleaner/internal/config"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/organisation"
)

type TypescriptImportStatementGroupRule struct {
	base         organisation.ImportStatementGroupRule[lp.ImportStatement]
	isTypeImport types.OptionalBool
}

func (r TypescriptImportStatementGroupRule) AppliesTo(statement TypescriptImportStatement) bool {
	if !r.base.AppliesTo(statement) {
		return false
	}

	if r.isTypeImport != types.OBNil {
		return r.isTypeImport.EqualsBool(statement.IsTypeImport)
	}

	return true
}

func NewTypescriptImportStatementFromDefinition(definition config.ImportStatementGroupRuleDefinition) organisation.ImportStatementGroupRule[TypescriptImportStatement] {
	return TypescriptImportStatementGroupRule{
		base:         organisation.NewGenericGroupingRuleFromDefinition(definition),
		isTypeImport: definition.IsTypeImport,
	}
}
