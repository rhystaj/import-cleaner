package fileprocessing

import (
	"bytes"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
	"strings"
)

type FileProcessor interface {
	ProcessFile(filePath string, parser parsing.ImportStatementParser) error
}

type FileProcessorImpl struct {
	FileManager iowrappers.FileManager
}

func (fp FileProcessorImpl) extractStatements(text string, parser parsing.ImportStatementParser) ([]parsing.ImportStatement, []string, error) {
	importStatements := make([]parsing.ImportStatement, 0)
	otherStatements := make([]string, 0)

	var ignoreEmptyStatement bool
	for statement := range parser.StatementsInText(text) {
		if parser.IsImportStatement(statement) {
			importStatement, err := parser.ParseImportStatement(statement)
			if err != nil {
				return make([]parsing.ImportStatement, 0), make([]string, 0), err
			}
			importStatements = append(importStatements, importStatement)
			ignoreEmptyStatement = true
			continue
		}

		if len(statement) > 0 || !ignoreEmptyStatement {
			otherStatements = append(otherStatements, statement)
		}

		ignoreEmptyStatement = false
	}

	return importStatements, otherStatements, nil
}

func (fp FileProcessorImpl) ProcessFile(filePath string, parser parsing.ImportStatementParser) error {
	fileContents, fileReadError := fp.FileManager.ReadStringFromFile(filePath)
	if fileReadError != nil {
		return fileReadError
	}

	importStatements, otherStatements, extractionError := fp.extractStatements(fileContents, parser)
	if extractionError != nil {
		return extractionError
	}

	statementGroups := organisation.OrganiseImportStatements(
		importStatements,
		[]organisation.ImportStatementGroupRule{},
	)

	var outputContentsBuffer bytes.Buffer
	for _, group := range statementGroups {
		if len(group) == 0 {
			continue
		}

		for _, statement := range group {
			statementAsText := parser.StatementAsString(statement)
			outputContentsBuffer.WriteString(statementAsText + "\n")
		}
		outputContentsBuffer.WriteString("\n")
	}

	outputContentsBuffer.WriteString(strings.Join(otherStatements, "\n"))

	fp.FileManager.WriteContentsToFile(filePath, outputContentsBuffer.String())

	return nil
}
