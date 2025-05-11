package config

import (
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

type ConfigValidationError struct {
	InvalidPaths []InvalidPath
}

func (e ConfigValidationError) Error() string {
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

func ProcessAndValidateConfig(rawConfig RawConfig) (Config, *ConfigValidationError) {
	processedGroupingRules := make([]organisation.ImportStatementGroupRule, len(rawConfig.GroupingRules))
	for i, ruleYaml := range rawConfig.GroupingRules {
		processedRule, _ := processAndValidateGroupingRule(ruleYaml)
		processedGroupingRules[i] = processedRule
	}

	return Config{
		IgnoredPaths:  rawConfig.IgnoredPaths,
		GroupingRules: processedGroupingRules,
	}, nil
}
