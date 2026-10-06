package languageprocessing

import (
	"bytes"
	"importcleaner/internal/datastructures"
	"iter"
	"slices"
	"strings"
)

type TypescriptLanguageDetails struct {
	IsDefault    bool
	IsWildcard   bool
	IsTypeImport bool
}

type TypescriptDependecy struct {
	Name         string
	Alias        string
	IsDefault    bool
	IsWildcard   bool
	IsTypeImport bool
}

type TypescriptImportStatement struct {
	ModulePathParts []string
	ModuleAlias     string
	Dependecies     []TypescriptDependecy
	External        bool
	IsTypeImport    bool
}

type TypescriptImportStatementParser struct{}

const (
	TextStateStandard = iota
	TextStateImportStart
	TextStateImportEnd
)

func (p TypescriptImportStatementParser) StatementsInText(text string) iter.Seq[string] {

	return func(yield func(string) bool) {
		var stringBuffer bytes.Buffer

		lookahead, _ := datastructures.CreateNewLookahead[string](2)

		textState := TextStateStandard
		for lah := range lookahead.ProcessSequence(DelimitedStringSequence(text, []rune{' ', ';', '\n'})) {
			if textState == TextStateStandard && strings.TrimSpace(*lah.PeekStart()) == "import" {
				if stringBuffer.Len() > 0 {
					if !yield(stringBuffer.String()) {
						return
					}
				}
				stringBuffer.Reset()

				stringBuffer.WriteString(*lah.PeekStart())
				textState = TextStateImportStart
				continue
			}

			if textState == TextStateImportStart && strings.TrimSpace(*lah.PeekStart()) == "from" {
				stringBuffer.WriteString(*lah.PeekStart())
				textState = TextStateImportEnd
				continue
			}

			if textState == TextStateImportEnd {
				stringBuffer.WriteString(*lah.PeekStart())

				if lah.PeekEnd() != nil && *lah.PeekEnd() == "\n" && strings.HasSuffix(*lah.PeekStart(), ";") {
					// New lines directly after semicolons should be included as part of the import statement,
					// so don't yield or transition state yet.
					continue
				}

				if !yield(stringBuffer.String()) {
					return
				}
				stringBuffer.Reset()

				textState = TextStateStandard
				continue
			}

			stringBuffer.WriteString(*lah.PeekStart())
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

func (p TypescriptImportStatementParser) ParseImportStatement(statementText string) (ImportStatement, error) {

	buildErrorResponse := func(errorDetail string) (TypescriptImportStatement, error) {
		return TypescriptImportStatement{}, &ImportStatementParseError{
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

func (p TypescriptImportStatement) dependencyAsStatementString(dependency *TypescriptDependecy) string {
	var statementStringBuilder strings.Builder

	if dependency.IsTypeImport {
		statementStringBuilder.WriteString("type ")
	}

	statementStringBuilder.WriteString(dependency.Name)

	if dependency.Alias != "" {
		statementStringBuilder.WriteString(" as ")
		statementStringBuilder.WriteString(dependency.Alias)
	}

	return statementStringBuilder.String()
}

func (p TypescriptImportStatement) AsText() (string, error) {
	if len(p.Dependecies) == 0 {
		return "", &ImportStatementWriteError{
			Statement: p,
			Detail:    "No dependencies were specified",
		}
	}

	var statementStringBuilder strings.Builder

	statementStringBuilder.WriteString("import ")

	if p.IsTypeImport {
		statementStringBuilder.WriteString("type ")
	}

	var defaultDependency *TypescriptDependecy
	otherDependencies := make([]*TypescriptDependecy, 0)
	for _, dependency := range p.Dependecies {
		if dependency.IsDefault {
			if defaultDependency != nil {
				return "", &ImportStatementWriteError{
					Statement: p,
					Detail:    "Multiple default dependencies are not allowed",
				}
			}

			defaultDependency = &dependency
		} else {
			otherDependencies = append(otherDependencies, &dependency)
		}
	}

	if defaultDependency != nil {
		statementStringBuilder.WriteString(p.dependencyAsStatementString(defaultDependency))
	}
	if len(otherDependencies) > 0 {
		if defaultDependency != nil {
			statementStringBuilder.WriteString(", ")
		}

		statementStringBuilder.WriteString("{ ")

		statementStrings := make([]string, len(otherDependencies))
		for i, dependency := range otherDependencies {
			statementStrings[i] = p.dependencyAsStatementString(dependency)
		}
		statementStringBuilder.WriteString(strings.Join(statementStrings, ", "))
		statementStringBuilder.WriteString(" }")
	}

	statementStringBuilder.WriteString(" from '")
	statementStringBuilder.WriteString(strings.Join(p.ModulePathParts, "/"))
	statementStringBuilder.WriteString("';")

	return statementStringBuilder.String(), nil
}

func (p TypescriptImportStatement) GetGenericDetails() ImportStatementGenericDetails {
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
