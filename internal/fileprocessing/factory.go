package fileprocessing

import (
	"importcleaner/internal/config"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
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
	return PythonFileProcessor{
		Parser: parsing.PythonImportStatementParser{
			WorkingDir:  f.WorkingDir,
			FileManager: f.FileManager,
		},
		FileManager: f.FileManager,
		Organiser: organisation.ImportStatementOrganiserImpl{
			GroupRules: f.Config.GroupingRules,
		},
	}, nil
}

func (f FileProcessorFactoryImpl) CreateTypescriptFileProcessor() (FileProcessor, error) {
	return TypescriptFileProcessor{
		Parser:      parsing.TypescriptImportStatementParser{},
		FileManager: f.FileManager,
		Organiser: organisation.ImportStatementOrganiserImpl{
			GroupRules: f.Config.GroupingRules,
		},
	}, nil
}
