package parsing

import (
	"fmt"
	"iter"
)

type Dependecy struct {
	Name  string
	Alias string
}

type ImportStatement struct {
	ModulePathParts []string
	ModuleAlias     string
	Dependecies     []Dependecy
	External        bool
}

type ImportStatementParseError struct {
	Statement string
}

func (e *ImportStatementParseError) Error() string {
	return fmt.Sprintf("%s is not a valid import statement", e.Statement)
}

func (e *ImportStatementParseError) Is(otherError error) bool {
	switch oe := otherError.(type) {
	case *ImportStatementParseError:
		return oe.Statement == e.Statement
	default:
		return false
	}
}

type ImportStatementParser interface {
	StatementsInText(text string) iter.Seq[string]
	IsImportStatement(statementText string) bool
	ParseImportStatement(statementText string) (ImportStatement, error)
	StatementAsString(statement ImportStatement) string
}
