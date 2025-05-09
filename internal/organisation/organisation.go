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
	ModulePathPrefix     string
}

func statementMatchesGroupRule(statement parsing.ImportStatement, groupRule ImportStatementGroupRule) bool {
	if groupRule.DependencySourceType == DependencySourceTypeNone {
		return true
	}

	return statement.External && groupRule.DependencySourceType == DependencySourceTypeExternal ||
		!statement.External && groupRule.DependencySourceType == DependencySourceTypeInternal
}

func OrganiseImportStatements(
	statements []parsing.ImportStatement,
	groupRules []ImportStatementGroupRule,
) [][]parsing.ImportStatement {

	groups := make([][]parsing.ImportStatement, len(groupRules)+1)
	for _, statement := range statements {
		statementGrouped := false
		for groupIndex, groupRule := range groupRules {
			if statementMatchesGroupRule(statement, groupRule) {
				groups[groupIndex] = append(groups[groupIndex], statement)
				statementGrouped = true
				break
			}
		}
		if !statementGrouped {
			groups[len(groupRules)] = append(groups[len(groupRules)], statement)
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
				return strings.Compare(a.ModulePathParts[0], b.ModulePathParts[0])
			})
	}

	return sortedGroups
}
