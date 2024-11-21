package parsing_test

import (
	"errors"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/parsing"
	"reflect"
	"testing"
)

func TestStatementsInText(t *testing.T) {
	testParser := parsing.PythonImportStatementParser{}

	testText := `from something import data
from big import (
	this,
	andthis,
	andalsothis
)

text = "some text" + /
	"some text on a new line"`

	expectedStatements := []string{
		"from something import data",
		"from big import (\n\tthis,\n\tandthis,\n\tandalsothis\n)",
		"",
		"text = \"some text\" + /\n\t\"some text on a new line\"",
	}

	result := testParser.StatementsInText(testText)

	resultingStatements := make([]string, 0)
	for statement := range result {
		resultingStatements = append(resultingStatements, statement)
	}

	if len(resultingStatements) != len(expectedStatements) {
		t.Errorf("Resulting and expected statements are not the same - %d statements were expected, but %d were retrieved.",
			len(expectedStatements), len(resultingStatements))
		return
	}

	for i, expectedStatement := range expectedStatements {
		if expectedStatement != resultingStatements[i] {
			t.Errorf("Unexpected element in index %d of result - expected '%s', but got '%s'", i, expectedStatement, resultingStatements[i])
		}
	}
}

func TestIsImportStatement(t *testing.T) {
	testParser := parsing.PythonImportStatementParser{}

	validStatements := []string{
		"import item",
		"from somepackage import thing1, thing2",
		"from big import (\n\tthis,\n\tandthis,\n\tandalsothis\n)",
	}

	invalidStatements := []string{
		"not valid",
		"from bad",
	}

	for _, statement := range validStatements {
		if !testParser.IsImportStatement(statement) {
			t.Fail()
		}
	}

	for _, statement := range invalidStatements {
		if testParser.IsImportStatement(statement) {
			t.Fail()
		}
	}
}

func AssertStatementParsedCorrectly(t *testing.T, statement string, directoryReader iowrappers.FileManager, expectedResult parsing.ImportStatement) {
	testParser := parsing.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: directoryReader,
	}

	resultingStatement, resultingError := testParser.ParseImportStatement(statement)

	if resultingError != nil {
		t.Errorf("Statement '%s' produced error '%s'", statement, resultingError.Error())
	}

	if !reflect.DeepEqual(resultingStatement, expectedResult) {
		t.Errorf("Statement '%s' produced result %+v, but %+v was expected", statement, resultingStatement, expectedResult)
	}
}

func TestParseImportStatement_ValidRootImports_NoAliases(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{
			{
				ItemName: "modulefile.py",
				IsDir:    false,
			},
			{
				ItemName: "moduledir",
				IsDir:    true,
			},
		},
	}

	AssertStatementParsedCorrectly(t, "import acoolpackage", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"acoolpackage"},
		Dependecies:     []parsing.Dependecy{},
		External:        true,
	})

	AssertStatementParsedCorrectly(t, "import modulefile", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"modulefile"},
		Dependecies:     []parsing.Dependecy{},
		External:        false,
	})

	AssertStatementParsedCorrectly(t, "import moduledir.subdep", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"moduledir", "subdep"},
		Dependecies:     []parsing.Dependecy{},
		External:        false,
	})
}

func TestParseImportStatement_ValidRootImports_Aliases(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{
			{
				ItemName: "modulefile.py",
				IsDir:    false,
			},
			{
				ItemName: "moduledir",
				IsDir:    true,
			},
		},
	}

	AssertStatementParsedCorrectly(t, "import acoolpackage as acp", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"acoolpackage"},
		ModuleAlias:     "acp",
		Dependecies:     []parsing.Dependecy{},
		External:        true,
	})

	AssertStatementParsedCorrectly(t, "import modulefile as mf", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "mf",
		Dependecies:     []parsing.Dependecy{},
		External:        false,
	})

	AssertStatementParsedCorrectly(t, "import moduledir.subdep as md", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"moduledir", "subdep"},
		ModuleAlias:     "md",
		Dependecies:     []parsing.Dependecy{},
		External:        false,
	})
}

func TestParseImportStatement_ValidSelectiveImports_SingleLine(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{
			{
				ItemName: "modulefile.py",
				IsDir:    false,
			},
			{
				ItemName: "moduledir",
				IsDir:    true,
			},
		},
	}

	AssertStatementParsedCorrectly(t, "from acoolpackage import somedependency", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"acoolpackage"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "somedependency", Alias: ""},
		},
		External: true,
	})

	AssertStatementParsedCorrectly(t, "from modulefile import somedependency", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "somedependency", Alias: ""},
		},
		External: false,
	})

	AssertStatementParsedCorrectly(t, "from moduledir.subdep import somedependency", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"moduledir", "subdep"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "somedependency", Alias: ""},
		},
		External: false,
	})

	AssertStatementParsedCorrectly(t, "from acoolpackage import somedependency as sd", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"acoolpackage"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "somedependency", Alias: "sd"},
		},
		External: true,
	})

	AssertStatementParsedCorrectly(t, "from modulefile import somedependency as sd", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "somedependency", Alias: "sd"},
		},
		External: false,
	})

	AssertStatementParsedCorrectly(t, "from moduledir.subdep import somedependency as sd", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"moduledir", "subdep"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "somedependency", Alias: "sd"},
		},
		External: false,
	})

	AssertStatementParsedCorrectly(t, "from acoolpackage import somedependency, anotherthing as at, somethingelse", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"acoolpackage"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "somedependency", Alias: ""},
			{Name: "anotherthing", Alias: "at"},
			{Name: "somethingelse", Alias: ""},
		},
		External: true,
	})
}

