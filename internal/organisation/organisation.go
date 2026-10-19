package organisation

import (
	lp "importcleaner/internal/languageprocessing"
	"iter"
	"slices"
	"strings"
)

type ImportStatementOrganiser[S lp.ImportStatement] interface {
	OrganiseImportStatements(statements []S) [][]S
}

type ImportStatementOrganiserImpl[S lp.ImportStatement] struct {
	GroupRules                  []ImportStatementGroupRule[S]
	ConsolidateImportStatements bool
}

func asSeq[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range slice {
			if !(yield(item)) {
				return
			}
		}
	}
}

func byModulePath[S lp.ImportStatement](a S, b S) int {
	aDetails := a.GetGenericDetails()
	bDetails := b.GetGenericDetails()

	return strings.Compare(
		strings.Join(aDetails.ModulePathParts, ""),
		strings.Join(bDetails.ModulePathParts, ""),
	)
}

func (o ImportStatementOrganiserImpl[S]) consolidateStatements(statements []S) []S {
	panic("not implemented")
}

func (o ImportStatementOrganiserImpl[S]) groupStatements(statements []S) [][]S {
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

	return nonEmptyGroups
}

func (o ImportStatementOrganiserImpl[S]) OrganiseImportStatements(statements []S) [][]S {
	if o.ConsolidateImportStatements {
		statements = o.consolidateStatements(statements)
	}

	groups := o.groupStatements(statements)

	var sortedGroups [][]S = make([][]S, len(groups))
	for i, group := range groups {
		sortedGroups[i] = slices.SortedFunc(asSeq(group), byModulePath)
	}

	return sortedGroups
}
