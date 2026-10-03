package parsing

import (
	"fmt"
	"reflect"
)

type ImportStatementParseError struct {
	Statement string
	Detail    string
}

type ImportStatementWriteError struct {
	Statement ImportStatement
	Detail    string
}

func (e *ImportStatementParseError) Error() string {
	message := fmt.Sprintf("%s is not a valid import statement", e.Statement)
	if e.Detail != "" {
		message += " - " + e.Detail
	}

	return message
}

func (e *ImportStatementParseError) Is(otherError error) bool {
	switch oe := otherError.(type) {
	case *ImportStatementParseError:
		return oe.Statement == e.Statement && oe.Detail == e.Detail
	default:
		return false
	}
}

func (e *ImportStatementWriteError) Error() string {
	message := fmt.Sprintf("Failed to write import statement: %+v", e.Statement)
	if e.Detail != "" {
		message += " - " + e.Detail
	}

	return message
}

func (e *ImportStatementWriteError) Is(otherError error) bool {
	switch oe := otherError.(type) {
	case *ImportStatementWriteError:
		return reflect.DeepEqual(e, oe)
	default:
		return false
	}
}
