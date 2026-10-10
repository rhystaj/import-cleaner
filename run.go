package main

import (
	"importcleaner/internal/config"
	"importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
	"importcleaner/program"
	"path/filepath"
)

const CONFIG_FILE_NAME = "imports-config.yml"

type RunArgs struct {
	WorkingDir    string
	TargetRootDir string
}

func Run(args RunArgs) error {
	fullTargetRootDirPath := filepath.Join(args.WorkingDir, args.TargetRootDir)
	fullConfigFilePath := filepath.Join(fullTargetRootDirPath, CONFIG_FILE_NAME)

	fileManager := iowrappers.FileManagerImpl{}

	rawConfig, configLoadError := config.LoadRawConfigFromYAMLFile(fileManager, fullConfigFilePath)
	if configLoadError != nil {
		return configLoadError
	}

	config, configValidationError := config.ProcessAndValidateConfig(rawConfig)
	if configValidationError != nil {
		return configValidationError
	}

	fileProcessorFactory := fileprocessing.FileProcessorFactoryImpl{
		Config:      config,
		WorkingDir:  fullTargetRootDirPath,
		FileManager: fileManager,
	}

	return program.CleanImports(config, fullTargetRootDirPath, fileManager, fileProcessorFactory)
}
