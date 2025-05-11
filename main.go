package main

import (
	"fmt"
	"importcleaner/internal/config"
	"importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
	"importcleaner/program"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME = "config.yml"

func main() {
	targetRootDir := os.Args[1]

	workingDir, _ := os.Getwd()

	fullTargetRootDirPath := filepath.Join(workingDir, targetRootDir)
	fullConfigFilePath := filepath.Join(fullTargetRootDirPath, CONFIG_FILE_NAME)

	fileManager := iowrappers.FileManagerImpl{}

	config, configReadError := config.LoadAndValidateConfigFromFile(fileManager, fullConfigFilePath)
	if configReadError != nil {
		fmt.Print(configReadError.Error())
		os.Exit(1)
	}

	parser := parsing.PythonImportStatementParser{
		WorkingDir:  fullTargetRootDirPath,
		FileManager: fileManager,
	}

	organiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: config.GroupingRules,
	}

	fileProcessor := fileprocessing.FileProcessorImpl{
		FileManager: fileManager,
		Organiser:   organiser,
	}

	program.CleanImports(config, fullTargetRootDirPath, fileManager, parser, fileProcessor)
}
