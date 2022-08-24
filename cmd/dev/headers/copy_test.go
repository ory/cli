package headers

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	t.Run("file --> non-existing path ending without slash", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCp(t, "test_src/SECURITY.md", "{{dstDir}}/SECURITY.md")
		workspace.verifyContent(t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("file --> non-existing path ending with slash", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifyCpAndCopyErr(t, "test_src/SECURITY.md", "{{dstDir}}/new/")
		workspace.done(t)
	})

	t.Run("file --> existing file", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.dstCopy.CreateFile("SECURITY.md", "existing content")
		workspace.dstCp.CreateFile("SECURITY.md", "existing content")
		workspace.verifySameBehaviorAsCp(t, "test_src/SECURITY.md", "{{dstDir}}/SECURITY.md")
		workspace.verifyContent(t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("file --> existing folder", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCp(t, "test_src/SECURITY.md", "{{dstDir}}")
		workspace.verifyContent(
			t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("folder", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifyCpAndCopyErr(t, "test_src", "{{dstDir}}")
		workspace.done(t)
	})
}

func TestCopyFileNoOverride(t *testing.T) {
	t.Run("file --> non-existing path ending without slash", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCpn(t, "test_src/SECURITY.md", "{{dstDir}}/SECURITY.md")
		workspace.verifyContent(t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("file --> non-existing path ending with slash", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifyCpnAndCopyErr(t, "test_src/SECURITY.md", "{{dstDir}}/new/")
		workspace.done(t)
	})

	t.Run("file --> existing file", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.dstCopy.CreateFile("SECURITY.md", "existing content")
		workspace.dstCp.CreateFile("SECURITY.md", "existing content")
		workspace.verifySameBehaviorAsCpn(t, "test_src/SECURITY.md", "{{dstDir}}/SECURITY.md")
		workspace.verifyContent(t, "test_copy_dst/SECURITY.md", `existing content`)
		workspace.done(t)
	})

	t.Run("file --> existing folder", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCpn(t, "test_src/SECURITY.md", "{{dstDir}}")
		workspace.verifyContent(
			t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("folder", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifyCpnAndCopyErr(t, "test_src", "{{dstDir}}")
		workspace.done(t)
	})
}

func TestCopyFiles(t *testing.T) {
	t.Run("folder --> folder", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCpr(t, "test_src", "{{dstDir}}")
		workspace.verifyContent(t,
			"test_copy_dst/test_src/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.verifyContent(t,
			"test_copy_dst/test_src/alpha/one.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/one.yml

name: alpha
number: one`)
		workspace.verifyContent(t,
			"test_copy_dst/test_src/alpha/two.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/two.yml

name: alpha
number: two`)
		workspace.verifyContent(t,
			"test_copy_dst/test_src/beta/one.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/beta/one.yml

name: beta
number: one`)
		workspace.done(t)
	})

	t.Run("subfolder --> folder", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCpr(t, "test_src/alpha", "{{dstDir}}")
		workspace.verifyContent(t,
			"test_copy_dst/alpha/one.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/one.yml

name: alpha
number: one`)
		workspace.verifyContent(t,
			"test_copy_dst/alpha/two.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/two.yml

name: alpha
number: two`)
		workspace.done(t)
	})

	t.Run("folder --> non-existing path", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCpr(t, "test_src", "{{dstDir}}/new")
		workspace.verifyContent(t,
			"test_copy_dst/new/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.verifyContent(t,
			"test_copy_dst/new/alpha/one.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/one.yml

name: alpha
number: one`)
		workspace.verifyContent(t,
			"test_copy_dst/new/alpha/two.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/alpha/two.yml

name: alpha
number: two`)
		workspace.verifyContent(t,
			"test_copy_dst/new/beta/one.yml", `
# AUTO-GENERATED, DO NOT EDIT!
# Please edit the original at https://github.com/ory/meta/blob/master/test_src/beta/one.yml

name: beta
number: one`)
		workspace.done(t)
	})

	t.Run("folder --> file", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifyCprAndCopyFilesErr(t, "test_src", "main.go")
		workspace.done(t)
	})

	t.Run("file --> file", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.dstCopy.CreateFile("SECURITY.md", "old content")
		workspace.dstCp.CreateFile("SECURITY.md", "old content")
		workspace.verifySameBehaviorAsCpr(t, "test_src/SECURITY.md", "{{dstDir}}/SECURITY.md")
		workspace.verifyContent(t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("file --> non-existing path", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCpr(t, "test_src/SECURITY.md", "{{dstDir}}/SECURITY.md")
		workspace.verifyContent(t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("file --> existing file", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.dstCopy.CreateFile("SECURITY.md", "existing content")
		workspace.dstCp.CreateFile("SECURITY.md", "existing content")
		workspace.verifySameBehaviorAsCpr(t, "test_src/SECURITY.md", "{{dstDir}}/SECURITY.md")
		workspace.verifyContent(t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})

	t.Run("file --> existing folder", func(t *testing.T) {
		workspace := createWorkspace()
		workspace.verifySameBehaviorAsCpr(t, "test_src/SECURITY.md", "{{dstDir}}")
		workspace.verifyContent(
			t,
			"test_copy_dst/SECURITY.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! -->
<!-- Please edit the original at https://github.com/ory/meta/blob/master/test_src/SECURITY.md -->

# header
text about security`)
		workspace.done(t)
	})
}

// directory structure for testing copy operations
type workspace struct {
	// the directory that contains the workspace
	root Dir
	// the directory that contains the folder tree to copy
	src Dir
	// the directory that contains the result of the built-in CopyFile(s) operation
	dstCopy Dir
	// the directory that contains the result of Unix's cp operation
	dstCp Dir
	// list of file paths whose content was verified
	verified []string
}

func createWorkspace() workspace {
	root := Dir{Path: "."}
	src := root.CreateDir("test_src")
	src.CreateFile("SECURITY.md", "# header\ntext about security")
	src.CreateFile("alpha/one.yml", "name: alpha\nnumber: one")
	src.CreateFile("alpha/two.yml", "name: alpha\nnumber: two")
	src.CreateFile("beta/one.yml", "name: beta\nnumber: one")
	dstCopy := root.CreateDir("test_copy_dst")
	dstCp := root.CreateDir("test_cp_dst")
	return workspace{root: root, src: src, dstCopy: dstCopy, dstCp: dstCp, verified: []string{}}
}

// removes this test workspace from the filesystem
func (ws *workspace) delete() {
	ws.src.Delete()
	ws.dstCopy.Delete()
	ws.dstCp.Delete()
}

// cleanup of this workspace at the end of a test
func (ws *workspace) done(t *testing.T) {
	ws.verifyAllFilesChecked(t)
	ws.delete()
}

// ensures that all files in the workspace have been verified with ws.verifyContent
func (ws *workspace) verifyAllFilesChecked(t *testing.T) {
	allFiles, err := ws.copiedFiles("")
	assert.NoError(t, err)
	assert.Equal(t, allFiles, ws.verified)
}

// ensures that the file with the given path in the test workspace
// contains the given content
func (ws *workspace) verifyContent(t *testing.T, filepath, want string) {
	t.Helper()
	have := ws.root.Content(filepath)
	assert.Equal(t, strings.Trim(want, "\n"), have)
	ws.verified = append(ws.verified, filepath[len(ws.dstCopy.Path):])
}

// ensures that the "CopyFile" function copies files
// the exact same way as the built-in "cp" command in Unix
func (ws *workspace) verifySameBehaviorAsCp(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "CopyFile"
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	err := CopyFile(src, dstCopy)
	assert.NoError(t, err)
	// run "cp"
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, err = exec.Command("cp", src, dstCp).CombinedOutput()
	assert.NoError(t, err)
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFile" function copies files
// the exact same way as the built-in "cp" command in Unix
func (ws *workspace) verifySameBehaviorAsCpn(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "CopyFileNoOverwrite"
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	err := CopyFileNoOverwrite(src, dstCopy)
	assert.NoError(t, err)
	// This function verifies that `ory dev headers cp` behaves the same as the built-in `cp` command on Linux.
	// The `cp` command on other OS like macOS has different behavior. This feature is used only on CI.
	if runtime.GOOS != "linux" {
		return
	}
	// run "cp -n"
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, err = exec.Command("cp", "-n", src, dstCp).CombinedOutput()
	assert.NoError(t, err)
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFiles" function copies files the exact same way as the built-in "cp -r" command in Unix.
func (ws *workspace) verifySameBehaviorAsCpr(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "CopyFiles"
	dst := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	err := CopyFiles(src, dst)
	assert.NoError(t, err)
	// This function verifies that `ory dev headers cp` behaves the same as the built-in `cp` command on Linux.
	// The `cp` command on other OS like macOS has different behavior. This feature is used only on CI.
	if runtime.GOOS != "linux" {
		return
	}
	// run "cp -r"
	dst = strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, err = exec.Command("cp", "-rv", src, dst).CombinedOutput()
	assert.NoError(t, err)
	// verify that both created the same files and folders
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFile" function and Unix "cp" tool
// both return an error
func (ws *workspace) verifyCpAndCopyErr(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "CopyFile"
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	copyErr := CopyFile(src, dstCopy)
	// This function verifies that `ory dev headers cp` behaves the same as the built-in `cp` command on Linux.
	// The `cp` command on other OS like macOS has different behavior. This feature is used only on CI.
	if runtime.GOOS != "linux" {
		return
	}
	// run "cp"
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, cpErr := exec.Command("cp", src, dstCp).CombinedOutput()
	// ensure both return errors
	if (copyErr == nil) || (cpErr == nil) {
		t.Fatalf("Unexpected success! cp: %v, copy: %v\n", cpErr, copyErr)
	}
	// verify that both created the same files and folders
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFileNoOverwrite" function and Unix "cp -n" tool
// both return an error
func (ws *workspace) verifyCpnAndCopyErr(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "CopyFileNoOverwrite"
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	copyErr := CopyFileNoOverwrite(src, dstCopy)
	// This function verifies that `ory dev headers cp` behaves the same as the built-in `cp` command on Linux.
	// The `cp` command on other OS like macOS has different behavior. This feature is used only on CI.
	if runtime.GOOS != "linux" {
		return
	}
	// run "cp -n"
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, cpErr := exec.Command("cp", "-n", src, dstCp).CombinedOutput()
	// ensure both return errors
	if (copyErr == nil) || (cpErr == nil) {
		t.Fatalf("Unexpected success! cp: %v, copy: %v\n", cpErr, copyErr)
	}
	// verify that both created the same files and folders
	ws.verifyEqualDstStructure(t)
}

// ensures that the "CopyFile" function and Unix "cp" tool
// both return an error
func (ws *workspace) verifyCprAndCopyFilesErr(t *testing.T, src, dstTemplate string) {
	t.Helper()
	// run "CopyFiles"
	dstCopyFiles := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	copyFilesErr := CopyFiles(src, dstCopyFiles)
	// This function verifies that `ory dev headers cp` behaves the same as the built-in `cp` command on Linux.
	// The `cp` command on other OS like macOS has different behavior. This feature is used only on CI.
	if runtime.GOOS != "linux" {
		return
	}
	// run "cp -r"
	dstCpr := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	_, cprErr := exec.Command("cp", "-r", src, dstCpr).CombinedOutput()
	// ensure both return errors
	if (copyFilesErr == nil) || (cprErr == nil) {
		t.Fatalf("Unexpected success! cp: %v, copy: %v\n", cprErr, copyFilesErr)
	}
	// verify that both created the same files and folders
	ws.verifyEqualDstStructure(t)
}

// ensures that the two given directories contain files with the same names
func (ws *workspace) verifyEqualDstStructure(t *testing.T) {
	t.Helper()
	copyEntries, err := ws.copiedFiles("dst")
	assert.NoError(t, err)
	cpEntries, err := ws.cpedFiles("dst")
	assert.NoError(t, err)
	assert.Equal(t, cpEntries, copyEntries)
}

// provides the relative paths of all files that were copied via `Copy` or `CopyFiles`
func (ws *workspace) copiedFiles(replacement string) ([]string, error) {
	return ws.files(ws.dstCopy.Path, replacement)
}

// provides the relative paths of all files that were copied via `cp` or `cp -r`
func (ws *workspace) cpedFiles(replacement string) ([]string, error) {
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

// a filesystem directory used for testing
type Dir struct {
	Path string
}

func CreateTmpDir() Dir {
	path, err := os.MkdirTemp("", "ory-license")
	if err != nil {
		panic(err)
	}
	return Dir{path}
}

func (t Dir) Content(path string) string {
	content, err := os.ReadFile(filepath.Join(t.Path, path))
	if err != nil {
		panic(err)
	}
	return string(content)
}

func (t Dir) CreateDir(name string) Dir {
	t.RemoveDir(name)
	path := filepath.Join(t.Path, name)
	err := os.Mkdir(path, 0744)
	if err != nil {
		panic(err)
	}
	return Dir{path}
}

func (t Dir) CreateFile(name, content string) string {
	path := filepath.Join(t.Path, name)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, []byte(content), 0744)
	if err != nil {
		panic(err)
	}
	return path
}

func (t Dir) Filename(base string) string {
	return filepath.Join(t.Path, base)
}

func (t Dir) RemoveDir(name string) {
	os.RemoveAll(filepath.Join(t.Path, name))
}

func (t Dir) Delete() {
	err := os.RemoveAll(t.Path)
	if err != nil {
		panic(err)
	}
}
