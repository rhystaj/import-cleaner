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

func parseDependencySourceType(str string) (valid bool, result organisation.DependencySourceType) {
	switch str {
	case "":
		return true, organisation.DependencySourceTypeNone
	case "internal":
		return true, organisation.DependencySourceTypeInternal
	case "external":
		return true, organisation.DependencySourceTypeExternal
	default:
		return false, organisation.DependencySourceTypeNone
	}
}

func processAndValidateGroupingRule(groupingRule GroupingRuleYaml) (organisation.ImportStatementGroupRule, error) {
	_, dependencySourceType := parseDependencySourceType(groupingRule.DependencySourceType)

	return organisation.ImportStatementGroupRule{
		DependencySourceType: dependencySourceType,
	}, nil
}

func LoadAndValidateConfigFromFile(fileManager iowrappers.FileManager, filePath string) (Config, *ConfigLoadError) {
	configYaml, _ := readConfigYAML(fileManager, filePath)

	processedGroupingRules := make([]organisation.ImportStatementGroupRule, len(configYaml.GroupingRules))
	for i, ruleYaml := range configYaml.GroupingRules {
		processedRule, _ := processAndValidateGroupingRule(ruleYaml)
		processedGroupingRules[i] = processedRule
	}

	return Config{
		IgnoredPaths:  configYaml.IgnoredPaths,
		GroupingRules: processedGroupingRules,
	}, nil
}
