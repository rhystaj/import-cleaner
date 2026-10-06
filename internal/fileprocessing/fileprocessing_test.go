package fileprocessing_test

import (
	"errors"
	"importcleaner/internal/config"
	fileprocessing "importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
	"testing"
)

func TestProcessFile_ValidPythonFileWithImports(t *testing.T) {
	testFileName := "test_file.py"

	testInputFileContents := `from internal import VALUE

x = VALUE

from dataclasses import dataclass, field
import datetime

class SomeClass:
    def __init__(self) -> None:
		from big import (
			this,
			andthis,
			andalsothis
		)
        pass

from submod import SUBVALUE
`

	expectedOutputFileContents := `from big import this, andthis, andalsothis
from dataclasses import dataclass, field
import datetime

from internal import VALUE
from submod import SUBVALUE


x = VALUE


class SomeClass:
    def __init__(self) -> None:
        pass

`

	fileManager := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile(testFileName, testInputFileContents),
		iowrappers.CreateMockFSFile("submod.py", ""),
		iowrappers.CreateMockFSDir("internal", make([]*iowrappers.MockFSNode, 0)),
	})

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeExternal,
			}),
		},
	}

	parser := parsing.PythonImportStatementParser{
		WorkingDir:  "",
		FileManager: fileManager,
	}

	testFileProcessor := fileprocessing.FileProcessorImpl{
		Parser:      parser,
		FileManager: fileManager,
		Organiser:   testOrganiser,
	}

	error := testFileProcessor.ProcessFile(testFileName)

	if error != nil {
		t.Errorf("File processing returned error: %s", error.Error())
	}

	updatedFileContents, _ := fileManager.ReadStringFromFile(testFileName)
	if updatedFileContents != expectedOutputFileContents {
		t.Errorf("File contents was '%s', \nbut '%s' \nwas expected.", updatedFileContents, expectedOutputFileContents)
	}
}

func TestProcessFile_ValidPythonFileWithNoImports(t *testing.T) {
	testFileName := "test_file.py"

	testInputFileContents := `x = VALUE

class SomeClass:
    def __init__(self) -> None:
        pass
`

	fileManager := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile(testFileName, testInputFileContents),
		iowrappers.CreateMockFSFile("submod.py", ""),
		iowrappers.CreateMockFSDir("internal", make([]*iowrappers.MockFSNode, 0)),
	})

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeExternal,
			}),
		},
	}

	parser := parsing.PythonImportStatementParser{
		WorkingDir:  "",
		FileManager: fileManager,
	}

	testFileProcessor := fileprocessing.FileProcessorImpl{
		Parser:      parser,
		FileManager: fileManager,
		Organiser:   testOrganiser,
	}

	error := testFileProcessor.ProcessFile(testFileName)

	if error != nil {
		t.Errorf("File processing returned error: %s", error.Error())
	}

	updatedFileContents, _ := fileManager.ReadStringFromFile(testFileName)
	if updatedFileContents != testInputFileContents {
		t.Errorf("File with no import should remain unchanged, but output was '%s'", updatedFileContents)
	}
}

func TestProcessFile_InvalidPythonFile(t *testing.T) {
	testFileName := "test_file.py"

	testInputFileContents := `from internal import VALUE

x = VALUE

from dataclasses import dataclass, field
import datetime as

class SomeClass:
    def __init__(self) -> None:
		from big import (
			this,
			andthis,
			andalsothis
		)
        pass

from submod import SUBVALUE
`

	expectedError := parsing.ImportStatementParseError{
		Statement: "import datetime as\n",
		Detail:    "Module alias expected",
	}

	fileManager := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSFile(testFileName, testInputFileContents),
		iowrappers.CreateMockFSFile("submod.py", ""),
		iowrappers.CreateMockFSDir("internal", make([]*iowrappers.MockFSNode, 0)),
	})

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeExternal,
			}),
		},
	}

	parser := parsing.PythonImportStatementParser{
		WorkingDir:  "",
		FileManager: fileManager,
	}

	testFileProcessor := fileprocessing.FileProcessorImpl{
		Parser:      parser,
		FileManager: fileManager,
		Organiser:   testOrganiser,
	}

	error := testFileProcessor.ProcessFile(testFileName)
	if !errors.Is(error, &expectedError) {
		t.Errorf("File processing expected error '%+v', but returned error '%+v'", expectedError, error)
	}

	updatedFileContents, _ := fileManager.ReadStringFromFile(testFileName)
	if updatedFileContents != testInputFileContents {
		t.Errorf("The file contents should not have been modified.")
	}
}
