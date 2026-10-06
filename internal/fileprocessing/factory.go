package fileprocessing

import (
	"importcleaner/internal/config"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/languageprocessing/python"
	"importcleaner/internal/languageprocessing/typescript"
	"importcleaner/internal/organisation"
)

type PythonFileProcessor = FileProcessorImpl
type TypescriptFileProcessor = FileProcessorImpl

type FileProcessorFactory interface {
	CreatePythonFileProcessor() (FileProcessor, error)
	CreateTypescriptFileProcessor() (FileProcessor, error)
}

type FileProcessorFactoryImpl struct {
	Config      config.Config
	WorkingDir  string
	FileManager iowrappers.FileManager
}

func (f FileProcessorFactoryImpl) CreatePythonFileProcessor() (FileProcessor, error) {
	groupingRules := make([]organisation.ImportStatementGroupRule, len(f.Config.GroupingRules))
	for i, rule := range f.Config.GroupingRules {
		groupingRules[i] = organisation.NewGroupingRuleFromDefinition(rule)
	}

	return PythonFileProcessor{
		Parser: python.PythonImportStatementParser{
			WorkingDir:  f.WorkingDir,
			FileManager: f.FileManager,
		},
		FileManager: f.FileManager,
		Organiser: organisation.ImportStatementOrganiserImpl{
			GroupRules: groupingRules,
		},
	}, nil
}

func (f FileProcessorFactoryImpl) CreateTypescriptFileProcessor() (FileProcessor, error) {
	groupingRules := make([]organisation.ImportStatementGroupRule, len(f.Config.GroupingRules))
	for i, rule := range f.Config.GroupingRules {
		groupingRules[i] = organisation.NewGroupingRuleFromDefinition(rule)
	}

	return TypescriptFileProcessor{
		Parser:      typescript.TypescriptImportStatementParser{},
		FileManager: f.FileManager,
		Organiser: organisation.ImportStatementOrganiserImpl{
			GroupRules: groupingRules,
		},
	}, nil
}
