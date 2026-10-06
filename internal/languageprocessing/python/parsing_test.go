package python_test

import (
	"errors"
	"importcleaner/internal/iowrappers"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/languageprocessing/python"
	"reflect"
	"testing"
)

func TestStatementsInText_Python(t *testing.T) {
	testParser := python.PythonImportStatementParser{}

	testText := `from something import data
from big import (
	this,
	andthis,
	andalsothis
)

text = "some text" + /
	"some text on a new line"`

	expectedStatements := []string{
		"from something import data\n",
		"from big import (\n\tthis,\n\tandthis,\n\tandalsothis\n)\n",
		"\n",
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
	testParser := python.PythonImportStatementParser{}

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
		if !testParser.IsIntendedImportStatement(statement) {
			t.Fail()
		}
	}

	for _, statement := range invalidStatements {
		if testParser.IsIntendedImportStatement(statement) {
			t.Fail()
		}
	}
}

func AssertPythonImportStatementParsedCorrectly(t *testing.T, statement string, directoryReader iowrappers.FileManager, expectedResult python.PythonImportStatement) {
	testParser := python.PythonImportStatementParser{
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
	testDirectoryReader := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile("modulefile.py", ""),
		iowrappers.CreateMockFSDir("moduledir", make([]*iowrappers.MockFSNode, 0)),
	})

	AssertPythonImportStatementParsedCorrectly(t, "import acoopythonackage", testDirectoryReader, python.PythonImportStatement{
		RawText:         "import acoopythonackage",
		ModulePathParts: []string{"acoopythonackage"},
		Dependecies:     []python.PythonDependecy{},
		External:        true,
	})

	AssertPythonImportStatementParsedCorrectly(t, "import modulefile", testDirectoryReader, python.PythonImportStatement{
		RawText:         "import modulefile",
		ModulePathParts: []string{"modulefile"},
		Dependecies:     []python.PythonDependecy{},
		External:        false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "import moduledir.subdep", testDirectoryReader, python.PythonImportStatement{
		RawText:         "import moduledir.subdep",
		ModulePathParts: []string{"moduledir", "subdep"},
		Dependecies:     []python.PythonDependecy{},
		External:        false,
	})
}

func TestParseImportStatement_ValidRootImports_Aliases(t *testing.T) {
	testDirectoryReader := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile("modulefile.py", ""),
		iowrappers.CreateMockFSDir("moduledir", make([]*iowrappers.MockFSNode, 0)),
	})

	AssertPythonImportStatementParsedCorrectly(t, "import acoopythonackage as acp", testDirectoryReader, python.PythonImportStatement{
		RawText:         "import acoopythonackage as acp",
		ModulePathParts: []string{"acoopythonackage"},
		ModuleAlias:     "acp",
		Dependecies:     []python.PythonDependecy{},
		External:        true,
	})

	AssertPythonImportStatementParsedCorrectly(t, "import modulefile as mf", testDirectoryReader, python.PythonImportStatement{
		RawText:         "import modulefile as mf",
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "mf",
		Dependecies:     []python.PythonDependecy{},
		External:        false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "import moduledir.subdep as md", testDirectoryReader, python.PythonImportStatement{
		RawText:         "import moduledir.subdep as md",
		ModulePathParts: []string{"moduledir", "subdep"},
		ModuleAlias:     "md",
		Dependecies:     []python.PythonDependecy{},
		External:        false,
	})
}

