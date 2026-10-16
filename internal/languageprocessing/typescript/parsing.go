package typescript

import (
	"bytes"
	lp "importcleaner/internal/languageprocessing"
	"iter"
	"slices"
	"strings"
)

type TypescriptImportStatementParser struct{}

const (
	TextStateStandard = iota
	TextStateImportStart
	TextStateImportEnd
)

func stringAsRuneSequence(text string) iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for _, rune := range text {
			if !yield(rune) {
				return
			}
		}
	}
}

func (p TypescriptImportStatementParser) StatementsInText(text string) iter.Seq[string] {
	return func(yield func(string) bool) {
		var stringBuffer bytes.Buffer

		processingImportStatement := false
		for frag := range fragmentByTypescriptStatementTerminators(text) {
			stringBuffer.WriteString(frag)

			if strings.HasPrefix(frag, "import") {
				processingImportStatement = true
			}

			if strings.Contains(frag, "from") {
				processingImportStatement = false
			}

			if processingImportStatement {
				continue
			}

			if !yield(stringBuffer.String()) {
				return
			}
			stringBuffer.Reset()
		}

		if stringBuffer.Len() > 0 {
			if !yield(stringBuffer.String()) {
				return
			}
		}
	}
}

func (p TypescriptImportStatementParser) IsIntendedImportStatement(statementText string) bool {
	statementStart, _, _ := strings.Cut(statementText, " ")
	return strings.TrimSpace(statementStart) == "import"
}

func (p TypescriptImportStatementParser) parseModulePathParts(modulePartsString string) ([]string, *string) {
	pathWithoutQuotes := strings.Trim(modulePartsString, "\"")
	pathWithoutQuotes = strings.Trim(pathWithoutQuotes, "'")

	return strings.Split(pathWithoutQuotes, "/"), nil
}

func (p TypescriptImportStatementParser) isExternalModule(modulePathParts []string) bool {
	return len(modulePathParts) > 0 && !slices.Contains([]string{".", ".."}, modulePathParts[0])
}

func (p TypescriptImportStatementParser) parseDependency(dependencyText string, isDefault bool) TypescriptDependecy {
	dependencyText = strings.TrimSpace(dependencyText)

	dependencyReferenceText, hasTypeKeyword := strings.CutPrefix(dependencyText, "type ")

	dependencyReferenceParts := strings.Split(dependencyReferenceText, " ")

	result := TypescriptDependecy{
		Name:         dependencyReferenceParts[0],
		IsDefault:    isDefault,
		IsTypeImport: hasTypeKeyword,
	}

	if len(dependencyReferenceParts) >= 3 {
		result.Alias = dependencyReferenceParts[2]
	}

	return result
}

func (p TypescriptImportStatementParser) parseDependenciesList(dependenciesText string) ([]TypescriptDependecy, string) {
	var defaultDependencies, nonDefaultDependencies []TypescriptDependecy

	dependenciesText = strings.TrimSpace(dependenciesText)
	if dependenciesText == "" && len(defaultDependencies) == 0 {
		return []TypescriptDependecy{}, "No dependencies were provided"
	}

	var defaultDependencyText string
	defaultDependencyText, dependenciesText, _ = strings.Cut(dependenciesText, "{")

	if defaultDependencyText != "" {
		defaultDependencyText = strings.TrimSpace(defaultDependencyText)
		defaultDependencyText = strings.TrimSuffix(defaultDependencyText, ",")

		defaultDependecy := p.parseDependency(defaultDependencyText, true)

		defaultDependencies = append(defaultDependencies, defaultDependecy)
	}

	var afterCloseBracket string
	var closingBracketFound bool
	dependenciesText = strings.TrimSpace(dependenciesText)
	dependenciesText, afterCloseBracket, closingBracketFound = strings.Cut(dependenciesText, "}")
	if dependenciesText != "" && !closingBracketFound {
		return []TypescriptDependecy{}, "Dependencies definition missing closing bracket"
	}
	if strings.TrimSpace(afterCloseBracket) != "" {
		return []TypescriptDependecy{}, "No further dependecy definitions expected after closing bracket"
	}

	dependenciesText = strings.TrimSpace(dependenciesText)
	if dependenciesText != "" {
		dependencyStrings := strings.Split(dependenciesText, ",")
		nonDefaultDependencies = make([]TypescriptDependecy, len(dependencyStrings))
		for i, dependencyString := range dependencyStrings {
			nonDefaultDependencies[i] = p.parseDependency(dependencyString, false)
		}
	}

	return append(defaultDependencies, nonDefaultDependencies...), ""
}

func (p TypescriptImportStatementParser) ParseImportStatement(statementText string) (lp.ImportStatement, error) {

	buildErrorResponse := func(errorDetail string) (TypescriptImportStatement, error) {
		return TypescriptImportStatement{}, &lp.ImportStatementParseError{
			Detail:    errorDetail,
			Statement: statementText,
		}
	}

	trimmedStatement := strings.TrimSpace(statementText)

	statementNoImportKeyword, importKeywordExisted := strings.CutPrefix(trimmedStatement, "import")
	if !importKeywordExisted {
		return buildErrorResponse("Statement doesn't start with 'import' keyword")
	}

	statementNoSemicolon, _ := strings.CutSuffix(statementNoImportKeyword, ";")

	dependenciesString, moduleString, fromExisted := strings.Cut(statementNoSemicolon, "from")
	if !fromExisted {
		return buildErrorResponse("Statement doesn't include 'from' keyword")
	}

	trimmedDependenciesString := strings.TrimSpace(dependenciesString)

	isTypeImport := false
	if strings.HasPrefix(trimmedDependenciesString, "type {") {
		isTypeImport = true
		trimmedDependenciesString, _ = strings.CutPrefix(trimmedDependenciesString, "type ")
	}

	dependencies, dependencyParseErrorDetail := p.parseDependenciesList(trimmedDependenciesString)
	if dependencyParseErrorDetail != "" {
		return buildErrorResponse(dependencyParseErrorDetail)
	}

	modulePathParts, moduleParseErrorDetail := p.parseModulePathParts(strings.TrimSpace(moduleString))
	if moduleParseErrorDetail != nil {
		return buildErrorResponse(*moduleParseErrorDetail)
	}

	return TypescriptImportStatement{
		ModulePathParts: modulePathParts,
		External:        p.isExternalModule(modulePathParts),
		Dependecies:     dependencies,
		IsTypeImport:    isTypeImport,
	}, nil
}
