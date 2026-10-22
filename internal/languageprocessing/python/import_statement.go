package python

import (
	"fmt"
	lp "importcleaner/internal/languageprocessing"
	"strings"
)

type PythonDependecy struct {
	Name  string
	Alias string
}

type PythonImportStatement struct {
	AncestoryLevels int
	ModulePathParts []string
	ModuleAlias     string
	Dependecies     []PythonDependecy
	External        bool
}

func (p PythonImportStatement) dependecyAsString(dependecy PythonDependecy) string {
	result := dependecy.Name

	if dependecy.Alias != "" {
		result += " as " + dependecy.Alias
	}

	return result
}

func (p PythonImportStatement) AsText() (string, error) {
	modulePath := strings.Repeat(".", p.AncestoryLevels) + strings.Join(p.ModulePathParts, ".")

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

func (p PythonImportStatement) GetGenericDetails() lp.ImportStatementGenericDetails {
	dependencies := make([]lp.GenericDependecy, len(p.Dependecies))
	for i, dependency := range p.Dependecies {
		dependencies[i] = lp.GenericDependecy{
			Name:  dependency.Name,
			Alias: dependency.Alias,
		}
	}

	return lp.ImportStatementGenericDetails{
		ModulePathParts: p.ModulePathParts,
		ModuleAlias:     p.ModuleAlias,
		Dependecies:     dependencies,
		External:        p.External,
	}
}

func (p PythonImportStatement) ConsolidationKey() string {
	var stringBuilder strings.Builder

	if len(p.Dependecies) == 0 {
		stringBuilder.WriteString("$")
	}

	if len(p.ModulePathParts) > 0 {
		stringBuilder.WriteString(strings.Join(p.ModulePathParts, "."))
	}

	if p.ModuleAlias != "" {
		stringBuilder.WriteString("(")
		stringBuilder.WriteString(p.ModuleAlias)
		stringBuilder.WriteString(")")
	}

	return stringBuilder.String()
}

func (p PythonImportStatement) ConsolidateWith(other lp.ImportStatement) (lp.ImportStatement, error) {
	otherPythonImportStatement, ok := other.(PythonImportStatement)
	if !ok {
		return nil, fmt.Errorf("other is not a Python import statement")
	}

	if p.ConsolidationKey() != otherPythonImportStatement.ConsolidationKey() {
		return nil, fmt.Errorf("consolidation keys do not match")
	}

	if p.External != otherPythonImportStatement.External {
		return nil, fmt.Errorf("externalities do not match")
	}

	result := PythonImportStatement{
		ModulePathParts: p.ModulePathParts,
		ModuleAlias:     p.ModuleAlias,
		Dependecies:     append(p.Dependecies, otherPythonImportStatement.Dependecies...),
		External:        p.External,
	}

	return result, nil
}
