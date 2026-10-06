package typescript_test

import (
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/languageprocessing/typescript"
	"testing"

	"github.com/stretchr/testify/assert"
)

type inidvidualDependencyParseCase struct {
	Text       string
	Dependency typescript.TypescriptDependecy
}

type generateIndividualDependencyCasesArgs struct {
	Name                   string
	Alias                  string
	IncludeTypeKeywordCase bool
}

func generateIndividualDependencyCases(args generateIndividualDependencyCasesArgs) []inidvidualDependencyParseCase {
	cases := make([]inidvidualDependencyParseCase, 0)

	cases = append(cases, inidvidualDependencyParseCase{
		Text: args.Name,
		Dependency: typescript.TypescriptDependecy{
			Name: args.Name,
		},
	})

	if args.Alias != "" {
		cases = append(cases, inidvidualDependencyParseCase{
			Text: args.Name + " as " + args.Alias,
			Dependency: typescript.TypescriptDependecy{
				Name:  args.Name,
				Alias: args.Alias,
			},
		})
	}

	if args.IncludeTypeKeywordCase {
		cases = append(cases, inidvidualDependencyParseCase{
			Text: "type " + args.Name,
			Dependency: typescript.TypescriptDependecy{
				Name:         args.Name,
				IsTypeImport: true,
			},
		})
	}

	if args.Alias != "" && args.IncludeTypeKeywordCase {
		cases = append(cases, inidvidualDependencyParseCase{
			Text: "type " + args.Name + " as " + args.Alias,
			Dependency: typescript.TypescriptDependecy{
				Name:         args.Name,
				Alias:        args.Alias,
				IsTypeImport: true,
			},
		})
	}

	return cases
}

type generateValidDependencySetPartCasesArgs struct {
	IncludeDefaultCases    bool
	IncludeTypeImportCases bool
	IncludeNewLineCases    bool
}

func generateValidDependencySetPartCases(args generateValidDependencySetPartCasesArgs) []dependencySetCase {
	defaultCases := make([]dependencySetCase, 0)
	if args.IncludeDefaultCases {
		for _, inidividualCase := range generateIndividualDependencyCases(generateIndividualDependencyCasesArgs{
			Name:                   "something",
			Alias:                  "someCoolAlias",
			IncludeTypeKeywordCase: args.IncludeTypeImportCases,
		}) {
			defaultCases = append(defaultCases, dependencySetCase{
				Text: inidividualCase.Text,
				Dependencies: []typescript.TypescriptDependecy{
					{
						Name:         inidividualCase.Dependency.Name,
						IsDefault:    true,
						IsTypeImport: inidividualCase.Dependency.IsTypeImport,
						Alias:        inidividualCase.Dependency.Alias,
					},
				},
			})
		}
	}

	nonDefaultCases := make([]dependencySetCase, 0)
	for _, inidividualCaseA := range generateIndividualDependencyCases(generateIndividualDependencyCasesArgs{
		Name:                   "this",
		Alias:                  "t",
		IncludeTypeKeywordCase: args.IncludeTypeImportCases,
	}) {
		for _, inidividualCaseB := range generateIndividualDependencyCases(generateIndividualDependencyCasesArgs{
			Name:                   "andThis",
			Alias:                  "at",
			IncludeTypeKeywordCase: args.IncludeTypeImportCases,
		}) {
			nonDefaultCases = append(nonDefaultCases, dependencySetCase{
				Text: "{ " + inidividualCaseA.Text + ", " + inidividualCaseB.Text + " }",
				Dependencies: []typescript.TypescriptDependecy{
					inidividualCaseA.Dependency,
					inidividualCaseB.Dependency,
				},
			})
			if args.IncludeNewLineCases {
				nonDefaultCases = append(nonDefaultCases, dependencySetCase{
					Text: "{\n\t" + inidividualCaseA.Text + ",\n\t" + inidividualCaseB.Text + "\n}",
					Dependencies: []typescript.TypescriptDependecy{
						inidividualCaseA.Dependency,
						inidividualCaseB.Dependency,
					},
				})
			}
		}
	}

	mixedCases := make([]dependencySetCase, len(defaultCases)*len(nonDefaultCases))
	for defaultCaseIndex, defaultCase := range defaultCases {
		for nonDefaultCaseIndex, nonDefaultCase := range nonDefaultCases {
			index := defaultCaseIndex*len(nonDefaultCases) + nonDefaultCaseIndex
			mixedCase := dependencySetCase{
				Text:         defaultCase.Text + ", " + nonDefaultCase.Text,
				Dependencies: append(defaultCase.Dependencies, nonDefaultCase.Dependencies...),
			}
			mixedCases[index] = mixedCase
		}
	}

	return append(defaultCases, append(nonDefaultCases, mixedCases...)...)
}

