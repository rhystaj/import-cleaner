package typescript_test

import (
	"fmt"
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

func Test_TypescriptImportStatement_GetGenericDetails(t *testing.T) {
	modulePathCases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"some", "module", "path"},
			expected: []string{"some", "module", "path"},
		},
		{
			input:    []string{"..", "some", "module", "path"},
			expected: []string{"some", "module", "path"},
		},
		{
			input:    []string{".", "some", "module", "path"},
			expected: []string{"some", "module", "path"},
		},
		{
			input:    []string{"..", ".", "some", "module", "path"},
			expected: []string{"some", "module", "path"},
		},
		{
			input:    []string{"..", ".", "..", "some", "module", "path"},
			expected: []string{"some", "module", "path"},
		},
		{
			input:    []string{"..", "..", "."},
			expected: []string{},
		},
	}

	for _, modulePathCase := range modulePathCases {

		testStatement := typescript.TypescriptImportStatement{
			External:        false,
			ModulePathParts: modulePathCase.input,
			Dependecies: []typescript.TypescriptDependecy{
				{
					Name:         "dependency",
					Alias:        "dep",
					IsTypeImport: true,
				},
			},
		}

		expectedGenericDetails := lp.ImportStatementGenericDetails{
			External:        false,
			ModulePathParts: modulePathCase.expected,
			Dependecies: []lp.GenericDependecy{
				{
					Name:  "dependency",
					Alias: "dep",
				},
			},
		}

		assert.Equal(t, expectedGenericDetails, testStatement.GetGenericDetails())
	}
}

func Test_TypescriptImportStatement_ConsolidationKey(t *testing.T) {
	cases := []struct {
		input    typescript.TypescriptImportStatement
		expected string
	}{
		{
			input: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
			},
			expected: "some.module.path",
		},
		{
			input: typescript.TypescriptImportStatement{
				ModulePathParts: []string{".", "some", "module", "path"},
			},
			expected: "some.module.path",
		},
		{
			input: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"..", "..", "some", "module", "path"},
			},
			expected: "some.module.path",
		},
		{
			input: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				IsTypeImport:    true,
			},
			expected: "*some.module.path",
		},
		{
			input: typescript.TypescriptImportStatement{
				ModulePathParts: []string{".", "some", "module", "path"},
				IsTypeImport:    true,
			},
			expected: "*some.module.path",
		},
		{
			input: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"..", "..", "some", "module", "path"},
				IsTypeImport:    true,
			},
			expected: "*some.module.path",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.input.ConsolidationKey())
	}
}

func Test_TypescriptImportStatement_ConsolidateWith_Success(t *testing.T) {
	cases := []struct {
		statementA typescript.TypescriptImportStatement
		statementB typescript.TypescriptImportStatement
		expected   typescript.TypescriptImportStatement
	}{
		{
			statementA: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies:     []typescript.TypescriptDependecy{{Name: "depA"}},
			},
			statementB: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies: []typescript.TypescriptDependecy{
					{Name: "depB", IsTypeImport: true},
					{Name: "depC"},
				},
			},
			expected: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies: []typescript.TypescriptDependecy{
					{Name: "depA"},
					{Name: "depB", IsTypeImport: true},
					{Name: "depC"},
				},
			},
		},
		{
			statementA: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies:     []typescript.TypescriptDependecy{{Name: "depA"}},
				IsTypeImport:    true,
			},
			statementB: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies: []typescript.TypescriptDependecy{
					{Name: "depB", IsTypeImport: true},
					{Name: "depC"},
				},
				IsTypeImport: true,
			},
			expected: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies: []typescript.TypescriptDependecy{
					{Name: "depA"},
					{Name: "depB", IsTypeImport: true},
					{Name: "depC"},
				},
				IsTypeImport: true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.expected.ConsolidationKey(), func(t *testing.T) {
			resultingImportStatement, resultingError := c.statementA.ConsolidateWith(c.statementB)
			assert.Nil(t, resultingError)
			assert.Equal(t, c.expected, resultingImportStatement)
		})
	}
}

type dummyImportStatement struct{}

func (d dummyImportStatement) AsText() (string, error) {
	return "dummy", nil
}

func (d dummyImportStatement) GetGenericDetails() lp.ImportStatementGenericDetails {
	return lp.ImportStatementGenericDetails{
		ModulePathParts: []string{"dummy"},
		ModuleAlias:     "d",
		Dependecies:     []lp.GenericDependecy{},
		External:        true,
	}
}

func (d dummyImportStatement) ConsolidationKey() string {
	return "dummy"
}

func (d dummyImportStatement) ConsolidateWith(other lp.ImportStatement) (lp.ImportStatement, error) {
	return d, nil
}

func (d dummyImportStatement) Compare(other lp.ImportStatement) int {
	return 0
}

func Test_TypescriptImportStatement_ConsolidateWith_Error(t *testing.T) {
	cases := []struct {
		statementA    typescript.TypescriptImportStatement
		statementB    lp.ImportStatement
		expectedError error
	}{
		{
			statementA: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies:     []typescript.TypescriptDependecy{{Name: "depA"}},
			},
			statementB: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "other", "module", "path"},
				Dependecies:     []typescript.TypescriptDependecy{{Name: "depB"}},
			},
			expectedError: fmt.Errorf("statements have different consolidation keys"),
		},
		{
			statementA: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies:     []typescript.TypescriptDependecy{{Name: "depA"}},
			},
			statementB: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies:     []typescript.TypescriptDependecy{{Name: "depB"}},
				IsTypeImport:    true,
			},
			expectedError: fmt.Errorf("statements have different consolidation keys"),
		},
		{
			statementA: typescript.TypescriptImportStatement{
				ModulePathParts: []string{"some", "module", "path"},
				Dependecies:     []typescript.TypescriptDependecy{{Name: "depA"}},
			},
			statementB:    dummyImportStatement{},
			expectedError: fmt.Errorf("other statement is not a TypescriptImportStatement"),
		},
	}

	for _, c := range cases {
		t.Run(c.expectedError.Error(), func(t *testing.T) {
			_, resultingError := c.statementA.ConsolidateWith(c.statementB)
			assert.NotNil(t, resultingError)
			assert.Equal(t, c.expectedError, resultingError)
		})
	}
}
