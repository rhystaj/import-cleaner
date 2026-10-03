package program

import (
	"fmt"
	"importcleaner/internal/config"
	"importcleaner/internal/datastructures"
	fileprocessing "importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
	"path/filepath"
	"strings"
)

type LanguageConfig struct {
	FileExtension string
	FileProcessor fileprocessing.FileProcessor
}

func isFileForLanguage(filePath string, language string) bool {
	switch language {
	case "python":
		return strings.HasSuffix(filePath, ".py")
	case "typescript":
		return strings.HasSuffix(filePath, ".ts")
	default:
		return false
	}
}

func getFileProcessorForLanguage(language string, fileProcessorFactory fileprocessing.FileProcessorFactory) (fileprocessing.FileProcessor, error) {
	switch language {
	case "python":
		return fileProcessorFactory.CreatePythonFileProcessor()
	case "typescript":
		return fileProcessorFactory.CreateTypescriptFileProcessor()
	default:
		return nil, fmt.Errorf("Language %s not supported", language)
	}
}

func CleanImports(
	config config.Config,
	targetDirPath string,
	fileManager iowrappers.FileManager,
	fileProcessorFactory fileprocessing.FileProcessorFactory) error {

	directoryStack := datastructures.CreateNewStack[string]()
	directoryStack.Push(&targetDirPath)

	fileProcessor, fileProcessorError := getFileProcessorForLanguage(config.Language, fileProcessorFactory)
	if fileProcessorError != nil {
		return fileProcessorError
	}

	for directoryStack.Size() > 0 {
		currentDir := *directoryStack.Pop()

		for dirItem := range fileManager.ReadDirectory(currentDir) {
			fullItemPath := filepath.Join(currentDir, dirItem.ItemName)
			relativeItemPath := strings.TrimPrefix(fullItemPath, targetDirPath+"/")

			ignorePath := false
			for _, ignoredPath := range config.IgnoredPaths {
				if strings.HasPrefix(relativeItemPath, ignoredPath) {
					ignorePath = true
					break
				}
			}
			if ignorePath {
				continue
			}

			if dirItem.IsDir {
				directoryStack.Push(&fullItemPath)
				continue
			}

			if !isFileForLanguage(fullItemPath, config.Language) {
				continue
			}

			fileProcessingError := fileProcessor.ProcessFile(fullItemPath)
			if fileProcessingError != nil {
				// TODO: Have these errors returned in the result and reported by the main program
				fmt.Printf("Error reading file '%s': %s\n", relativeItemPath, fileProcessingError.Error())
				continue
			}
		}
	}

	return nil
}
