//go:build integration

package integrationtest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertDirectoriesEqual(t *testing.T, leftDir, rightDir string, fileSuffixes []string) {
	t.Helper()

	leftEntries, err := collectTree(leftDir)
	require.NoError(t, err)
	rightEntries, err := collectTree(rightDir)
	require.NoError(t, err)

	leftEntries = filterBySuffixes(leftEntries, fileSuffixes)
	rightEntries = filterBySuffixes(rightEntries, fileSuffixes)

	require.Equal(t, treeSignature(leftEntries), treeSignature(rightEntries), "directory structure differs")

	for _, entry := range leftEntries {
		if entry.isDir {
			continue
		}

		leftContents, err := os.ReadFile(filepath.Join(leftDir, entry.relPath))
		require.NoError(t, err)
		rightContents, err := os.ReadFile(filepath.Join(rightDir, entry.relPath))
		require.NoError(t, err)

		require.Equal(t, string(leftContents), string(rightContents), "file contents differ: %s", entry.relPath)
	}
}
