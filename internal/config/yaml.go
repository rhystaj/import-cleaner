package config

import (
	"importcleaner/internal/iowrappers"

	"github.com/go-yaml/yaml"
)

type GroupingRuleYaml struct {
	DependencySourceType string `yaml:"dependencySourceType"`
}

type ConfigYaml struct {
	IgnoredPaths  []string           `yaml:"ignoredPaths"`
	GroupingRules []GroupingRuleYaml `yaml:"groupingRule"`
}

func readConfigYAML(fileManager iowrappers.FileManager, filePath string) (ConfigYaml, *ConfigLoadError) {
	var config ConfigYaml

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
