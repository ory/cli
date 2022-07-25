package headers

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/tests"
	"github.com/stretchr/testify/assert"
)

func TestCopyFileToFolderNoSlash(t *testing.T) {
	workspace := setupCopyFile()
	err := CopyFile("test_copy_src/README.md", "test_copy_dst")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		workspace.root.Content("test_copy_dst/README.md"))
	err = cp("test_copy_src/README.md", "test_cp_dst")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	workspace.cleanup()
}

func TestCopyFileToFolderSlash(t *testing.T) {
	workspace := setupCopyFile()
	err := CopyFile("test_copy_src/README.md", "test_copy_dst/")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		workspace.root.Content("test_copy_dst/README.md"))
	err = cp("test_copy_src/README.md", "test_cp_dst/")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	workspace.cleanup()
}

func TestCopyFileToFilepath(t *testing.T) {
	workspace := setupCopyFile()
	err := CopyFile("test_copy_src/README.md", "test_copy_dst/README.md")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		workspace.root.Content("test_copy_dst/README.md"))
	err = cp("test_copy_src/README.md", "test_cp_dst/README.md")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	workspace.cleanup()
}

func TestCopyFilesNoSlash(t *testing.T) {
	workspace := setupCopyFile()
	workspace.src.CreateFile("alpha/one.md", "# Alpha\nOne")
	workspace.src.CreateFile("alpha/two.md", "# Alpha\nTwo")
	workspace.src.CreateFile("beta/one.md", "# Beta\nOne")
	err := CopyFiles("test_copy_src", "test_copy_dst")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/one.md. -->\n\n# Alpha\nOne",
		workspace.dstCopy.Content("alpha/one.md"))
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/two.md. -->\n\n# Alpha\nTwo",
		workspace.dstCopy.Content("alpha/two.md"))
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/beta/one.md. -->\n\n# Beta\nOne",
		workspace.dstCopy.Content("beta/one.md"))
	err = cpr("test_copy_src", "test_cp_dst")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	workspace.cleanup()
}

func TestCopyFilesSlash(t *testing.T) {
	workspace := setupCopyFile()
	workspace.src.CreateFile("alpha/one.md", "# Alpha\nOne")
	workspace.src.CreateFile("alpha/two.md", "# Alpha\nTwo")
	workspace.src.CreateFile("beta/one.md", "# Beta\nOne")
	err := CopyFiles("test_copy_src/", "test_copy_dst/")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/one.md. -->\n\n# Alpha\nOne",
		workspace.dstCopy.Content("alpha/one.md"))
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/two.md. -->\n\n# Alpha\nTwo",
		workspace.dstCopy.Content("alpha/two.md"))
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/beta/one.md. -->\n\n# Beta\nOne",
		workspace.dstCopy.Content("beta/one.md"))
	err = cpr("test_copy_src/", "test_cp_dst/")
	assert.NoError(t, err)
	verifyEqualFolderStructure(t, "test_copy_dst", "test_cp_dst")
	workspace.cleanup()
}

func TestDstPathCpFilePath(t *testing.T) {
	t.Parallel()
	root := tests.CreateTmpDir()
	dst := root.CreateDir("dst")
	give := dst.Path
	want := dst.Filename("foo.md")
	have := destPathCp("origin/foo.md", give)
	assert.Equal(t, want, have)
}

func TestDstPathCpDirPath(t *testing.T) {
	t.Parallel()
	root := tests.CreateTmpDir()
	dir := root.CreateDir("dst")
	file := dir.CreateFile("foo.md", "")
	have := destPathCp("origin/foo.md", file)
	assert.Equal(t, file, have)
}

func TestDstPathCprRoot(t *testing.T) {
	have := dstPathCpr("src/README.md", "src", "dst")
	want := "dst/README.md"
	assert.Equal(t, want, have)
}

func TestDstPathCprSubfolder(t *testing.T) {
	have := dstPathCpr("src/sub1/sub2/README.md", "src", "dst")
	want := "dst/sub1/sub2/README.md"
	assert.Equal(t, want, have)
}

func setupCopyFile() workspace {
	root := tests.Dir{Path: "."}
	src := root.CreateDir("test_copy_src")
	src.CreateFile("README.md", "# the readme\ntext")
	dstCopy := root.CreateDir("test_copy_dst")
	dstCp := root.CreateDir("test_cp_dst")
	cleanup := func() {
		src.Cleanup()
		dstCopy.Cleanup()
		dstCp.Cleanup()
	}
	return workspace{
		root,
		src,
		dstCopy,
		dstCp,
		cleanup,
	}
}

// directory structure for testing copy operations
type workspace struct {
	// the directory that contains the workspace
	root tests.Dir
	// the directory that contains the folder tree to copy
	src tests.Dir
	// the directory that contains the result of the built-in CopyFile(s) operation
	dstCopy tests.Dir
	// the directory that contains the result of Unix's cp operation
	dstCp tests.Dir
	// removes this workspace
	cleanup func()
}

// executes the unix "cp" command
func cp(src, dst string) error {
	_, err := exec.Command("cp", src, dst).CombinedOutput()
	return err
}

// executes the unix "cp -r" command
func cpr(src, dst string) error {
	_, err := exec.Command("cp", "-r", src, dst).CombinedOutput()
	return err
}

// ensures that the two given directories contain files with the same names
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