func TestParseImportStatement_ValidSelectiveImports_SingleLine(t *testing.T) {
	testDirectoryReader := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile("modulefile.py", ""),
		iowrappers.CreateMockFSDir("moduledir", make([]*iowrappers.MockFSNode, 0)),
	})

	AssertPythonImportStatementParsedCorrectly(t, "from acoopythonackage import somedependency", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from acoopythonackage import somedependency",
		ModulePathParts: []string{"acoopythonackage"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "somedependency", Alias: ""},
		},
		External: true,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from modulefile import somedependency", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from modulefile import somedependency",
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "somedependency", Alias: ""},
		},
		External: false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from moduledir.subdep import somedependency", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from moduledir.subdep import somedependency",
		ModulePathParts: []string{"moduledir", "subdep"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "somedependency", Alias: ""},
		},
		External: false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from acoopythonackage import somedependency as sd", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from acoopythonackage import somedependency as sd",
		ModulePathParts: []string{"acoopythonackage"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "somedependency", Alias: "sd"},
		},
		External: true,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from modulefile import somedependency as sd", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from modulefile import somedependency as sd",
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "somedependency", Alias: "sd"},
		},
		External: false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from moduledir.subdep import somedependency as sd", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from moduledir.subdep import somedependency as sd",
		ModulePathParts: []string{"moduledir", "subdep"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "somedependency", Alias: "sd"},
		},
		External: false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from acoopythonackage import somedependency, anotherthing as at, somethingelse", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from acoopythonackage import somedependency, anotherthing as at, somethingelse",
		ModulePathParts: []string{"acoopythonackage"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "somedependency", Alias: ""},
			{Name: "anotherthing", Alias: "at"},
			{Name: "somethingelse", Alias: ""},
		},
		External: true,
	})
}

func TestParseImportStatement_ValidSelectiveImports_MultipleLines(t *testing.T) {
	testDirectoryReader := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile("modulefile.py", ""),
		iowrappers.CreateMockFSDir("moduledir", make([]*iowrappers.MockFSNode, 0)),
	})

	AssertPythonImportStatementParsedCorrectly(t, "from acoopythonackage import (\n\tthis as t,\n\tandthis,\n\tandalsothis\n)", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from acoopythonackage import (\n\tthis as t,\n\tandthis,\n\tandalsothis\n)",
		ModulePathParts: []string{"acoopythonackage"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "this", Alias: "t"},
			{Name: "andthis", Alias: ""},
			{Name: "andalsothis", Alias: ""},
		},
		External: true,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from modulefile import (\n\tthis,\n\tandthis as at,\n\tandalsothis\n)", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from modulefile import (\n\tthis,\n\tandthis as at,\n\tandalsothis\n)",
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "this", Alias: ""},
			{Name: "andthis", Alias: "at"},
			{Name: "andalsothis", Alias: ""},
		},
		External: false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from moduledir import (\n   this,\n   andthis,\n   andalsothis as aat\n)", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from moduledir import (\n   this,\n   andthis,\n   andalsothis as aat\n)",
		ModulePathParts: []string{"moduledir"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "this", Alias: ""},
			{Name: "andthis", Alias: ""},
			{Name: "andalsothis", Alias: "aat"},
		},
		External: false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from acoopythonackage import (\n   this as t,\n   andthis,\n   andalsothis\n)", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from acoopythonackage import (\n   this as t,\n   andthis,\n   andalsothis\n)",
		ModulePathParts: []string{"acoopythonackage"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "this", Alias: "t"},
			{Name: "andthis", Alias: ""},
			{Name: "andalsothis", Alias: ""},
		},
		External: true,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from modulefile import (\n   this,\n   andthis as at,\n   andalsothis\n)", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from modulefile import (\n   this,\n   andthis as at,\n   andalsothis\n)",
		ModulePathParts: []string{"modulefile"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "this", Alias: ""},
			{Name: "andthis", Alias: "at"},
			{Name: "andalsothis", Alias: ""},
		},
		External: false,
	})

	AssertPythonImportStatementParsedCorrectly(t, "from moduledir import (\n\tthis,\n\tandthis,\n\tandalsothis as aat\n)", testDirectoryReader, python.PythonImportStatement{
		RawText:         "from moduledir import (\n\tthis,\n\tandthis,\n\tandalsothis as aat\n)",
		ModulePathParts: []string{"moduledir"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "this", Alias: ""},
			{Name: "andthis", Alias: ""},
			{Name: "andalsothis", Alias: "aat"},
		},
		External: false,
	})
}