func TestStatementsInText_Typescript(t *testing.T) {

	testText := `import something from 'somecoolmodule';
import {
	this,
	andthis,
	andalsothis
} from "bigimport";

const VAL = 1
import {aValue} from 'tricky/case'

if(something === 1) {
	aFunction(value)
	import * from 'zzz'

	let x = new Thing().make();
}

while(true) { console.log() }

import anotherThing from 'anothercoolmodule
`

	expectedStatements := []string{
		"import something from 'somecoolmodule';\n",
		"import {\n\tthis,\n\tandthis,\n\tandalsothis\n} from \"bigimport\";\n",
		"\nconst VAL = 1\n",
		"import {aValue} from 'tricky/case'\n",
		"\nif(something === 1) {\n\taFunction(value)\n",
		"\timport * from 'zzz'\n",
		"\n\tlet x = new Thing().make();\n}\n\nwhile(true) { console.log() }\n\n",
		"import anotherThing from 'anothercoolmodule\n",
	}

	testParser := typescript.TypescriptImportStatementParser{}

	testResult := testParser.StatementsInText(testText)
	resultingStatements := make([]string, 0)
	for statement := range testResult {
		resultingStatements = append(resultingStatements, statement)
	}

	assert.Equal(t, expectedStatements, resultingStatements)

}

func TestInIntendedImportStatement_True_Typescript(t *testing.T) {

	testParser := typescript.TypescriptImportStatementParser{}

	testValues := []string{
		"import * from 'something",
		"import {\n\tthis,\n\tandthis\n} from bigimport",
		"import sjdhnfksdhfjhk&&&", //Technically indended to be an import statement, but not valid
	}

	for _, testValue := range testValues {
		t.Run(testValue, func(t_ *testing.T) {
			assert.True(t_, testParser.IsIntendedImportStatement(testValue))
		})
	}

}

func TestIsIntendedImportStatement_False_Typescript(t *testing.T) {

	testParser := typescript.TypescriptImportStatementParser{}

	testValues := []string{
		"nope",
		"let x = 1",
		"from x import y", //Wrong language
	}

	for _, testValue := range testValues {
		t.Run(testValue, func(t_ *testing.T) {
			assert.False(t_, testParser.IsIntendedImportStatement(testValue))
		})
	}

}

type dependencySetCase struct {
	Text                string
	Dependencies        []typescript.TypescriptDependecy
	ExpectedErrorDetail string
}

type moduleParseCase struct {
	Text                    string
	ExpectedExternal        bool
	ExpectedModulePathPaths []string
}

func validModuleParseCases() []moduleParseCase {
	return []moduleParseCase{
		{
			Text:                    "from 'somecoolmodule'",
			ExpectedExternal:        true,
			ExpectedModulePathPaths: []string{"somecoolmodule"},
		},
		{
			Text:                    "from \"somecoolmodule\"",
			ExpectedExternal:        true,
			ExpectedModulePathPaths: []string{"somecoolmodule"},
		},
		{
			Text:                    "from 'grandparent/parent/child'",
			ExpectedExternal:        true,
			ExpectedModulePathPaths: []string{"grandparent", "parent", "child"},
		},
		{
			Text:                    "from \"grandparent/parent/child\"",
			ExpectedExternal:        true,
			ExpectedModulePathPaths: []string{"grandparent", "parent", "child"},
		},
		{
			Text:                    "from './sibling'",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{".", "sibling"},
		},
		{
			Text:                    "from \"./sibling\"",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{".", "sibling"},
		},
		{
			Text:                    "from '../../../cousin'",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{"..", "..", "..", "cousin"},
		},
		{
			Text:                    "from \"../../../cousin\"",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{"..", "..", "..", "cousin"},
		},
		{
			Text:                    "from './sibling/child'",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{".", "sibling", "child"},
		},
		{
			Text:                    "from \"./sibling/child\"",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{".", "sibling", "child"},
		},
		{
			Text:                    "from '../../../cousin/child'",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{"..", "..", "..", "cousin", "child"},
		},
		{
			Text:                    "from \"../../../cousin/child\"",
			ExpectedExternal:        false,
			ExpectedModulePathPaths: []string{"..", "..", "..", "cousin", "child"},
		},
	}
}

