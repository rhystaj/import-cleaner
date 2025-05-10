package config

import (
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/organisation"
)

type Config struct {
	IgnoredPaths  []string
	GroupingRules []organisation.ImportStatementGroupRule
}

type ConfigLoadError struct{}

func (e ConfigLoadError) Error() string {
	return "Config is invalid and couldn't be loaded."
}

func LoadAndValidateConfigFromFile(fileManager iowrappers.FileManager, filePath string) (Config, *ConfigLoadError) {
	configYaml, yamlReadError := readConfigYAML(fileManager, filePath)

	panic("Not implemented")
}
