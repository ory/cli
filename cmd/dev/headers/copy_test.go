package headers_test

import (
	"os"
	"testing"

	"github.com/ory/cli/cmd/dev/headers"
	"github.com/stretchr/testify/assert"
)

func TestCopy_singleFile_fullDestPath(t *testing.T) {
	t.Parallel()
	rootDir := testDir{path: "."}
	srcDir := rootDir.createDir("test_copy_src")
	dstDir := rootDir.createDir("test_copy_dst")
	srcDir.createFile("README.md", "# the readme\ntext")
	err := headers.CopyFile("test_copy_src/README.md", dstDir.filename("README.md"))
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/test_copy_src/README.md. -->\n\n# the readme\ntext",
		dstDir.content("README.md"))
	os.RemoveAll(srcDir.path)
	os.RemoveAll(dstDir.path)
}