func TestParseImportStatement_ValidSelectiveImports_MultipleLines(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{
			{
				ItemName: "modulefile.py",
				IsDir:    false,
			},
			{
				ItemName: "moduledir",
				IsDir:    true,
			},
		},
	}

	AssertStatementParsedCorrectly(t, "from acoolpackage import (\n\tthis as t,\n\tandthis,\n\tandalsothis\n)", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"acoolpackage"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "this", Alias: "t"},
			{Name: "andthis", Alias: ""},
			{Name: "andalsothis", Alias: ""},
		},
		External: true,
	})

	AssertStatementParsedCorrectly(t, "from modulefile import (\n\tthis,\n\tandthis as at,\n\tandalsothis\n)", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "this", Alias: ""},
			{Name: "andthis", Alias: "at"},
			{Name: "andalsothis", Alias: ""},
		},
		External: false,
	})

	AssertStatementParsedCorrectly(t, "from moduledir import (\n\tthis,\n\tandthis,\n\tandalsothis as aat\n)", testDirectoryReader, parsing.ImportStatement{
		ModulePathParts: []string{"moduledir"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "this", Alias: ""},
			{Name: "andthis", Alias: ""},
			{Name: "andalsothis", Alias: "aat"},
		},
		External: false,
	})
}

func AssertParseErrorHandledCorrectly(t *testing.T, statement string, expectedError error, parser parsing.ImportStatementParser) {
	expectedImportStatement := parsing.ImportStatement{
		ModulePathParts: []string{},
		ModuleAlias:     "",
		Dependecies:     []parsing.Dependecy{},
		External:        false,
	}

	result, error := parser.ParseImportStatement(statement)

	if !reflect.DeepEqual(result, expectedImportStatement) {
		t.Errorf("Failure in handling of invalid import statement '%s': %+v was returned but '%+v' was expected", statement, result, expectedImportStatement)
	}

	if !errors.Is(error, expectedError) {
		t.Errorf("Incorrect error returned in parsing of statement '%s': '%+v' was returned, but '%+v' was expected.", statement, error, expectedError)
	}
}

func TestParseImportStatement_NotImportStatement(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{},
	}

	testParser := parsing.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: testDirectoryReader,
	}

	testStatement := "not at all valid"
	expectedError := parsing.ImportStatementParseError{
		Statement: testStatement,
		Detail:    "Not import statement",
	}

	AssertParseErrorHandledCorrectly(t, testStatement, &expectedError, testParser)
}

func TestParseImportStatement_Malformed(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{},
	}

	testParser := parsing.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: testDirectoryReader,
	}

	AssertParseErrorHandledCorrectly(t, "import", &parsing.ImportStatementParseError{Statement: "import", Detail: "No module"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something", &parsing.ImportStatementParseError{Statement: "from something", Detail: "'import' keyword expected"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something import", &parsing.ImportStatementParseError{Statement: "from something import", Detail: "No dependecies listed"}, testParser)
}

func TestParseImportStatement_MalformedAliases(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{},
	}

	testParser := parsing.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: testDirectoryReader,
	}

	AssertParseErrorHandledCorrectly(t, "import mod as", &parsing.ImportStatementParseError{Statement: "import mod as", Detail: "Module alias expected"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something import dependency, baddep as", &parsing.ImportStatementParseError{Statement: "from something import dependency, baddep as", Detail: "Dependency alias expected"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something import dependency, bad dep", &parsing.ImportStatementParseError{Statement: "from something import dependency, bad dep", Detail: "Malformed dependecy 'bad dep'"}, testParser)
}

func AssertStatementCorrectlyConvertedToString(t *testing.T, statement parsing.ImportStatement, expectedString string, parser parsing.ImportStatementParser) {
	result := parser.StatementAsString(statement)

	if result != expectedString {
		t.Errorf("Failure to convert statement %+v to string: Expected '%s', but retrieved '%s'", statement, expectedString, result)
	}
}

func TestStatementAsString(t *testing.T) {
	testDirectoryReader := iowrappers.MockFileManager{
		Items: []iowrappers.DirectoryItemInfo{},
	}

	testParser := parsing.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: testDirectoryReader,
	}

	AssertStatementCorrectlyConvertedToString(t, parsing.ImportStatement{
		ModulePathParts: []string{"somemod"},
		ModuleAlias:     "",
		Dependecies:     []parsing.Dependecy{},
		External:        true,
	}, "import somemod", testParser)

	AssertStatementCorrectlyConvertedToString(t, parsing.ImportStatement{
		ModulePathParts: []string{"bigmod"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "something", Alias: ""},
			{Name: "thingwithalias", Alias: "twa"},
			{Name: "anotherthing", Alias: ""},
		},
		External: true,
	}, "from bigmod import something, thingwithalias as twa, anotherthing", testParser)

	AssertStatementCorrectlyConvertedToString(t, parsing.ImportStatement{
		ModulePathParts: []string{"mod", "with", "multiple", "parts"},
		ModuleAlias:     "",
		Dependecies: []parsing.Dependecy{
			{Name: "something", Alias: "s"},
			{Name: "thingwithalias", Alias: ""},
			{Name: "anotherthing", Alias: "at"},
		},
		External: true,
	}, "from mod.with.multiple.parts import something as s, thingwithalias, anotherthing as at", testParser)
}
