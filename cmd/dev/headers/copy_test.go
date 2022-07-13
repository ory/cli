package headers

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/tests"
	"github.com/stretchr/testify/assert"
)

func TestCopyFileToFolderNoSlash(t *testing.T) {
	rootDir, cleanup := setupCopyFile()
	err := CopyFile("test_copy_src/README.md", "test_copy_dst")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		rootDir.Content("test_copy_dst/README.md"))
	err = cp("test_copy_src/README.md", "test_cp_dst")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	cleanup()
}

func TestCopyFileToFolderSlash(t *testing.T) {
	rootDir, cleanup := setupCopyFile()
	err := CopyFile("test_copy_src/README.md", "test_copy_dst/")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		rootDir.Content("test_copy_dst/README.md"))
	err = cp("test_copy_src/README.md", "test_cp_dst/")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	cleanup()
}

func TestCopyFileToFilepath(t *testing.T) {
	rootDir, cleanup := setupCopyFile()
	err := CopyFile("test_copy_src/README.md", "test_copy_dst/README.md")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		rootDir.Content("test_copy_dst/README.md"))
	err = cp("test_copy_src/README.md", "test_cp_dst/README.md")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	cleanup()
}

func setupCopyFile() (tests.Dir, func()) {
	rootDir := tests.Dir{Path: "."}
	srcDir := rootDir.CreateDir("test_copy_src")
	srcDir.CreateFile("README.md", "# the readme\ntext")
	dstCopy := rootDir.CreateDir("test_copy_dst")
	dstCp := rootDir.CreateDir("test_cp_dst")
	return srcDir, func() {
		srcDir.Cleanup()
		dstCopy.Cleanup()
		dstCp.Cleanup()
	}
}

// cp executes the unix "cp" command
func cp(src, dst string) error {
	_, err := exec.Command("cp", src, dst).CombinedOutput()
	return err
}

func verifyEqualFolderStructure(t *testing.T, copyDir string, cpDir string) {
	t.Helper()
	copyEntries := []string{}
	filepath.WalkDir(cpDir, func(path string, entry fs.DirEntry, err error) error {
		copyEntries = append(copyEntries, path)
		return nil
	})
	cpEntries := []string{}
	filepath.WalkDir(cpDir, func(path string, entry fs.DirEntry, err error) error {
		cpEntries = append(cpEntries, path)
		return nil
	})
}

func TestCopyRecursive(t *testing.T) {
	rootDir := tests.Dir{Path: "."}
	srcDir := rootDir.CreateDir("test_copy_src")
	srcDir.CreateFile("alpha/one.md", "# Alpha\nOne")
	srcDir.CreateFile("alpha/two.md", "# Alpha\nTwo")
	srcDir.CreateFile("beta/one.md", "# Beta\nOne")
	dstDir := rootDir.CreateDir("test_copy_dst")
	err := CopyFiles("test_copy_src", dstDir.Path)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# Alpha\nOne",
		dstDir.Content("alpha/one.md"))
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# Alpha\nTwo",
		dstDir.Content("alpha/two.md"))
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# Beta\nOne",
		dstDir.Content("beta/one.md"))
	os.RemoveAll(srcDir.Path)
	os.RemoveAll(dstDir.Path)
}

func TestDetermineDestPath_filePath(t *testing.T) {
	t.Parallel()
	root := tests.CreateTmpDir()
	dir := root.CreateDir("dst")
	have := determineDestPath("origin/foo.md", dir.Path)
	assert.Equal(t, dir.Filename("foo.md"), have)
}

func TestDetermineDestPath_dirPath(t *testing.T) {
	t.Parallel()
	root := tests.CreateTmpDir()
	dir := root.CreateDir("dst")
	file := dir.CreateFile("foo.md", "")
	have := determineDestPath("origin/foo.md", file)
	assert.Equal(t, file, have)
}
