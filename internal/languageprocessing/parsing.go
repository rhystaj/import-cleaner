package languageprocessing

import (
	"iter"
)

type GenericDependecy struct {
	Name  string
	Alias string
}

type ImportStatementGenericDetails struct {
	ModulePathParts []string
	ModuleAlias     string
	Dependecies     []GenericDependecy
	External        bool
}

type TextOutputOptions struct {
	BreakLinesAtNumChars int
}

type ImportStatement interface {
	AsText(options TextOutputOptions) (string, error)
	GetGenericDetails() ImportStatementGenericDetails
	// Statements are considered to be able to be consolidated if they have the same consolidation key
	ConsolidationKey() string
	ConsolidateWith(other ImportStatement) (ImportStatement, error)
	Compare(other ImportStatement) int
}

type ImportStatementParser interface {
	StatementsInText(text string) iter.Seq[string]
	IsIntendedImportStatement(statementText string) bool
	ParseImportStatement(statementText string) (ImportStatement, error)
}
