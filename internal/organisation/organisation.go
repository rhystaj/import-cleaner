package organisation

import (
	lp "importcleaner/internal/languageprocessing"
	"slices"
	"strings"
)

type ImportStatementOrganiser[S lp.ImportStatement] interface {
	OrganiseImportStatements(statements []S) [][]S
}

type ImportStatementOrganiserImpl[S lp.ImportStatement] struct {
	GroupRules []ImportStatementGroupRule[S]
}

func (o ImportStatementOrganiserImpl[S]) OrganiseImportStatements(statements []S) [][]S {

	groups := make([][]S, len(o.GroupRules)+1)
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

	nonEmptyGroups := make([][]S, 0)
	for _, group := range groups {
		if len(group) > 0 {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	var sortedGroups [][]S = make([][]S, len(nonEmptyGroups))
	for i, group := range nonEmptyGroups {
		sortedGroups[i] = slices.SortedFunc(
			func(yield func(S) bool) {
				for _, item := range group {
					if !(yield(item)) {
						return
					}
				}
			},
			func(a S, b S) int {
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
