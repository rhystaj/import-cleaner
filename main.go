package main

import (
	"fmt"
	"importcleaner/internal/config"
	"importcleaner/internal/datastructures"
	fileprocessing "importcleaner/internal/file_processing"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/parsing"
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_FILE_NAME = "config.yml"

func main() {
	targetRootDir := os.Args[1]

	workingDir, _ := os.Getwd()

	fullTargetRootDirPath := filepath.Join(workingDir, targetRootDir)
	fullConfigFilePath := filepath.Join(fullTargetRootDirPath, CONFIG_FILE_NAME)

	fileManager := iowrappers.FileManagerImpl{}

	config, configReadError := config.ReadConfigYAML(fileManager, fullConfigFilePath)
	if configReadError != nil {
		fmt.Print(configReadError.Error())
		os.Exit(1)
	}

	fmt.Printf("Config: %+v", config)

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
