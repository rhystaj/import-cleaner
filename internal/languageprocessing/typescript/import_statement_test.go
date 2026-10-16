package typescript_test

import (
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/languageprocessing/typescript"
	"testing"

	"github.com/stretchr/testify/assert"
)

type moduleAsStatementStringCase struct {
	Description       string
	ModulePathParts   []string
	ExpectedStatement string
}

func moduleAsStatementStringCases() []moduleAsStatementStringCase {
	return []moduleAsStatementStringCase{
		{
			Description:       "SingleModule",
			ModulePathParts:   []string{"somecoolmodule"},
			ExpectedStatement: "'somecoolmodule'",
		},
		{
			Description:       "NestedModule",
			ModulePathParts:   []string{"somecoolmodule", "child"},
			ExpectedStatement: "'somecoolmodule/child'",
		},
		{
			Description:       "NestedModuleWithMultipleParts",
			ModulePathParts:   []string{"somecoolmodule", "child", "grandchild"},
			ExpectedStatement: "'somecoolmodule/child/grandchild'",
		},
		{
			Description:       "SiblingModule",
			ModulePathParts:   []string{".", "sibling"},
			ExpectedStatement: "'./sibling'",
		},
		{
			Description:       "CousinModule",
			ModulePathParts:   []string{"..", "..", "cousin"},
			ExpectedStatement: "'../../cousin'",
		},
		{
			Description:       "NestedSiblingModule",
			ModulePathParts:   []string{".", "sibling", "child"},
			ExpectedStatement: "'./sibling/child'",
		},
		{
			Description:       "NestedCousinModule",
			ModulePathParts:   []string{"..", "..", "cousin", "child"},
			ExpectedStatement: "'../../cousin/child'",
		},
	}
}

func Test_TypescriptImportStatement_AsText_ValidNonTypeImport(t *testing.T) {
	dependencyCases := generateValidDependencySetPartCases(generateValidDependencySetPartCasesArgs{
		IncludeDefaultCases:    true,
		IncludeTypeImportCases: true,
	})

	moduleCases := moduleAsStatementStringCases()

	for _, dependencyCase := range dependencyCases {
		for _, moduleCase := range moduleCases {
			expectedStatement := "import " + dependencyCase.Text + " from " + moduleCase.ExpectedStatement + ";"
			testStatement := typescript.TypescriptImportStatement{
				ModulePathParts: moduleCase.ModulePathParts,
				External:        false,
				Dependecies:     dependencyCase.Dependencies,
			}

			result, error := testStatement.AsText()

			t.Run(expectedStatement, func(t_ *testing.T) {
				assert.Nil(t_, error)
				assert.Equal(t_, expectedStatement, result)
			})
		}
	}
}

func Test_TypescriptImportStatement_AsText_ValidTypeImport(t *testing.T) {
	dependencyCases := generateValidDependencySetPartCases(generateValidDependencySetPartCasesArgs{})

	moduleCases := moduleAsStatementStringCases()

	for _, dependencyCase := range dependencyCases {
		for _, moduleCase := range moduleCases {
			expectedStatement := "import type " + dependencyCase.Text + " from " + moduleCase.ExpectedStatement + ";"
			testStatement := typescript.TypescriptImportStatement{
				ModulePathParts: moduleCase.ModulePathParts,
				External:        false,
				Dependecies:     dependencyCase.Dependencies,
				IsTypeImport:    true,
			}

			result, error := testStatement.AsText()

			t.Run(expectedStatement, func(t_ *testing.T) {
				assert.Nil(t_, error)
				assert.Equal(t_, expectedStatement, result)
			})
		}
	}
}

func AssertCorrectErrorHandlingForInvalidTypescriotImportStatementWrite(
	t *testing.T, description string, statement typescript.TypescriptImportStatement, expectedErrorDetail string) {

	expectedError := lp.ImportStatementWriteError{
		Statement: statement,
		Detail:    expectedErrorDetail,
	}

	resultingImportStatement, resultingError := statement.AsText()

	t.Run(description, func(t_ *testing.T) {
		assert.NotNil(t_, resultingError)
		assert.Equal(t_, "", resultingImportStatement)
		assert.Equal(t_, &expectedError, resultingError)
	})
}

func Test_TypescriptImportStatement_AsText_NoDependencies_Error(t *testing.T) {
	const EXPECTED_ERROR_DETAIL = "No dependencies were specified"

	for _, moduleCase := range moduleAsStatementStringCases() {
		description := moduleCase.Description
		testStatement := typescript.TypescriptImportStatement{
			ModulePathParts: moduleCase.ModulePathParts,
			External:        false,
			Dependecies:     []typescript.TypescriptDependecy{},
		}
		AssertCorrectErrorHandlingForInvalidTypescriotImportStatementWrite(t, description, testStatement, EXPECTED_ERROR_DETAIL)
	}
}

func Test_TypescriptImportStatement_AsText_MultipleDefaultDependencies_Error(t *testing.T) {
	const EXPECTED_ERROR_DETAIL = "Multiple default dependencies are not allowed"

	for _, moduleCase := range moduleAsStatementStringCases() {
		description := moduleCase.Description
		testStatement := typescript.TypescriptImportStatement{
			ModulePathParts: moduleCase.ModulePathParts,
			External:        false,
			Dependecies: []typescript.TypescriptDependecy{
				{Name: "something", IsDefault: true},
				{Name: "anotherThing", IsDefault: true},
			},
		}
		AssertCorrectErrorHandlingForInvalidTypescriotImportStatementWrite(t, description, testStatement, EXPECTED_ERROR_DETAIL)
	}
}
