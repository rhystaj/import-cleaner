package config

import (
	"fmt"
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
		pathListStringBuilder.WriteString("\n\t" + invalidPath.Path + " - " + invalidPath.ErrorDescription)
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

func processAndValidateGroupingRule(rawGroupingRule RawGroupingRule) (organisation.ImportStatementGroupRule, []InvalidPath) {
	invalidPaths := make([]InvalidPath, 0)

	dependencySourceTypeValid, dependencySourceType := parseDependencySourceType(rawGroupingRule.DependencySourceType)
	if !dependencySourceTypeValid {
		invalidPaths = append(invalidPaths, InvalidPath{
			Path:             "groupingRules.dependencySourceType",
			ErrorDescription: fmt.Sprintf("Must be 'internal', or 'external', but was '%s'", rawGroupingRule.DependencySourceType),
		})
	}

	return organisation.ImportStatementGroupRule{
		DependencySourceType: dependencySourceType,
	}, invalidPaths
}

func ProcessAndValidateConfig(rawConfig RawConfig) (Config, *ConfigValidationError) {
	var invalidPaths []InvalidPath

	processedGroupingRules := make([]organisation.ImportStatementGroupRule, len(rawConfig.GroupingRules))
	for i, ruleYaml := range rawConfig.GroupingRules {
		var processedRule organisation.ImportStatementGroupRule
		processedRule, invalidPaths = processAndValidateGroupingRule(ruleYaml)
		processedGroupingRules[i] = processedRule
	}

	var validationError *ConfigValidationError
	if len(invalidPaths) > 0 {
		validationError = &ConfigValidationError{
			InvalidPaths: invalidPaths,
		}
	}

	return Config{
		IgnoredPaths:  rawConfig.IgnoredPaths,
		GroupingRules: processedGroupingRules,
	}, validationError
}
