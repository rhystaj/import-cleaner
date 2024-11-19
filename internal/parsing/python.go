package parsing

import (
	"bytes"
	"fmt"
	"importcleaner/internal/datastructures"
	"importcleaner/internal/iowrappers"
	"iter"
	"strings"
)

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
			if r == '\n' && bracketStack.Size() <= 0 {
				if ignoreNextLineBreak {
					ignoreNextLineBreak = false
					statementBuffer.WriteRune(r)
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

			statementBuffer.WriteRune(r)
		}

		yield(statementBuffer.String())
	}

}

func (p PythonImportStatementParser) IsImportStatement(statementText string) bool {
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

func (p PythonImportStatementParser) isExternalModule(moduleName string) bool {

	if strings.HasPrefix(moduleName, ".") {
		//Relative import, so can't be external
		return false
	}

	for e := range p.FileManager.ReadDirectory(p.WorkingDir) {
		if !e.IsDir && e.ItemName == moduleName+".py" {
			return false
		}

		if e.IsDir && e.ItemName == moduleName {
			return false
		}
	}

	return true
}

func (p PythonImportStatementParser) parseModuleInfo(moduleInfoText string) (moduleName string, moduleAlias string, successful bool) {
	parts := strings.Split(moduleInfoText, " ")

	if len(parts) <= 1 {
		return parts[0], "", true
	}

	if len(parts) != 3 || parts[1] != "as" {
		return "", "", false
	}

	return parts[0], parts[2], true
}

func (p PythonImportStatementParser) parseDependency(dependecyText string) (dependency Dependecy, successful bool) {
	parts := strings.Split(dependecyText, " ")

	if len(parts) <= 1 {
		return Dependecy{
			Name:  parts[0],
			Alias: "",
		}, true
	}

	if len(parts) != 3 || parts[1] != "as" {
		return Dependecy{}, false
	}

	return Dependecy{
		Name:  parts[0],
		Alias: parts[2],
	}, true
}

func (p PythonImportStatementParser) parseDependenciesList(dependenciesListText string) (dependencies []Dependecy, successful bool) {
	result := make([]Dependecy, 0)

	currentStatementPart := ""
	statementRemaining := dependenciesListText

	for len(statementRemaining) > 0 {
		currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, ",")

		charsToTrim := []string{"(", "\t", "\n", ")"}
		dependecyText := strings.TrimSpace(currentStatementPart)
		for _, char := range charsToTrim {
			dependecyText = strings.ReplaceAll(dependecyText, char, "")
		}

		dependecy, dependencyParseSuccessful := p.parseDependency(dependecyText)
		if !dependencyParseSuccessful {
			return []Dependecy{}, false
		}

		if len(dependecyText) > 0 {
			result = append(result, dependecy)
		}
	}

	return result, true
}

func (p PythonImportStatementParser) generateErrorResultForStatement(statement string) (ImportStatement, *ImportStatementParseError) {
	error := &ImportStatementParseError{
		Statement: statement,
	}

	return ImportStatement{"", "", []Dependecy{}, false}, error
}

func (p PythonImportStatementParser) ParseImportStatement(statementText string) (ImportStatement, error) {
	currentStatementPart := ""
	statementRemaining := strings.TrimSpace(statementText)

	currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, " ")
	if currentStatementPart != "import" && currentStatementPart != "from" {
		return p.generateErrorResultForStatement(statementText)
	}

	if currentStatementPart == "import" {
		if statementRemaining == "" {
			return p.generateErrorResultForStatement(statementText)
		}

		moduleName, moduleAlias, moduleParseSuccessful := p.parseModuleInfo(statementRemaining)
		if !moduleParseSuccessful {
			return p.generateErrorResultForStatement(statementText)
		}

		return ImportStatement{
			ModuleName:  moduleName,
			ModuleAlias: moduleAlias,
			Dependecies: []Dependecy{},
			External:    p.isExternalModule(moduleName),
		}, nil
	}

	currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, " ")
	moduleName := currentStatementPart

	currentStatementPart, statementRemaining, _ = strings.Cut(statementRemaining, " ")
	if currentStatementPart != "import" {
		return p.generateErrorResultForStatement(statementText)
	}

	if statementRemaining == "" {
		return p.generateErrorResultForStatement(statementText)
	}

	dependencies, dependeciesListParseSuccessful := p.parseDependenciesList(statementRemaining)
	if !dependeciesListParseSuccessful {
		return p.generateErrorResultForStatement(statementText)
	}

	external := p.isExternalModule(moduleName)

	return ImportStatement{moduleName, "", dependencies, external}, nil
}

func (p PythonImportStatementParser) dependecyAsString(dependecy Dependecy) string {
	result := dependecy.Name

	if dependecy.Alias != "" {
		result += " as " + dependecy.Alias
	}

	return result
}

func (p PythonImportStatementParser) StatementAsString(statement ImportStatement) string {
	if len(statement.Dependecies) > 0 {
		dependencyNames := make([]string, 0)
		for _, dep := range statement.Dependecies {
			dependencyNames = append(dependencyNames, p.dependecyAsString(dep))
		}

		return fmt.Sprintf("from %s import %s", statement.ModuleName, strings.Join(dependencyNames, ", "))
	} else {
		return fmt.Sprintf("import %s", statement.ModuleName)
	}
}
