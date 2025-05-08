package iowrappers

import (
	"fmt"
	"iter"
	"os"
)

type DirectoryItemInfo struct {
	IsDir    bool
	ItemName string
}

type FileReadError struct {
	filePath string
}

func (e FileReadError) Error() string {
	return fmt.Sprintf("File '%s' could not be read successfully", e.filePath)
}

type FileManager interface {
	ReadDirectory(directoryName string) iter.Seq[DirectoryItemInfo]
	ReadBytesFromFile(filepath string) ([]byte, *FileReadError)
	ReadStringFromFile(filepath string) (string, *FileReadError)
	WriteContentsToFile(filepath string, contents string) error
}

type FileManagerImpl struct{}

func (dr FileManagerImpl) ReadDirectory(directoryName string) iter.Seq[DirectoryItemInfo] {
	return func(yield func(DirectoryItemInfo) bool) {
		files, workingDirReadError := os.ReadDir(directoryName)
		if workingDirReadError != nil {
			return
		}

		for _, file := range files {
			if !yield(DirectoryItemInfo{
				IsDir:    file.IsDir(),
				ItemName: file.Name(),
			}) {
				return
			}
		}
	}
}

func (dr FileManagerImpl) ReadBytesFromFile(filepath string) ([]byte, *FileReadError) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return make([]byte, 0), &FileReadError{
			filePath: filepath,
		}
	}

	return contents, nil
}

func (dr FileManagerImpl) ReadStringFromFile(filepath string) (string, *FileReadError) {
	contents, err := dr.ReadBytesFromFile(filepath)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func (dr FileManagerImpl) WriteContentsToFile(filepath string, contents string) error {
	file, createError := os.Create(filepath)
	if createError != nil {
		return createError
	}
	defer file.Close()

	file.WriteString(contents)

	return nil
}
