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

func Test_CopyFile_ToFolder_NoSlash(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCp(t, "test_copy_src/README.md", "{{dstDir}}")
	workspace.verifyContent(
		t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->

# the readme
text`)
	workspace.cleanup()
}

func Test_CopyFile_ToFolder_Slash(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCp(t, "test_copy_src/README.md", "{{dstDir}}/")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->

# the readme
text`)
	workspace.cleanup()
}

func Test_CopyFile_ToFilepath(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCp(t, "test_copy_src/README.md", "{{dstDir}}/README.md")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->

# the readme
text`)
	workspace.cleanup()
}

func Test_CopyFiles_NoSlash(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCpr(t, "test_copy_src", "test_copy_dst")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->

# the readme
text`)
	workspace.verifyContent(t,
		"test_copy_dst/alpha/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/one.md. -->

# Alpha
One`)
	workspace.verifyContent(t,
		"test_copy_dst/alpha/two.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/two.md. -->

# Alpha
Two`)
	workspace.verifyContent(t,
		"test_copy_dst/beta/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/beta/one.md. -->

# Beta
One`)
	workspace.cleanup()
}

func Test_CopyFiles_Slash(t *testing.T) {
	workspace := createWorkspace()
	workspace.verifySameBehaviorAsCpr(t, "test_copy_src", "test_copy_dst/")
	workspace.verifyContent(t,
		"test_copy_dst/README.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->

# the readme
text`)
	workspace.verifyContent(t,
		"test_copy_dst/alpha/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/one.md. -->

# Alpha
One`)
	workspace.verifyContent(t,
		"test_copy_dst/alpha/two.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/alpha/two.md. -->

# Alpha
Two`)
	workspace.verifyContent(t,
		"test_copy_dst/beta/one.md", `
<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/beta/one.md. -->

# Beta
One`)
	workspace.cleanup()
}

func createWorkspace() workspace {
	root := tests.Dir{Path: "."}
	src := root.CreateDir("test_copy_src")
	src.CreateFile("README.md", "# the readme\ntext")
	src.CreateFile("alpha/one.md", "# Alpha\nOne")
	src.CreateFile("alpha/two.md", "# Alpha\nTwo")
	src.CreateFile("beta/one.md", "# Beta\nOne")
	dstCopy := root.CreateDir("test_copy_dst")
	dstCp := root.CreateDir("test_cp_dst")
	return workspace{root, src, dstCopy, dstCp}
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
}

func (ws workspace) cleanup() {
	ws.src.Cleanup()
	ws.dstCopy.Cleanup()
	ws.dstCp.Cleanup()
}

func (ws workspace) verifyContent(t *testing.T, filepath, want string) {
	have := ws.root.Content(filepath)
	assert.Equal(t, tests.Trim(want), have)
}

// verifies that the "CopyFile" function copies files the exact same way as the built-in "cp" command in Unix.
func (ws workspace) verifySameBehaviorAsCp(t *testing.T, src, dstTemplate string) {
	// run "cp"
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	err := cp(src, dstCp)
	assert.NoError(t, err)
	// run "CopyFile"
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	err = CopyFile(src, dstCopy)
	assert.NoError(t, err)
	// verify that both created the same files and folders
	verifyEqualFolderStructure(t, dstCp, dstCopy)
}

// verifies that the "CopyFile" function copies files the exact same way as the built-in "cp" command in Unix.
func (ws workspace) verifySameBehaviorAsCpr(t *testing.T, src, dstTemplate string) {
	// run "cp -r"
	dstCp := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCp.Path, 1)
	err := cpr(src, dstCp)
	assert.NoError(t, err)
	// run "CopyFile"
	dstCopy := strings.Replace(dstTemplate, "{{dstDir}}", ws.dstCopy.Path, 1)
	err = CopyFiles(src, dstCopy)
	assert.NoError(t, err)
	// verify that both created the same files and folders
	verifyEqualFolderStructure(t, dstCp, dstCopy)
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
	assert.Equal(t, cpEntries, copyEntries)
}
