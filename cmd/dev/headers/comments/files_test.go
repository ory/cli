package comments_test

import (
	"os"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/ory/cli/cmd/dev/headers/tests"
	"github.com/stretchr/testify/assert"
)

func TestFileContentWithoutHeader_knownFile(t *testing.T) {
	give := tests.Trim(`
<!-- copyright Ory -->
<!-- all rights reserved -->

file content`)
	want := tests.Trim(`
file content`)
	tests.CreateTestFile(t, "testfile.md", give)
	defer os.Remove("testfile.md")
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_otherCommentFirst(t *testing.T) {
	give := tests.Trim(`
<!-- another comment block -->

<!-- copyright Ory -->
<!-- all rights reserved -->

file content`)
	want := tests.Trim(`
<!-- another comment block -->

file content`)
	tests.CreateTestFile(t, "testfile.md", give)
	defer os.Remove("testfile.md")
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_unknownFile(t *testing.T) {
	give := "file content"
	want := "file content"
	tests.CreateTestFile(t, "testfile.txt", give)
	defer os.Remove("testfile.txt")
	have, err := comments.FileContentWithoutHeader("testfile.txt", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}
