package config

import (
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/organisation"
	"strings"
)

type Config struct {
	IgnoredPaths  []string
	GroupingRules []organisation.ImportStatementGroupRule
}

type InvalidPath struct {
	Path             string
	ErrorDescription string
}

type ConfigLoadError struct {
	InvalidPaths []InvalidPath
}

func (e ConfigLoadError) Error() string {
	if len(e.InvalidPaths) == 0 {
		return "Config load error returned, but no invalid paths were specified."
	}

	var pathListStringBuilder strings.Builder
	for _, invalidPath := range e.InvalidPaths {
		pathListStringBuilder.WriteString("\n" + invalidPath.Path + " - " + invalidPath.ErrorDescription)
	}

	return "Errors detected in config: " + pathListStringBuilder.String()
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
