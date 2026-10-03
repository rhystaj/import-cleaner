package program_test

import (
	"importcleaner/internal/config"
	fileprocessing "importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
	"importcleaner/program"
	"testing"

	"github.com/ovechkin-dm/mockio/v2/mock"
)

func TestCleanImports_NoFilesIgnored(t *testing.T) {

	config := config.Config{
		Language: "python",
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

	mockController := mock.NewMockController(t)
	mockFileProcessor := mock.Mock[fileprocessing.FileProcessor](mockController)
	mockFileProcessorFactory := mock.Mock[fileprocessing.FileProcessorFactory](mockController)

	mock.When(mockFileProcessor.ProcessFile(mock.AnyString())).ThenReturn(nil)
	mock.When(mockFileProcessorFactory.CreatePythonFileProcessor()).ThenReturn(mockFileProcessor, nil)

	expectedFilesProcessed := []string{
		"test/main.py",
		"test/sys/dangerous_code.py",
		"test/sys/really_dangerous_code.py",
		"test/subdir/some_code.py",
		"test/subdir/some_more_code.py",
		"test/subdir/subsubdir/dont.py",
		"test/subdir/subsubdir/do.py",
	}

	program.CleanImports(config, filepath, fileManager, mockFileProcessorFactory)

	for _, filepath := range expectedFilesProcessed {
		mock.Verify(mockFileProcessor, mock.Once()).ProcessFile(filepath)
	}
}

func TestCleanImports_FilesIgnored(t *testing.T) {

	config := config.Config{
		Language: "python",
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

	mockController := mock.NewMockController(t)
	mockFileProcessor := mock.Mock[fileprocessing.FileProcessor](mockController)
	mockFileProcessorFactory := mock.Mock[fileprocessing.FileProcessorFactory](mockController)

	mock.When(mockFileProcessor.ProcessFile(mock.AnyString())).ThenReturn(nil)
	mock.When(mockFileProcessorFactory.CreatePythonFileProcessor()).ThenReturn(mockFileProcessor, nil)

	expectedFilesProcessed := []string{
		"test/main.py",
		"test/subdir/some_code.py",
		"test/subdir/some_more_code.py",
		"test/subdir/subsubdir/do.py",
	}

	program.CleanImports(config, filepath, fileManager, mockFileProcessorFactory)

	for _, filepath := range expectedFilesProcessed {
		mock.Verify(mockFileProcessor, mock.Once()).ProcessFile(filepath)
	}
}
