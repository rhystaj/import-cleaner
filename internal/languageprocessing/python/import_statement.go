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
