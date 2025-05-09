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

func OrganiseImportStatements(
	statements []parsing.ImportStatement,
	groupRules []ImportStatementGroupRule,
) [][]parsing.ImportStatement {

	groups := [][]parsing.ImportStatement{
		make([]parsing.ImportStatement, 0),
		make([]parsing.ImportStatement, 0),
	}

	for _, statement := range statements {
		if statement.External {
			groups[0] = append(groups[0], statement)
		} else {
			groups[1] = append(groups[1], statement)
		}
	}

	var sortedGroups [][]parsing.ImportStatement = make([][]parsing.ImportStatement, 2)
	for i, group := range groups {
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
