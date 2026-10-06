package languageprocessing

import (
	"bytes"
	"fmt"
	"importcleaner/internal/datastructures"
	"importcleaner/internal/iowrappers"
	"iter"
	"strings"
)

type PythonDependecy struct {
	Name  string
	Alias string
}

type PythonImportStatement struct {
	RawText         string
	ModulePathParts []string
	ModuleAlias     string
	Dependecies     []PythonDependecy
	External        bool
}

type PythonImportStatementParser struct {
	WorkingDir  string
	FileManager iowrappers.FileManager
}

func (p PythonImportStatementParser) StatementsInText(text string) iter.Seq[string] {
	ignoreNextLineBreak := false
	bracketStack := datastructures.CreateNewStack[rune]()

	var statementBuffer bytes.Buffer

	return func(yield func(string) bool) {
		for _, r := range text {
			statementBuffer.WriteRune(r)

			if r == '\n' && bracketStack.Size() <= 0 {

				if ignoreNextLineBreak {
					ignoreNextLineBreak = false
					continue
				}

				if !(yield(statementBuffer.String())) {
					return
				}

				statementBuffer.Reset()
				continue
			}

			if r == '/' {
				ignoreNextLineBreak = true
			} else if IsOpenBracket(r) {
				bracketStack.Push(&r)
			} else if bracketStack.Size() > 0 && IsCorrespondingCloseBracket(r, *bracketStack.Peek()) {
				bracketStack.Pop()
			}
		}

		yield(statementBuffer.String())
	}

}

func (p PythonImportStatementParser) IsIntendedImportStatement(statementText string) bool {
	statementParts := strings.Split(statementText, " ")

	if len(statementParts) < 2 {
		return false
	}

	if strings.TrimSpace(statementParts[0]) == "import" {
		return true
	}

	if len(statementParts) < 4 {
		return false
	}

	return strings.TrimSpace(statementParts[0]) == "from" && strings.TrimSpace(statementParts[2]) == "import"
}

func (p PythonImportStatementParser) isExternalModule(modulePathParts []string) bool {

	if modulePathParts[0] == "" {
		//Relative import, so can't be external
		return false
	}

	for e := range p.FileManager.ReadDirectory(p.WorkingDir) {
		if !e.IsDir && e.ItemName == modulePathParts[0]+".py" {
			return false
		}

		if e.IsDir && e.ItemName == modulePathParts[0] {
			return false
		}
	}

	return true
}

func (p PythonImportStatementParser) parseModuleInfo(moduleInfoText string) (modulePathParts []string, moduleAlias string, errorMessage string) {
	parts := strings.Split(moduleInfoText, " ")

	if len(parts) == 0 {
		return []string{}, "", "No module"
	}

	modulePathParts = strings.Split(parts[0], ".")
	if len(parts) <= 1 {
		return modulePathParts, "", ""
	}

	if len(parts) != 3 || parts[1] != "as" {
		return []string{}, "", "Module alias expected"
	}

	return modulePathParts, parts[2], ""
}

func (p PythonImportStatementParser) parseDependency(dependecyText string) (dependency PythonDependecy, errorMessage string) {
	parts := strings.Split(strings.Trim(dependecyText, " "), " ")

	if len(parts) <= 1 {
		return PythonDependecy{
			Name:  parts[0],
			Alias: "",
		}, ""
	}

	if len(parts) >= 2 && parts[1] != "as" {
		return PythonDependecy{}, fmt.Sprintf("Malformed dependecy '%s'", dependecyText)
	}

	if parts[1] == "as" && len(parts) == 2 {
		return PythonDependecy{}, "Dependency alias expected"
	}

	return PythonDependecy{
		Name:  parts[0],
		Alias: parts[2],
	}, ""
}

