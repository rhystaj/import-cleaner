package config

import (
	"importcleaner/internal/iowrappers"

	"github.com/go-yaml/yaml"
)

type Config struct {
	IgnoredDirs []string
}

type ConfigLoadError struct{}

func (e ConfigLoadError) Error() string {
	return "Config is invalid and couldn't be loaded."
}

func ReadConfigYAML(fileManager iowrappers.FileManager, filePath string) (Config, *ConfigLoadError) {
	var config Config

	fileBytes, fileReadError := fileManager.ReadBytesFromFile(filePath)
	if fileReadError != nil {
		return config, &ConfigLoadError{}
	}

	unmarshallingError := yaml.Unmarshal(fileBytes, &config)
	if unmarshallingError != nil {
		return config, &ConfigLoadError{}
	}

	return config, nil
}
