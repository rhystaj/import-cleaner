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
}

type ImportStatementParser interface {
	StatementsInText(text string) iter.Seq[string]
	IsIntendedImportStatement(statementText string) bool
	ParseImportStatement(statementText string) (ImportStatement, error)
}
