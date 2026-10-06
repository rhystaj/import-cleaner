package organisation

import (
	lp "importcleaner/internal/languageprocessing"
	"slices"
	"strings"
)

type ImportStatementOrganiser interface {
	OrganiseImportStatements(statements []lp.ImportStatement) [][]lp.ImportStatement
}

type ImportStatementOrganiserImpl struct {
	GroupRules []ImportStatementGroupRule
}

func (o ImportStatementOrganiserImpl) OrganiseImportStatements(statements []lp.ImportStatement) [][]lp.ImportStatement {

	groups := make([][]lp.ImportStatement, len(o.GroupRules)+1)
	for _, statement := range statements {
		statementGrouped := false
		for groupIndex, groupRule := range o.GroupRules {
			if groupRule.AppliesTo(statement) {
				groups[groupIndex] = append(groups[groupIndex], statement)
				statementGrouped = true
				break
			}
		}
		if !statementGrouped {
			groups[len(o.GroupRules)] = append(groups[len(o.GroupRules)], statement)
		}
	}

	nonEmptyGroups := make([][]lp.ImportStatement, 0)
	for _, group := range groups {
		if len(group) > 0 {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	var sortedGroups [][]lp.ImportStatement = make([][]lp.ImportStatement, len(nonEmptyGroups))
	for i, group := range nonEmptyGroups {
		sortedGroups[i] = slices.SortedFunc(
			func(yield func(lp.ImportStatement) bool) {
				for _, item := range group {
					if !(yield(item)) {
						return
					}
				}
			},
			func(a lp.ImportStatement, b lp.ImportStatement) int {
				aDetails := a.GetGenericDetails()
				bDetails := b.GetGenericDetails()

				return strings.Compare(
					strings.Join(aDetails.ModulePathParts, ""),
					strings.Join(bDetails.ModulePathParts, ""),
				)
			})
	}

	return sortedGroups
}
