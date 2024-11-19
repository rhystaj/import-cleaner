package iowrappers

import (
	"iter"
)

type MockFileManager struct {
	Items        []DirectoryItemInfo
	FileContents map[string]string
}

func (dr MockFileManager) ReadDirectory(_ string) iter.Seq[DirectoryItemInfo] {
	return func(yield func(DirectoryItemInfo) bool) {
		for _, item := range dr.Items {
			if !(yield(item)) {
				return
			}
		}
	}
}

func (dr MockFileManager) ReadContentsFromFile(fileName string) (string, *FileReadError) {
	return dr.FileContents[fileName], nil
}

func (dr MockFileManager) WriteContentsToFile(fileName string, contents string) error {
	dr.FileContents[fileName] = contents
	return nil
}
