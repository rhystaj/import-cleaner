package main

import (
	"fmt"
	"importcleaner/internal/config"
	"importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
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

	config, configReadError := config.ReadConfigYAML(fileManager, fullConfigFilePath)
	if configReadError != nil {
		fmt.Print(configReadError.Error())
		os.Exit(1)
	}

	parser := parsing.PythonImportStatementParser{
		WorkingDir:  fullTargetRootDirPath,
		FileManager: fileManager,
	}

	fileProcessor := fileprocessing.FileProcessorImpl{
		FileManager: fileManager,
	}

	program.CleanImports(config, fullTargetRootDirPath, fileManager, parser, fileProcessor)
}
