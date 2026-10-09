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
	return r.base.AppliesTo(statement) && r.isTypeImport.EqualsBool(statement.IsTypeImport)
}

func NewTypescriptImportStatementFromDefinition(definition config.ImportStatementGroupRuleDefinition) organisation.ImportStatementGroupRule[TypescriptImportStatement] {
	return TypescriptImportStatementGroupRule{
		base:         organisation.NewGenericGroupingRuleFromDefinition(definition),
		isTypeImport: definition.IsTypeImport,
	}
}
