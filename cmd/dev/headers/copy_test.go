package headers

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/tests"
	"github.com/stretchr/testify/assert"
)

func Test_CopyFile_fromFile_toNonExistingPath_noSlash(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCp(t, "test_src/README.md", "{{dstDir}}/README.md")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.done(t)
}

func Test_CopyFile_fromFile_toNonExistingPath_withSlash(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifyCpAndCopyErr(t, "test_src/README.md", "{{dstDir}}/new/")
	workspace.done(t)
}

func Test_CopyFile_fromFile_toExistingFile(t *testing.T) {
	workspace := createWorkspace()
	workspace.dstCopy.CreateFile("README.md", "existing content")
	workspace.dstCp.CreateFile("README.md", "existing content")
	workspace.verifySameBehaviorAsCp(t, "test_src/README.md", "{{dstDir}}/README.md")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.done(t)
}

func Test_CopyFile_fromFile_toExistingFolder(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCp(t, "test_src/README.md", "{{dstDir}}")
	workspace.verifyContent(
		t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.done(t)
}

func Test_CopyFile_fromFolder(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifyCpAndCopyErr(t, "test_src", "{{dstDir}}")
	workspace.done(t)
}

func Test_CopyFiles_fromFolder_toFolder(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCpr(t, "test_src", "{{dstDir}}")
	workspace.verifyContent(t,
		"test_copy_dst/test_src/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.verifyContent(t,
		"test_copy_dst/test_src/alpha/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/one.md. -->

# Alpha
One`)
	workspace.verifyContent(t,
		"test_copy_dst/test_src/alpha/two.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/two.md. -->

# Alpha
Two`)
	workspace.verifyContent(t,
		"test_copy_dst/test_src/beta/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/beta/one.md. -->

# Beta
One`)
	workspace.done(t)
}

func Test_CopyFiles_fromFolder_toNonExistingPath(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCpr(t, "test_src", "{{dstDir}}/new")
	workspace.verifyContent(t,
		"test_copy_dst/new/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.verifyContent(t,
		"test_copy_dst/new/alpha/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/one.md. -->

# Alpha
One`)
	workspace.verifyContent(t,
		"test_copy_dst/new/alpha/two.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/two.md. -->

# Alpha
Two`)
	workspace.verifyContent(t,
		"test_copy_dst/new/beta/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/beta/one.md. -->

# Beta
One`)
	workspace.done(t)
}

func Test_CopyFiles_fromFolder_toFile(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifyCprAndCopyFilesErr(t, "test_src", "main.go")
	workspace.done(t)
}

func Test_CopyFiles_fromFile_toFile(t *testing.T) {
	workspace := createWorkspace()
	workspace.dstCopy.CreateFile("README.md", "old content")
	workspace.dstCp.CreateFile("README.md", "old content")
	workspace.verifySameBehaviorAsCpr(t, "test_src/README.md", "{{dstDir}}/README.md")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.done(t)
}

func Test_CopyFiles_fromFile_toNonExistingPath(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCpr(t, "test_src/README.md", "{{dstDir}}/README.md")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.done(t)
}

func Test_CopyFiles_fromFile_toExistingFile(t *testing.T) {
	workspace := createWorkspace()
	workspace.dstCopy.CreateFile("README.md", "existing content")
	workspace.dstCp.CreateFile("README.md", "existing content")
	workspace.verifySameBehaviorAsCpr(t, "test_src/README.md", "{{dstDir}}/README.md")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.done(t)
}

func Test_CopyFiles_fromFile_toExistingFolder(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCpr(t, "test_src/README.md", "{{dstDir}}")
	workspace.verifyContent(
		t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_src/README.md. -->

# readme header
readme text`)
	workspace.done(t)
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
	// list of file paths whose content was verified
	verified []string
}

func createWorkspace() workspace {
	root := tests.Dir{Path: "."}
	src := root.CreateDir("test_src")
	src.CreateFile("README.md", "# readme header\nreadme text")
	src.CreateFile("alpha/one.md", "# Alpha\nOne")
	src.CreateFile("alpha/two.md", "# Alpha\nTwo")
	src.CreateFile("beta/one.md", "# Beta\nOne")
	dstCopy := root.CreateDir("test_copy_dst")
	dstCp := root.CreateDir("test_cp_dst")
	return workspace{root: root, src: src, dstCopy: dstCopy, dstCp: dstCp, verified: []string{}}
}

// removes this test workspace from the filesystem
func (ws *workspace) cleanup() {
	ws.src.Cleanup()
	ws.dstCopy.Cleanup()
	ws.dstCp.Cleanup()
}

func (ws *workspace) done(t *testing.T) {
	ws.verifyAllFilesCompared(t)
	ws.cleanup()
}

// ensures that all files in the workspace have been verified with ws.verifyContent
func (ws *workspace) verifyAllFilesCompared(t *testing.T) {
	allFiles, err := ws.copyFiles("")
	assert.NoError(t, err)
	assert.Equal(t, allFiles, ws.verified)
}

// ensures that the file with the given path in the test workspace
// contains the given content
func (ws *workspace) verifyContent(t *testing.T, filepath, want string) {
	t.Helper()
	have := ws.root.Content(filepath)
	assert.Equal(t, tests.Trim(want), have)
	ws.verified = append(ws.verified, filepath[len(ws.dstCopy.Path):])
}

// ensures that the "CopyFile" function copies files
// the exact same way as the built-in "cp" command in Unix
func (ws *workspace) verifySameBehaviorAsCp(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "cp"
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, err := exec.Command("cp", src, dstCp).CombinedOutput()
	assert.NoError(t, err)
	// run "CopyFile"
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	err = CopyFile(src, dstCopy)
	assert.NoError(t, err)
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFiles" function copies files the exact same way as the built-in "cp -r" command in Unix.
func (ws *workspace) verifySameBehaviorAsCpr(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "cp -r"
	dst := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, err := exec.Command("cp", "-rv", src, dst).CombinedOutput()
	assert.NoError(t, err)
	// run "CopyFile"
	dst = strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	err = CopyFiles(src, dst)
	assert.NoError(t, err)
	// verify that both created the same files and folders
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFile" function and Unix "cp" tool
// both return an error
func (ws *workspace) verifyCpAndCopyErr(t *testing.T, src, dstTemplate string) {
	t.Helper()
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, cpErr := exec.Command("cp", src, dstCp).CombinedOutput()
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	copyErr := CopyFile(src, dstCopy)
	if (copyErr == nil) || (cpErr == nil) {
		t.Fatalf("Unexpected success! cp: %v, copy: %v\n", cpErr, copyErr)
	}
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFile" function and Unix "cp" tool
// both return an error
func (ws *workspace) verifyCprAndCopyFilesErr(t *testing.T, src, dstTemplate string) {
	t.Helper()
	dstCpr := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, cprErr := exec.Command("cp", "-r", src, dstCpr).CombinedOutput()
	dstCopyFiles := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	copyFilesErr := CopyFiles(src, dstCopyFiles)
	if (copyFilesErr == nil) || (cprErr == nil) {
		t.Fatalf("Unexpected success! cp: %v, copy: %v\n", cprErr, copyFilesErr)
	}
	ws.verifyEqualDstStructure(t)
}

// ensures that the two given directories contain files with the same names
func (ws *workspace) verifyEqualDstStructure(t *testing.T) {
	t.Helper()
	copyEntries, err := ws.copyFiles("dst")
	assert.NoError(t, err)
	cpEntries, err := ws.cpFiles("dst")
	assert.NoError(t, err)
	assert.Equal(t, cpEntries, copyEntries)
}

// provides the relative paths of all files that were copied via `Copy` or `CopyFiles`
func (ws *workspace) copyFiles(replacement string) ([]string, error) {
	return ws.files(ws.dstCopy.Path, replacement)
}

// provides the relative paths of all files that were copied via `cp` or `cp -r`
func (ws *workspace) cpFiles(replacement string) ([]string, error) {
	return ws.files(ws.dstCp.Path, replacement)
}

// provides the relative paths of all files in the given folder
func (ws *workspace) files(dir, replacement string) ([]string, error) {
	result := []string{}
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		result = append(result, strings.Replace(path, dir, replacement, 1))
		return nil
	})
	return result, err
}
