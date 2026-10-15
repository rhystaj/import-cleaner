//go:build integration

package integrationtest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertDirectoriesEqual(t *testing.T, leftDir, rightDir string, fileSuffixes []string) {
	t.Helper()

	leftEntries, err := collectTree(leftDir)
	assert.NoError(t, err)
	rightEntries, err := collectTree(rightDir)
	assert.NoError(t, err)

	leftEntries = filterBySuffixes(leftEntries, fileSuffixes)
	rightEntries = filterBySuffixes(rightEntries, fileSuffixes)

	assert.Equal(t, treeSignature(leftEntries), treeSignature(rightEntries), "directory structure differs")

	for _, entry := range leftEntries {
		if entry.isDir {
			continue
		}

		leftContents, err := os.ReadFile(filepath.Join(leftDir, entry.relPath))
		assert.NoError(t, err)
		rightContents, err := os.ReadFile(filepath.Join(rightDir, entry.relPath))
		assert.NoError(t, err)

		assert.Equal(t, string(leftContents), string(rightContents), "file contents differ: %s", entry.relPath)
	}
}
