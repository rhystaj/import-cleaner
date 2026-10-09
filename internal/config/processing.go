package config

import (
	"fmt"
	types "importcleaner/internal"
	"strings"
)

type DependencySourceType int

const (
	DependencySourceTypeNone DependencySourceType = iota
	DependencySourceTypeInternal
	DependencySourceTypeExternal
)

type ImportStatementGroupRuleDefinition struct {
	DependencySourceType DependencySourceType
	ModulePaths          [][]string
	IsTypeImport         types.OptionalBool
}

type Config struct {
	Language      string
	IgnoredPaths  []string
	GroupingRules []ImportStatementGroupRuleDefinition
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
		pathListStringBuilder.WriteString("\n\t")
		pathListStringBuilder.WriteString(invalidPath.Path)
		pathListStringBuilder.WriteString(" - ")
		pathListStringBuilder.WriteString(invalidPath.ErrorDescription)
	}

	return "Errors detected in config: " + pathListStringBuilder.String()
}

func parseDependencySourceType(str string) (valid bool, result DependencySourceType) {
	switch str {
	case "":
		return true, DependencySourceTypeNone
	case "internal":
		return true, DependencySourceTypeInternal
	case "external":
		return true, DependencySourceTypeExternal
	default:
		return false, DependencySourceTypeNone
	}
}

func parseOptionalBool(str string) (valid bool, result types.OptionalBool) {
	switch str {
	case "":
		return true, types.OBNil
	case "false":
		return true, types.OBFalse
	case "true":
		return true, types.OBTrue
	default:
		return false, types.OBNil
	}
}

func processAndValidateGroupingRule(rawGroupingRule RawGroupingRuleDefinition) (ImportStatementGroupRuleDefinition, []InvalidPath) {
	invalidPaths := make([]InvalidPath, 0)

	dependencySourceTypeValid, dependencySourceType := parseDependencySourceType(rawGroupingRule.DependencySourceType)
	if !dependencySourceTypeValid {
		invalidPaths = append(invalidPaths, InvalidPath{
			Path:             "groupingRules.dependencySourceType",
			ErrorDescription: fmt.Sprintf("Must be 'internal', or 'external', but was '%s'", rawGroupingRule.DependencySourceType),
		})
	}

	processedModulePaths := make([][]string, len(rawGroupingRule.ModulePaths))
	for i, rawModulePath := range rawGroupingRule.ModulePaths {
		processedModulePaths[i] = strings.Split(rawModulePath, ".")
	}

	isTypeImportValid, isTypeImport := parseOptionalBool(rawGroupingRule.IsTypeImport)
	if !isTypeImportValid {
		invalidPaths = append(invalidPaths, InvalidPath{
			Path:             "groupingRules.isTypeImport",
			ErrorDescription: fmt.Sprintf("Must be 'true', or 'false', but was '%s'", rawGroupingRule.IsTypeImport),
		})
	}

	return ImportStatementGroupRuleDefinition{
		DependencySourceType: dependencySourceType,
		ModulePaths:          processedModulePaths,
		IsTypeImport:         isTypeImport,
	}, invalidPaths
}

func ProcessAndValidateConfig(rawConfig RawConfig) (Config, *ConfigValidationError) {
	var configInvalidPaths []InvalidPath

	processedGroupingRules := make([]ImportStatementGroupRuleDefinition, len(rawConfig.GroupingRules))
	for i, ruleYaml := range rawConfig.GroupingRules {
		processedRule, groupRuleInvalidPaths := processAndValidateGroupingRule(ruleYaml)
		processedGroupingRules[i] = processedRule
		configInvalidPaths = append(configInvalidPaths, groupRuleInvalidPaths...)
	}

	var validationError *ConfigValidationError
	if len(configInvalidPaths) > 0 {
		validationError = &ConfigValidationError{
			InvalidPaths: configInvalidPaths,
		}
	}

	return Config{
		Language:      rawConfig.Language,
		IgnoredPaths:  rawConfig.IgnoredPaths,
		GroupingRules: processedGroupingRules,
	}, validationError
}
