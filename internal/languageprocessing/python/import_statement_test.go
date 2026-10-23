package python_test

import (
	"fmt"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/languageprocessing/python"
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectedOutput struct {
	noFormatting   string
	withLineBreaks string
}

func Test_PythonImportStatement_AsText(t *testing.T) {
	lineBreakAtNumChars := 66

	testCases := []struct {
		statement      python.PythonImportStatement
		expectedOutput expectedOutput
	}{
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"somemod"},
				ModuleAlias:     "",
				Dependecies:     []python.PythonDependecy{},
				External:        true,
			},
			expectedOutput: expectedOutput{
				noFormatting:   "import somemod",
				withLineBreaks: "import somemod",
			},
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"bigmod"},
				ModuleAlias:     "",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: ""},
					{Name: "thingwithalias", Alias: "twa"},
					{Name: "anotherthing", Alias: ""},
				},
				External: true,
			},
			expectedOutput: expectedOutput{
				noFormatting:   "from bigmod import something, thingwithalias as twa, anotherthing",
				withLineBreaks: "from bigmod import something, thingwithalias as twa, anotherthing",
			},
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"mod", "with", "multiple", "parts"},
				ModuleAlias:     "",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
					{Name: "thingwithalias", Alias: ""},
					{Name: "anotherthing", Alias: "at"},
				},
				External: true,
			},
			expectedOutput: expectedOutput{
				noFormatting:   "from mod.with.multiple.parts import something as s, thingwithalias, anotherthing as at",
				withLineBreaks: "from mod.with.multiple.parts import (\n\tsomething as s,\n\tthingwithalias,\n\tanotherthing as at\n)",
			},
		},
		{
			statement: python.PythonImportStatement{
				AncestoryLevels: 1,
				ModulePathParts: []string{"somemod"},
				ModuleAlias:     "",
				Dependecies:     []python.PythonDependecy{},
				External:        false,
			},
			expectedOutput: expectedOutput{
				noFormatting:   "import .somemod",
				withLineBreaks: "import .somemod",
			},
		},
		{
			statement: python.PythonImportStatement{
				AncestoryLevels: 2,
				ModulePathParts: []string{"bigmod"},
				ModuleAlias:     "",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: ""},
					{Name: "thingwithalias", Alias: "twa"},
					{Name: "anotherthing", Alias: ""},
				},
				External: false,
			},
			expectedOutput: expectedOutput{
				noFormatting:   "from ..bigmod import something, thingwithalias as twa, anotherthing",
				withLineBreaks: "from ..bigmod import (\n\tsomething,\n\tthingwithalias as twa,\n\tanotherthing\n)",
			},
		},
		{
			statement: python.PythonImportStatement{
				AncestoryLevels: 3,
				ModulePathParts: []string{"mod", "with", "multiple", "parts"},
				ModuleAlias:     "",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
					{Name: "thingwithalias", Alias: ""},
					{Name: "anotherthing", Alias: "at"},
				},
				External: false,
			},
			expectedOutput: expectedOutput{
				noFormatting:   "from ...mod.with.multiple.parts import something as s, thingwithalias, anotherthing as at",
				withLineBreaks: "from ...mod.with.multiple.parts import (\n\tsomething as s,\n\tthingwithalias,\n\tanotherthing as at\n)",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedOutput.noFormatting+" - No formatting", func(t *testing.T) {
			result, _ := testCase.statement.AsText(lp.TextOutputOptions{})
			assert.Equal(t, testCase.expectedOutput.noFormatting, result)
		})
		t.Run(testCase.expectedOutput.noFormatting+" - Line breaks after num chars", func(t *testing.T) {
			result, _ := testCase.statement.AsText(lp.TextOutputOptions{
				BreakLinesAtNumChars: lineBreakAtNumChars,
			})
			assert.Equal(t, testCase.expectedOutput.withLineBreaks, result)
		})
	}
}

