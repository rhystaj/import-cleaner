package typescript

import (
	lp "importcleaner/internal/languageprocessing"
	"strings"
)

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
		return "", &lp.ImportStatementWriteError{
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
				return "", &lp.ImportStatementWriteError{
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

func (p TypescriptImportStatement) GetGenericDetails() lp.ImportStatementGenericDetails {
	dependencies := make([]lp.GenericDependecy, len(p.Dependecies))
	for i, dependency := range p.Dependecies {
		dependencies[i] = lp.GenericDependecy{
			Name:  dependency.Name,
			Alias: dependency.Alias,
		}
	}

	return lp.ImportStatementGenericDetails{
		ModulePathParts: stripRelativePrefixFromModulePath(p.ModulePathParts),
		ModuleAlias:     p.ModuleAlias,
		Dependecies:     dependencies,
		External:        p.External,
	}
}

func (p TypescriptImportStatement) ConsolidationKey() string {
	modulePath := strings.Join(stripRelativePrefixFromModulePath(p.ModulePathParts), ".")

	if p.IsTypeImport {
		modulePath = "*" + modulePath
	}

	return modulePath
}

func (p TypescriptImportStatement) ConsolidateWith(other lp.ImportStatement) (lp.ImportStatement, error) {
	panic("Not implemented")
}
