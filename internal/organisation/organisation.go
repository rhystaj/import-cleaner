package organisation

import (
	"fmt"
	lp "importcleaner/internal/languageprocessing"
	"iter"
	"slices"
)

type ImportStatementOrganiser[S lp.ImportStatement] interface {
	OrganiseImportStatements(statements []S) ([][]S, error)
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

func (o ImportStatementOrganiserImpl[S]) consolidateStatements(statements []S) ([]S, error) {
	statementsByConsolidationKey := make(map[string][]S)
	for _, statement := range statements {
		consolidationKey := statement.ConsolidationKey()
		statementsByConsolidationKey[consolidationKey] = append(statementsByConsolidationKey[consolidationKey], statement)
	}

	consolidatedStatements := make([]S, 0)
	for _, statements := range statementsByConsolidationKey {
		if len(statements) <= 0 {
			// Defensive
			continue
		}

		consolidatedStatement := statements[0]
		if len(statements) > 1 {
			for _, statement := range statements[1:] {
				newStatement, consolidationError := consolidatedStatement.ConsolidateWith(statement)
				if consolidationError != nil {
					return make([]S, 0), consolidationError
				}

				var isCorrectType bool = false
				consolidatedStatement, isCorrectType = newStatement.(S)
				if !isCorrectType {
					return make([]S, 0), fmt.Errorf("Consolidated statement is not of type %T", consolidatedStatement)
				}
			}
		}

		consolidatedStatements = append(consolidatedStatements, consolidatedStatement)
	}

	return consolidatedStatements, nil
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

func (o ImportStatementOrganiserImpl[S]) OrganiseImportStatements(statements []S) ([][]S, error) {
	if o.ConsolidateImportStatements {
		var consolidaionError error = nil
		statements, consolidaionError = o.consolidateStatements(statements)
		if consolidaionError != nil {
			return make([][]S, 0), consolidaionError
		}
	}

	groups := o.groupStatements(statements)

	var sortedGroups [][]S = make([][]S, len(groups))
	for i, group := range groups {
		sortedGroups[i] = slices.SortedFunc(asSeq(group), func(a S, b S) int {
			return a.Compare(b)
		})
	}

	return sortedGroups, nil
}