func AssertParseErrorHandledCorrectly(t *testing.T, statement string, expectedError error, parser python.PythonImportStatementParser) {
	expectedImportStatement := python.PythonImportStatement{
		ModulePathParts: []string{},
		ModuleAlias:     "",
		Dependecies:     []python.PythonDependecy{},
		External:        false,
	}

	result, error := parser.ParseImportStatement(statement)

	if !reflect.DeepEqual(result, expectedImportStatement) {
		t.Errorf("Failure in handling of invalid import statement '%s': %+v was returned but '%+v' was expected", statement, result, expectedImportStatement)
	}

	if !errors.Is(error, expectedError) {
		t.Errorf("Incorrect error returned in python of statement '%s': '%+v' was returned, but '%+v' was expected.", statement, error, expectedError)
	}
}

func TestParseImportStatement_NotImportStatement(t *testing.T) {
	testDirectoryReader := iowrappers.InitialiseMockFileManager(make([]*iowrappers.MockFSNode, 0))

	testParser := python.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: testDirectoryReader,
	}

	testStatement := "not at all valid"
	expectedError := lp.ImportStatementParseError{
		Statement: testStatement,
		Detail:    "Not import statement",
	}

	AssertParseErrorHandledCorrectly(t, testStatement, &expectedError, testParser)
}

func TestParseImportStatement_Malformed(t *testing.T) {
	testDirectoryReader := iowrappers.InitialiseMockFileManager(make([]*iowrappers.MockFSNode, 0))

	testParser := python.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: testDirectoryReader,
	}

	AssertParseErrorHandledCorrectly(t, "import", &lp.ImportStatementParseError{Statement: "import", Detail: "No module"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something", &lp.ImportStatementParseError{Statement: "from something", Detail: "'import' keyword expected"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something import", &lp.ImportStatementParseError{Statement: "from something import", Detail: "No dependecies listed"}, testParser)
}

func TestParseImportStatement_MalformedAliases(t *testing.T) {
	testDirectoryReader := iowrappers.InitialiseMockFileManager(make([]*iowrappers.MockFSNode, 0))

	testParser := python.PythonImportStatementParser{
		WorkingDir:  "", //Doesn't matter for tesing purposes
		FileManager: testDirectoryReader,
	}

	AssertParseErrorHandledCorrectly(t, "import mod as", &lp.ImportStatementParseError{Statement: "import mod as", Detail: "Module alias expected"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something import dependency, baddep as", &lp.ImportStatementParseError{Statement: "from something import dependency, baddep as", Detail: "Dependency alias expected"}, testParser)
	AssertParseErrorHandledCorrectly(t, "from something import dependency, bad dep", &lp.ImportStatementParseError{Statement: "from something import dependency, bad dep", Detail: "Malformed dependecy 'bad dep'"}, testParser)
}

func AssertStatementCorrectlyConvertedToText(t *testing.T, statement python.PythonImportStatement, expectedString string) {
	result, _ := statement.AsText()

	if result != expectedString {
		t.Errorf("Failure to convert statement %+v to string: Expected '%s', but retrieved '%s'", statement, expectedString, result)
	}
}

func Test_PythonImportStatement_AsText(t *testing.T) {
	AssertStatementCorrectlyConvertedToText(t, python.PythonImportStatement{
		ModulePathParts: []string{"somemod"},
		ModuleAlias:     "",
		Dependecies:     []python.PythonDependecy{},
		External:        true,
	}, "import somemod")

	AssertStatementCorrectlyConvertedToText(t, python.PythonImportStatement{
		ModulePathParts: []string{"bigmod"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "something", Alias: ""},
			{Name: "thingwithalias", Alias: "twa"},
			{Name: "anotherthing", Alias: ""},
		},
		External: true,
	}, "from bigmod import something, thingwithalias as twa, anotherthing")

	AssertStatementCorrectlyConvertedToText(t, python.PythonImportStatement{
		ModulePathParts: []string{"mod", "with", "multiple", "parts"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "something", Alias: "s"},
			{Name: "thingwithalias", Alias: ""},
			{Name: "anotherthing", Alias: "at"},
		},
		External: true,
	}, "from mod.with.multiple.parts import something as s, thingwithalias, anotherthing as at")
}
