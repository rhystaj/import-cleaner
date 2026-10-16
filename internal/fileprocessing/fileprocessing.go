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

type FileProcessorImpl[S lp.ImportStatement] struct {
	Parser      lp.ImportStatementParser
	FileManager iowrappers.FileManager
	Organiser   organisation.ImportStatementOrganiser[S]
}

func (fp FileProcessorImpl[S]) extractStatements(text string) ([]S, []string, error) {
	importStatements := make([]S, 0)
	otherStatements := make([]string, 0)

	passedImportsAtTop := false
	for statement := range fp.Parser.StatementsInText(text) {
		if strings.TrimSpace(statement) == "" && !passedImportsAtTop {
			//Whitespace between imports at the top of the file should be ignored
			continue
		}

		if !fp.Parser.IsIntendedImportStatement(statement) {
			otherStatements = append(otherStatements, statement)
			passedImportsAtTop = true // Reached first non-whitespace, non-import statement
			continue
		}

		importStatement, err := fp.Parser.ParseImportStatement(statement)
		if err != nil {
			return make([]S, 0), make([]string, 0), err
		}
		importStatements = append(importStatements, importStatement.(S))
	}

	return importStatements, otherStatements, nil
}

func (fp FileProcessorImpl[S]) ProcessFile(filePath string) error {
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
