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

	rawConfig, configLoadError := config.LoadRawConfigFromYAMLFile(fileManager, fullConfigFilePath)
	if configLoadError != nil {
		fmt.Print(configLoadError.Error())
	}

	fmt.Printf("Raw config: %+v", rawConfig)

	config, configValidationError := config.ProcessAndValidateConfig(rawConfig)
	if configValidationError != nil {
		fmt.Print(configValidationError.Error())
		os.Exit(1)
	}

	fmt.Printf("Config: %+v", config)

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