func Test_PythonImportStatement_ConsolidationKey(t *testing.T) {

	cases := []struct {
		statement   python.PythonImportStatement
		expectedKey string
	}{
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"somemod"},
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
				},
			},
			expectedKey: "somemod",
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"somemod"},
				ModuleAlias:     "sm",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
				},
			},
			expectedKey: "somemod(sm)",
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
				},
			},
			expectedKey: "some.mod.with.multiple.parts",
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
				},
			},
			expectedKey: "some.mod.with.multiple.parts(smwmp)",
		},

		// For python we need to distiguish between imports with and without dependencies as these are
		// distict cases.
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"somemod"},
			},
			expectedKey: "$somemod",
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"somemod"},
				ModuleAlias:     "sm",
			},
			expectedKey: "$somemod(sm)",
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
			},
			expectedKey: "$some.mod.with.multiple.parts",
		},
		{
			statement: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
			},
			expectedKey: "$some.mod.with.multiple.parts(smwmp)",
		},
	}

	for _, c := range cases {
		t.Run(c.expectedKey, func(t *testing.T) {
			key := c.statement.ConsolidationKey()
			assert.Equal(t, c.expectedKey, key)
		})
	}
}

func Test_PythonImportStatement_ConsolidateWith_Success(t *testing.T) {
	cases := []struct {
		statementA     python.PythonImportStatement
		statementB     python.PythonImportStatement
		expectedResult python.PythonImportStatement
	}{
		{
			statementA: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
					{Name: "anotherthing", Alias: "at"},
				},
				External: true,
			},
			statementB: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "extrathing", Alias: "et"},
					{Name: "onelastthingiswear", Alias: "oltwis"},
				},
				External: true,
			},
			expectedResult: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
					{Name: "anotherthing", Alias: "at"},
					{Name: "extrathing", Alias: "et"},
					{Name: "onelastthingiswear", Alias: "oltwis"},
				},
				External: true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.expectedResult.ConsolidationKey(), func(t *testing.T) {
			result, _ := c.statementA.ConsolidateWith(c.statementB)
			assert.Equal(t, c.expectedResult, result)
		})
	}
}

type dummyImportStatement struct{}

func (d dummyImportStatement) AsText(options lp.TextOutputOptions) (string, error) {
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

func Test_PythonImportStatement_ConsolidateWith_Failure(t *testing.T) {
	cases := []struct {
		statementA    python.PythonImportStatement
		statementB    lp.ImportStatement
		expectedError error
	}{
		{
			statementA: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
					{Name: "anotherthing", Alias: "at"},
				},
				External: true,
			},
			statementB: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				Dependecies: []python.PythonDependecy{
					{Name: "extrathing", Alias: "et"},
					{Name: "onelastthingiswear", Alias: "oltwis"},
				},
				External: true,
			},
			expectedError: fmt.Errorf("consolidation keys do not match"),
		},
		{
			statementA: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
					{Name: "anotherthing", Alias: "at"},
				},
				External: true,
			},
			statementB:    dummyImportStatement{},
			expectedError: fmt.Errorf("other is not a Python import statement"),
		},
		{
			statementA: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "something", Alias: "s"},
					{Name: "anotherthing", Alias: "at"},
				},
				External: true,
			},
			statementB: python.PythonImportStatement{
				ModulePathParts: []string{"some", "mod", "with", "multiple", "parts"},
				ModuleAlias:     "smwmp",
				Dependecies: []python.PythonDependecy{
					{Name: "extrathing", Alias: "et"},
					{Name: "onelastthingiswear", Alias: "oltwis"},
				},
				External: false,
			},
			expectedError: fmt.Errorf("externalities do not match"),
		},
	}

	for _, c := range cases {
		t.Run(c.expectedError.Error(), func(t *testing.T) {
			_, err := c.statementA.ConsolidateWith(c.statementB)
			assert.Equal(t, c.expectedError, err)
		})
	}
}
