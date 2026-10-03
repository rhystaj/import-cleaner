package main

import (
	"fmt"
	"importcleaner/internal/config"
	"importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
	"importcleaner/program"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME = "imports-config.yml"

func main() {
	targetRootDir := os.Args[1]

	workingDir, _ := os.Getwd()

	fullTargetRootDirPath := filepath.Join(workingDir, targetRootDir)
	fullConfigFilePath := filepath.Join(fullTargetRootDirPath, CONFIG_FILE_NAME)

	fileManager := iowrappers.FileManagerImpl{}

	rawConfig, configLoadError := config.LoadRawConfigFromYAMLFile(fileManager, fullConfigFilePath)
	if configLoadError != nil {
		fmt.Println(configLoadError.Error())
		os.Exit(1)
	}

	config, configValidationError := config.ProcessAndValidateConfig(rawConfig)
	if configValidationError != nil {
		fmt.Println(configValidationError.Error())
		os.Exit(1)
	}

	fileProcessorFactory := fileprocessing.FileProcessorFactoryImpl{
		Config:      config,
		WorkingDir:  fullTargetRootDirPath,
		FileManager: fileManager,
	}

	cleanImportsError := program.CleanImports(config, fullTargetRootDirPath, fileManager, fileProcessorFactory)
	if cleanImportsError != nil {
		fmt.Println(cleanImportsError.Error())
		os.Exit(1)
	}
}
