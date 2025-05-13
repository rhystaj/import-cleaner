package organisation

import (
	"importcleaner/internal/parsing"
	"slices"
	"strings"
)

type DependencySourceType int

const (
	DependencySourceTypeNone DependencySourceType = iota
	DependencySourceTypeInternal
	DependencySourceTypeExternal
)

type ImportStatementGroupRule struct {
	DependencySourceType DependencySourceType
	ModulePaths          [][]string
}

func statementModulePathMatchesGroupRuleModulePath(statementModulePath []string, groupRuleModulePath []string) bool {
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

func statementMatchesGroupRule(statement parsing.ImportStatement, groupRule ImportStatementGroupRule) bool {
	var dependecySourceTypeMatches bool
	switch groupRule.DependencySourceType {
	case DependencySourceTypeNone:
		dependecySourceTypeMatches = true
	case DependencySourceTypeInternal:
		dependecySourceTypeMatches = !statement.External
	case DependencySourceTypeExternal:
		dependecySourceTypeMatches = statement.External
	}

	hasMatchingPath := true
	if len(groupRule.ModulePaths) > 0 {
		hasMatchingPath = false
		for _, groupRuleModulePath := range groupRule.ModulePaths {
			if statementModulePathMatchesGroupRuleModulePath(statement.ModulePathParts, groupRuleModulePath) {
				hasMatchingPath = true
				break
			}
		}
	}

	return dependecySourceTypeMatches && hasMatchingPath
}

type ImportStatementOrganiser interface {
	OrganiseImportStatements(statements []parsing.ImportStatement) [][]parsing.ImportStatement
}

type ImportStatementOrganiserImpl struct {
	GroupRules []ImportStatementGroupRule
}

func (o ImportStatementOrganiserImpl) OrganiseImportStatements(statements []parsing.ImportStatement) [][]parsing.ImportStatement {

	groups := make([][]parsing.ImportStatement, len(o.GroupRules)+1)
	for _, statement := range statements {
		statementGrouped := false
		for groupIndex, groupRule := range o.GroupRules {
			if statementMatchesGroupRule(statement, groupRule) {
				groups[groupIndex] = append(groups[groupIndex], statement)
				statementGrouped = true
				break
			}
		}
		if !statementGrouped {
			groups[len(o.GroupRules)] = append(groups[len(o.GroupRules)], statement)
		}
	}

	nonEmptyGroups := make([][]parsing.ImportStatement, 0)
	for _, group := range groups {
		if len(group) > 0 {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	var sortedGroups [][]parsing.ImportStatement = make([][]parsing.ImportStatement, len(nonEmptyGroups))
	for i, group := range nonEmptyGroups {
		sortedGroups[i] = slices.SortedFunc(
			func(yield func(parsing.ImportStatement) bool) {
				for _, item := range group {
					if !(yield(item)) {
						return
					}
				}
			},
			func(a parsing.ImportStatement, b parsing.ImportStatement) int {
				return strings.Compare(
					strings.Join(a.ModulePathParts, ""),
					strings.Join(b.ModulePathParts, ""),
				)
			})
	}

	return sortedGroups
}
