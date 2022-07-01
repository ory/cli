package headers_test

import (
	"testing"

	"github.com/ory/cli/cmd/dev/headers"
	"github.com/stretchr/testify/assert"
)

func TestCopy_singleFile(t *testing.T) {
	t.Parallel()
	rootDir := createTmpDir()
	srcDir := rootDir.createDir("source")
	dstDir := rootDir.createDir("dst")
	readme := srcDir.createFile("README.md", "# the readme\ntext")
	err := headers.CopyFile(readme, dstDir.filename("README.md"))
	assert.NoError(t, err)
	assert.Equal(
		t,
		"<!-- AUTO-GENERATED, DO NOT EDIT! Please edit the original at https://github.com/ory/meta/blob/master/source/README.md -->\n#the readme\ntext",
		dstDir.content("README.md"))
}
