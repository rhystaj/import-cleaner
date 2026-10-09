package config

import (
	"fmt"
	"importcleaner/internal/iowrappers"

	"github.com/go-yaml/yaml"
)

type RawGroupingRuleDefinition struct {
	DependencySourceType string   `yaml:"dependencySourceType"`
	ModulePaths          []string `yaml:"modulePaths"`
	IsTypeImport         string   `yaml:"isTypeImport"`
}

type RawConfig struct {
	Language      string                      `yaml:"language"`
	IgnoredPaths  []string                    `yaml:"ignoredPaths"`
	GroupingRules []RawGroupingRuleDefinition `yaml:"groupingRules"`
}

type ConfigLoadError struct {
	SourceDescription string
}

func (e ConfigLoadError) Error() string {
	return fmt.Sprintf("Could not load config as source '%s' is invalid.", e.SourceDescription)
}

func LoadRawConfigFromYAMLFile(fileManager iowrappers.FileManager, filePath string) (RawConfig, *ConfigLoadError) {
	var config RawConfig

	fileBytes, fileReadError := fileManager.ReadBytesFromFile(filePath)
	if fileReadError != nil {
		return config, &ConfigLoadError{
			SourceDescription: filePath,
		}
	}

	unmarshallingError := yaml.Unmarshal(fileBytes, &config)
	if unmarshallingError != nil {
		return config, &ConfigLoadError{
			SourceDescription: filePath,
		}
	}

	return config, nil
}
