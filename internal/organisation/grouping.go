package organisation

import (
	"importcleaner/internal/config"
	"importcleaner/internal/parsing"
)

type ImportStatementGroupRule interface {
	AppliesTo(statement parsing.ImportStatement) bool
}

type GenericImportStatementGroupRule struct {
	dependencySourceType config.DependencySourceType
	modulePaths          [][]string
}

func NewGroupingRuleFromDefinition(definition config.ImportStatementGroupRuleDefinition) ImportStatementGroupRule {
	return GenericImportStatementGroupRule{
		dependencySourceType: definition.DependencySourceType,
		modulePaths:          definition.ModulePaths,
	}
}

func (gr GenericImportStatementGroupRule) statementModulePathMatchesGroupRuleModulePath(statementModulePath []string, groupRuleModulePath []string) bool {
	for groupRuleModulePathIndex, groupRuleModulePathPart := range groupRuleModulePath {
		if groupRuleModulePathIndex > len(statementModulePath)-1 {
			return false
		}
		if groupRuleModulePathPart != statementModulePath[groupRuleModulePathIndex] {
			return false
		}
	}
	return true
}

func (gr GenericImportStatementGroupRule) AppliesTo(statement parsing.ImportStatement) bool {
	statmentDetails := statement.GetGenericDetails()

	var dependecySourceTypeMatches bool
	switch gr.dependencySourceType {
	case config.DependencySourceTypeNone:
		dependecySourceTypeMatches = true
	case config.DependencySourceTypeInternal:
		dependecySourceTypeMatches = !statmentDetails.External
	case config.DependencySourceTypeExternal:
		dependecySourceTypeMatches = statmentDetails.External
	}

	hasMatchingPath := true
	if len(gr.modulePaths) > 0 {
		hasMatchingPath = false
		for _, groupRuleModulePath := range gr.modulePaths {
			if gr.statementModulePathMatchesGroupRuleModulePath(statmentDetails.ModulePathParts, groupRuleModulePath) {
				hasMatchingPath = true
				break
			}
		}
	}

	return dependecySourceTypeMatches && hasMatchingPath
}
