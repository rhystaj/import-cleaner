package iowrappers

import (
	"iter"
	"path/filepath"
)

type MockFSNode struct {
	info     DirectoryItemInfo
	children []*MockFSNode
	contents *string
}

type MockFileManager struct {
	nodeIndex map[string]*MockFSNode
}

func buildNodeIndexRecursively(currentPath *string, nodes []*MockFSNode, nodeIndex *map[string]*MockFSNode) {
	for _, node := range nodes {
		path := node.info.ItemName
		if currentPath != nil {
			path = filepath.Join(*currentPath, path)
		}

		(*nodeIndex)[path] = node
		buildNodeIndexRecursively(&path, node.children, nodeIndex)
	}
}

func CreateMockFSFile(name string, contents string) *MockFSNode {
	return &MockFSNode{
		info: DirectoryItemInfo{
			ItemName: name,
			IsDir:    false,
		},
		children: make([]*MockFSNode, 0),
		contents: &contents,
	}
}

func CreateMockFSDir(name string, children []*MockFSNode) *MockFSNode {
	return &MockFSNode{
		info: DirectoryItemInfo{
			ItemName: name,
			IsDir:    true,
		},
		children: children,
		contents: nil,
	}
}

func InitialiseMockFileManager(rootNodes []*MockFSNode) MockFileManager {
	nodeIndex := make(map[string]*MockFSNode)
	nodeIndex[""] = CreateMockFSDir("", rootNodes)

	buildNodeIndexRecursively(nil, rootNodes, &nodeIndex)
	return MockFileManager{nodeIndex: nodeIndex}
}

func (dr MockFileManager) ReadDirectory(directoryName string) iter.Seq[DirectoryItemInfo] {
	return func(yield func(DirectoryItemInfo) bool) {
		for _, item := range dr.nodeIndex[directoryName].children {
			if !(yield(item.info)) {
				return
			}
		}
	}
}

func (dr MockFileManager) ReadBytesFromFile(fileName string) ([]byte, *FileReadError) {
	node, fileExists := dr.nodeIndex[fileName]
	if !fileExists {
		return make([]byte, 0), &FileReadError{FilePath: fileName}
	}

	return []byte(*node.contents), nil
}

func (dr MockFileManager) ReadStringFromFile(fileName string) (string, *FileReadError) {
	node, fileExists := dr.nodeIndex[fileName]
	if !fileExists {
		return "", &FileReadError{FilePath: fileName}
	}

	return *node.contents, nil
}

func (dr MockFileManager) WriteContentsToFile(fileName string, contents string) error {
	node, fileExists := dr.nodeIndex[fileName]
	if !fileExists {
		return &FileReadError{FilePath: fileName}
	}

	node.contents = &contents
	return nil
}