func (p PythonImportStatementParser) parseDependenciesList(dependenciesListText string) (dependencies []PythonDependecy, errorMessage string) {
	result := make([]PythonDependecy, 0)

	currentStatementPart := ""
	statementRemaining := dependenciesListText

	for len(statementRemaining) > 0 {
		currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, ",")

		charsToTrim := []string{"(", "\t", "\n", ")"}
		dependecyText := strings.TrimSpace(currentStatementPart)
		for _, char := range charsToTrim {
			dependecyText = strings.ReplaceAll(dependecyText, char, "")
		}

		dependecy, dependecyParseErrorMessage := p.parseDependency(dependecyText)
		if dependecyParseErrorMessage != "" {
			return []PythonDependecy{}, dependecyParseErrorMessage
		}

		if len(dependecyText) > 0 {
			result = append(result, dependecy)
		}
	}

	return result, ""
}

func (p PythonImportStatementParser) generateErrorResultForStatement(statement string, detail string) (ImportStatement, *ImportStatementParseError) {
	error := &ImportStatementParseError{
		Statement: statement,
		Detail:    detail,
	}

	return PythonImportStatement{"", []string{}, "", []PythonDependecy{}, false}, error
}

func (p PythonImportStatementParser) ParseImportStatement(statementText string) (ImportStatement, error) {
	currentStatementPart := ""
	statementRemaining := strings.TrimSpace(statementText)

	currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, " ")
	if currentStatementPart != "import" && currentStatementPart != "from" {
		return p.generateErrorResultForStatement(statementText, "Not import statement")
	}

	if currentStatementPart == "import" {
		if statementRemaining == "" {
			return p.generateErrorResultForStatement(statementText, "No module")
		}

		modulePathParts, moduleAlias, moduleParseErrorMessage := p.parseModuleInfo(statementRemaining)
		if moduleParseErrorMessage != "" {
			return p.generateErrorResultForStatement(statementText, moduleParseErrorMessage)
		}

		return PythonImportStatement{
			RawText:         statementText,
			ModulePathParts: modulePathParts,
			ModuleAlias:     moduleAlias,
			Dependecies:     []PythonDependecy{},
			External:        p.isExternalModule(modulePathParts),
		}, nil
	}

	currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, " ")
	modulePathParts, _, moduleParseErrorMessage := p.parseModuleInfo(currentStatementPart)
	if moduleParseErrorMessage != "" {
		return p.generateErrorResultForStatement(statementText, moduleParseErrorMessage)
	}

	currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, " ")
	if currentStatementPart != "import" {
		return p.generateErrorResultForStatement(statementText, "'import' keyword expected")
	}

	if statementRemaining == "" {
		return p.generateErrorResultForStatement(statementText, "No dependecies listed")
	}

	dependencies, dependeciesListParseErrorMessage := p.parseDependenciesList(statementRemaining)
	if dependeciesListParseErrorMessage != "" {
		return p.generateErrorResultForStatement(statementText, dependeciesListParseErrorMessage)
	}

	external := p.isExternalModule(modulePathParts)

	return PythonImportStatement{statementText, modulePathParts, "", dependencies, external}, nil
}

func (p PythonImportStatement) dependecyAsString(dependecy PythonDependecy) string {
	result := dependecy.Name

	if dependecy.Alias != "" {
		result += " as " + dependecy.Alias
	}

	return result
}

func (p PythonImportStatement) AsText() (string, error) {
	modulePath := strings.Join(p.ModulePathParts, ".")

	if len(p.Dependecies) > 0 {
		dependencyNames := make([]string, 0)
		for _, dep := range p.Dependecies {
			dependencyNames = append(dependencyNames, p.dependecyAsString(dep))
		}

		result := fmt.Sprintf("from %s import %s", modulePath, strings.Join(dependencyNames, ", "))
		return result, nil
	} else {
		result := fmt.Sprintf("import %s", modulePath)
		return result, nil
	}
}

func (p PythonImportStatement) GetGenericDetails() ImportStatementGenericDetails {
	dependencies := make([]GenericDependecy, len(p.Dependecies))
	for i, dependency := range p.Dependecies {
		dependencies[i] = GenericDependecy{
			Name:  dependency.Name,
			Alias: dependency.Alias,
		}
	}

	return ImportStatementGenericDetails{
		ModulePathParts: p.ModulePathParts,
		ModuleAlias:     p.ModuleAlias,
		Dependecies:     dependencies,
		External:        p.External,
	}
}