func TestParseImportStatement_ValidNonTypeImport_Typescript(t *testing.T) {

	testParser := typescript.TypescriptImportStatementParser{}

	for _, dependencyCase := range generateValidDependencySetPartCases(generateValidDependencySetPartCasesArgs{
		IncludeDefaultCases:    true,
		IncludeTypeImportCases: true,
		IncludeNewLineCases:    true,
	}) {
		for _, moduleCase := range validModuleParseCases() {
			for _, statementSuffix := range []string{";", "", "\n", ";\n"} {
				statement := "import " + dependencyCase.Text + " " + moduleCase.Text + statementSuffix

				expectedResult := typescript.TypescriptImportStatement{
					ModulePathParts: moduleCase.ExpectedModulePathPaths,
					External:        moduleCase.ExpectedExternal,
					Dependecies:     dependencyCase.Dependencies,
				}

				result, err := testParser.ParseImportStatement(statement)

				t.Run(statement, func(t_ *testing.T) {
					assert.Nil(t_, err)
					assert.Equal(t_, expectedResult, result)
				})
			}
		}
	}
}

func TestParseImportStatement_ValidTypeImport_Typescript(t *testing.T) {

	testParser := typescript.TypescriptImportStatementParser{}

	for _, dependencyCase := range generateValidDependencySetPartCases(generateValidDependencySetPartCasesArgs{
		IncludeDefaultCases:    false,
		IncludeTypeImportCases: false,
		IncludeNewLineCases:    true,
	}) {
		for _, moduleCase := range validModuleParseCases() {
			for _, statementSuffix := range []string{";", "", "\n", ";\n"} {
				statement := "import type " + dependencyCase.Text + " " + moduleCase.Text + statementSuffix

				expectedResult := typescript.TypescriptImportStatement{
					ModulePathParts: moduleCase.ExpectedModulePathPaths,
					External:        moduleCase.ExpectedExternal,
					Dependecies:     dependencyCase.Dependencies,
					IsTypeImport:    true,
				}

				result, err := testParser.ParseImportStatement(statement)

				t.Run(statement, func(t_ *testing.T) {
					assert.Nil(t_, err)
					assert.Equal(t_, expectedResult, result)
				})
			}
		}
	}
}

func AssertCorrectErrorHandlingForInvalidTypescriotImportStatement(t *testing.T, statement string, expectedError error) {

	testParser := typescript.TypescriptImportStatementParser{}

	resultingImportStatement, resultingError := testParser.ParseImportStatement(statement)

	t.Run(statement, func(t_ *testing.T) {
		assert.NotNil(t_, resultingError)
		assert.Equal(t_, typescript.TypescriptImportStatement{}, resultingImportStatement)
		assert.ErrorIs(t_, expectedError, resultingError)
	})

}

func TestParseImportStatement_InvalidStatement_NoFromKeyword_Typescript(t *testing.T) {
	const EXPECTED_ERROR_DETAIL = "Statement doesn't include 'from' keyword"

	AssertCorrectErrorHandlingForInvalidTypescriotImportStatement(t, "import", &lp.ImportStatementParseError{
		Statement: "import",
		Detail:    EXPECTED_ERROR_DETAIL,
	})

	for _, depCase := range generateValidDependencySetPartCases(generateValidDependencySetPartCasesArgs{
		IncludeDefaultCases:    true,
		IncludeTypeImportCases: true,
		IncludeNewLineCases:    true,
	}) {
		statementNoFrom := "import " + depCase.Text
		AssertCorrectErrorHandlingForInvalidTypescriotImportStatement(t, statementNoFrom, &lp.ImportStatementParseError{
			Statement: statementNoFrom,
			Detail:    EXPECTED_ERROR_DETAIL,
		})
	}
}

func TestParseImportStatement_InvalidStatement_InvalidDependencyDeclaration(t *testing.T) {
	testCases := []dependencySetCase{
		{
			Text:                "",
			ExpectedErrorDetail: "No dependencies were provided",
		},
		{
			Text:                "{ this, andThis }, something",
			ExpectedErrorDetail: "No further dependecy definitions expected after closing bracket",
		},
		{
			Text:                "{ this, andThis ",
			ExpectedErrorDetail: "Dependencies definition missing closing bracket",
		},
	}

	for _, testCase := range testCases {
		for _, moduleCase := range validModuleParseCases() {
			statement := "import " + testCase.Text + " from " + moduleCase.Text
			AssertCorrectErrorHandlingForInvalidTypescriotImportStatement(t, statement, &lp.ImportStatementParseError{
				Statement: statement,
				Detail:    testCase.ExpectedErrorDetail,
			})
		}
	}
}

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
