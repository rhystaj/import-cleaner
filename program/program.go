package program

import (
	"fmt"
	"importcleaner/internal/config"
	"importcleaner/internal/datastructures"
	fileprocessing "importcleaner/internal/fileprocessing"
	"importcleaner/internal/iowrappers"
	"importcleaner/internal/parsing"
	"path/filepath"
	"strings"
)

func CleanImports(
	config config.Config,
	targetDirPath string,
	fileManager iowrappers.FileManager,
	parser parsing.ImportStatementParser,
	fileProcessor fileprocessing.FileProcessor) {

	fmt.Printf("Config: %+v\n", config)

	directoryStack := datastructures.CreateNewStack[string]()
	directoryStack.Push(&targetDirPath)

	for directoryStack.Size() > 0 {
		currentDir := *directoryStack.Pop()

		for dirItem := range fileManager.ReadDirectory(currentDir) {
			fullItemPath := filepath.Join(currentDir, dirItem.ItemName)

			if dirItem.IsDir {
				directoryStack.Push(&fullItemPath)
				continue
			}

			if !strings.HasSuffix(fullItemPath, ".py") {
				continue
			}

			relativeItemPath := strings.TrimPrefix(fullItemPath, targetDirPath)
			fileProcessingError := fileProcessor.ProcessFile(fullItemPath, parser)
			if fileProcessingError != nil {
				fmt.Printf("Error reading file '%s': %s\n", relativeItemPath, fileProcessingError.Error())
				continue
			}

			fmt.Printf("File '%s' processed successfully\n", relativeItemPath)
		}
	}

}
