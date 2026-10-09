package typescript

import (
	"importcleaner/internal/config"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/organisation"
)

type TypescriptImportStatementGroupRule struct {
	base         organisation.ImportStatementGroupRule[lp.ImportStatement]
	isTypeImport bool
}

func (r TypescriptImportStatementGroupRule) AppliesTo(statement TypescriptImportStatement) bool {
	return r.base.AppliesTo(statement) && r.isTypeImport == statement.IsTypeImport
}

func NewTypescriptImportStatementFromDefinition(definition config.ImportStatementGroupRuleDefinition) organisation.ImportStatementGroupRule[TypescriptImportStatement] {
	return TypescriptImportStatementGroupRule{
		base:         organisation.NewGenericGroupingRuleFromDefinition(definition),
		isTypeImport: definition.IsTypeImport,
	}
}
