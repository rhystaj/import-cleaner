package main

import (
	"fmt"
	"importcleaner/internal/datastructures"
	fileprocessing "importcleaner/internal/file_processing"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/parsing"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	targetRootDir := os.Args[1]

	workingDir, _ := os.Getwd()

	fullTargetRootDirPath := filepath.Join(workingDir, targetRootDir)

	fileManager := iowrappers.FileManagerImpl{}

	parser := parsing.PythonImportStatementParser{
		WorkingDir:  fullTargetRootDirPath,
		FileManager: fileManager,
	}

	fileProcessor := fileprocessing.FileProcessor{
		FileManager: fileManager,
	}

	directoryStack := datastructures.CreateNewStack[string]()
	directoryStack.Push(&fullTargetRootDirPath)

	for directoryStack.Size() > 0 {
		currentDir := *directoryStack.Pop()

		for dirItem := range fileManager.ReadDirectory(currentDir) {
			fullItemPath := filepath.Join(currentDir, dirItem.ItemName)

			if dirItem.IsDir {
				directoryStack.Push(&fullItemPath)
				continue
			}

			if !strings.HasSuffix(fullItemPath, ".py") {
				continue
			}

			relativeItemPath := strings.TrimPrefix(fullItemPath, fullTargetRootDirPath)
			fileProcessingError := fileProcessor.ProcessFile(fullItemPath, parser)
			if fileProcessingError != nil {
				fmt.Printf("Error reading file '%s': %s\n", relativeItemPath, fileProcessingError.Error())
				continue
			}

			fmt.Printf("File '%s' processed successfully\n", relativeItemPath)
		}
	}
}
