package headers_test

import (
	"os"
	"testing"

	"github.com/ory/cli/cmd/dev/headers"
	"github.com/ory/cli/cmd/dev/headers/tests"
	"github.com/stretchr/testify/assert"
)

func TestCopy_singleFile_fullDestPath(t *testing.T) {
	rootDir := tests.Dir{Path: "."}
	srcDir := rootDir.CreateDir("test_copy_src")
	dstDir := rootDir.CreateDir("test_copy_dst")
	srcDir.CreateFile("README.md", "# the readme\ntext")
	err := headers.CopyFile("test_copy_src/README.md", dstDir.Filename("README.md"))
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		dstDir.Content("README.md"))
	rootDir.RemoveDir("test_copy_src")
	rootDir.RemoveDir("test_copy_dst")
}

func TestCopy_singleFile_destDir(t *testing.T) {
	rootDir := tests.Dir{Path: "."}
	srcDir := rootDir.CreateDir("test_copy_src")
	dstDir := rootDir.CreateDir("test_copy_dst")
	srcDir.CreateFile("README.md", "# the readme\ntext")
	err := headers.CopyFile("test_copy_src/README.md", dstDir.Path)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		dstDir.Content("README.md"))
	os.RemoveAll(srcDir.Path)
	os.RemoveAll(dstDir.Path)
}

func TestDetermineDestPath_filePath(t *testing.T) {
	t.Parallel()
	root := tests.CreateTmpDir()
	dir := root.CreateDir("dst")
	have := headers.DetermineDestPath("origin/foo.md", dir.Path)
	assert.Equal(t, dir.Filename("foo.md"), have)
}

func TestDetermineDestPath_dirPath(t *testing.T) {
	t.Parallel()
	root := tests.CreateTmpDir()
	dir := root.CreateDir("dst")
	file := dir.CreateFile("foo.md", "")
	have := headers.DetermineDestPath("origin/foo.md", file)
	assert.Equal(t, file, have)
}
