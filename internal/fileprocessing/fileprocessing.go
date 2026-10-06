package fileprocessing

import (
	"bytes"
	"importcleaner/internal/iowrappers"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/organisation"
	"strings"
)

type FileProcessor interface {
	ProcessFile(filePath string) error
}

type FileProcessorImpl struct {
	Parser      lp.ImportStatementParser
	FileManager iowrappers.FileManager
	Organiser   organisation.ImportStatementOrganiser
}

func (fp FileProcessorImpl) extractStatements(text string) ([]lp.ImportStatement, []string, error) {
	importStatements := make([]lp.ImportStatement, 0)
	otherStatements := make([]string, 0)

	var ignoreEmptyStatement bool
	for statement := range fp.Parser.StatementsInText(text) {
		if fp.Parser.IsIntendedImportStatement(statement) {
			importStatement, err := fp.Parser.ParseImportStatement(statement)
			if err != nil {
				return make([]lp.ImportStatement, 0), make([]string, 0), err
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

func (fp FileProcessorImpl) ProcessFile(filePath string) error {
	fileContents, fileReadError := fp.FileManager.ReadStringFromFile(filePath)
	if fileReadError != nil {
		return fileReadError
	}

	importStatements, otherStatements, extractionError := fp.extractStatements(fileContents)
	if extractionError != nil {
		return extractionError
	}

	statementGroups := fp.Organiser.OrganiseImportStatements(importStatements)

	var outputContentsBuffer bytes.Buffer
	for _, group := range statementGroups {
		if len(group) == 0 {
			continue
		}

		for _, statement := range group {
			statementAsText, _ := statement.AsText() //TODO: Properly handle error
			outputContentsBuffer.WriteString(statementAsText)
			outputContentsBuffer.WriteString("\n")
		}
		outputContentsBuffer.WriteString("\n")
	}

	outputContentsBuffer.WriteString(strings.Join(otherStatements, ""))

	fp.FileManager.WriteContentsToFile(filePath, outputContentsBuffer.String())

	return nil
}
