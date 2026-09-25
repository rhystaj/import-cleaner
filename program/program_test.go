package program_test

import (
	"importcleaner/internal/config"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/parsing"
	"importcleaner/program"
	"iter"
	"reflect"
	"testing"
)

type MockParser struct{}

func (p MockParser) StatementsInText(text string) iter.Seq[string] {
	return func(yield func(string) bool) {}
}

func (p MockParser) IsImportStatement(statementText string) bool {
	return false
}

func (p MockParser) ParseImportStatement(statementText string) (parsing.ImportStatement, error) {
	return parsing.ImportStatement{}, nil
}

func (p MockParser) StatementAsString(statement parsing.ImportStatement) string {
	return ""
}

type MockFileProcessor struct {
	processedFilepaths []string
}

func (fp *MockFileProcessor) ProcessFile(filePath string, parser parsing.ImportStatementParser) error {
	fp.processedFilepaths = append(fp.processedFilepaths, filePath)
	return nil
}

func TestCleanImports_NoFilesIgnored(t *testing.T) {

	config := config.Config{}
	filepath := "test"
	fileManager := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSDir(filepath, []*iowrappers.MockFSNode{
			iowrappers.CreateMockFSFile("Readme.md", ""),
			iowrappers.CreateMockFSFile("main.py", ""),
			iowrappers.CreateMockFSDir("subdir", []*iowrappers.MockFSNode{
				iowrappers.CreateMockFSFile("some_code.py", ""),
				iowrappers.CreateMockFSFile("not_python.c", ""),
				iowrappers.CreateMockFSFile("some_more_code.py", ""),
				iowrappers.CreateMockFSDir("subsubdir", []*iowrappers.MockFSNode{
					iowrappers.CreateMockFSFile("dont.py", ""),
					iowrappers.CreateMockFSFile("do.py", ""),
				}),
			}),
			iowrappers.CreateMockFSDir("sys", []*iowrappers.MockFSNode{
				iowrappers.CreateMockFSFile("dangerous_code.py", ""),
				iowrappers.CreateMockFSFile("not_python_either.c", ""),
				iowrappers.CreateMockFSFile("really_dangerous_code.py", ""),
			}),
		}),
	})
	parser := MockParser{}
	fileProcessor := MockFileProcessor{}

	expectedFilesProcessed := []string{
		"test/main.py",
		"test/sys/dangerous_code.py",
		"test/sys/really_dangerous_code.py",
		"test/subdir/some_code.py",
		"test/subdir/some_more_code.py",
		"test/subdir/subsubdir/dont.py",
		"test/subdir/subsubdir/do.py",
	}

	program.CleanImports(config, filepath, fileManager, parser, &fileProcessor)

	if !reflect.DeepEqual(expectedFilesProcessed, fileProcessor.processedFilepaths) {
		t.Errorf("Incorrect files processed - Expected %+v, but recieved %+v",
			expectedFilesProcessed, fileProcessor.processedFilepaths)
	}
}

func TestCleanImports_FilesIgnored(t *testing.T) {

	config := config.Config{
		IgnoredPaths: []string{
			"sys",
			"subdir/subsubdir/dont.py",
		},
	}
	filepath := "test"
	fileManager := iowrappers.InitialiseMockFileManager([]*iowrappers.MockFSNode{
		iowrappers.CreateMockFSDir(filepath, []*iowrappers.MockFSNode{
			iowrappers.CreateMockFSFile("Readme.md", ""),
			iowrappers.CreateMockFSFile("main.py", ""),
			iowrappers.CreateMockFSDir("subdir", []*iowrappers.MockFSNode{
				iowrappers.CreateMockFSFile("some_code.py", ""),
				iowrappers.CreateMockFSFile("not_python.c", ""),
				iowrappers.CreateMockFSFile("some_more_code.py", ""),
				iowrappers.CreateMockFSDir("subsubdir", []*iowrappers.MockFSNode{
					iowrappers.CreateMockFSFile("dont.py", ""),
					iowrappers.CreateMockFSFile("do.py", ""),
				}),
			}),
			iowrappers.CreateMockFSDir("sys", []*iowrappers.MockFSNode{
				iowrappers.CreateMockFSFile("dangerous_code.py", ""),
				iowrappers.CreateMockFSFile("not_python_either.c", ""),
				iowrappers.CreateMockFSFile("really_dangerous_code.py", ""),
			}),
		}),
	})
	parser := MockParser{}
	fileProcessor := MockFileProcessor{}

	expectedFilesProcessed := []string{
		"test/main.py",
		"test/subdir/some_code.py",
		"test/subdir/some_more_code.py",
		"test/subdir/subsubdir/do.py",
	}

	program.CleanImports(config, filepath, fileManager, parser, &fileProcessor)

	if !reflect.DeepEqual(expectedFilesProcessed, fileProcessor.processedFilepaths) {
		t.Errorf("Incorrect files processed - Expected %+v, but recieved %+v",
			expectedFilesProcessed, fileProcessor.processedFilepaths)
	}
}
