package fileprocessing

import (
	"importcleaner/internal/config"
	"importcleaner/internal/iowrappers"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/languageprocessing/python"
	"importcleaner/internal/languageprocessing/typescript"
	"importcleaner/internal/organisation"
)

type PythonFileProcessor = FileProcessorImpl[lp.ImportStatement]
type TypescriptFileProcessor = FileProcessorImpl[typescript.TypescriptImportStatement]

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
	groupingRules := make([]organisation.ImportStatementGroupRule[lp.ImportStatement], len(f.Config.GroupingRules))
	for i, rule := range f.Config.GroupingRules {
		groupingRules[i] = organisation.NewGenericGroupingRuleFromDefinition(rule)
	}

	return PythonFileProcessor{
		Parser: python.PythonImportStatementParser{
			WorkingDir:  f.WorkingDir,
			FileManager: f.FileManager,
		},
		FileManager: f.FileManager,
		Organiser: organisation.ImportStatementOrganiserImpl[lp.ImportStatement]{
			GroupRules: groupingRules,
		},
	}, nil
}

func (f FileProcessorFactoryImpl) CreateTypescriptFileProcessor() (FileProcessor, error) {
	groupingRules := make([]organisation.ImportStatementGroupRule[typescript.TypescriptImportStatement], len(f.Config.GroupingRules))
	for i, rule := range f.Config.GroupingRules {
		groupingRules[i] = typescript.NewTypescriptImportStatementFromDefinition(rule)
	}

	return TypescriptFileProcessor{
		Parser:      typescript.TypescriptImportStatementParser{},
		FileManager: f.FileManager,
		Organiser: organisation.ImportStatementOrganiserImpl[typescript.TypescriptImportStatement]{
			GroupRules: groupingRules,
		},
	}, nil
}
