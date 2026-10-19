package languageprocessing

import (
	"iter"
)

type GenericDependecy struct {
	Name  string
	Alias string
}

type ImportStatementGenericDetails struct {
	RawText         string
	ModulePathParts []string
	ModuleAlias     string
	Dependecies     []GenericDependecy
	External        bool
}

type ImportStatement interface {
	AsText() (string, error)
	GetGenericDetails() ImportStatementGenericDetails
	// Statements are considered to be able to be consolidated if they have the same consolidation key
	ConsolidationKey() string
	ConsolidateWith(other ImportStatement) (ImportStatement, error)
}

type ImportStatementParser interface {
	StatementsInText(text string) iter.Seq[string]
	IsIntendedImportStatement(statementText string) bool
	ParseImportStatement(statementText string) (ImportStatement, error)
}
