package comments_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/cli/cmd/dev/headers/comments"
)

func TestFileContentWithoutHeader_knownFile(t *testing.T) {
	give := strings.Trim(`
// copyright Ory
// all rights reserved

file content`, "\n")
	want := strings.Trim(`
file content`, "\n")
	createTestFile(t, "testfile.go", give)
	defer os.Remove("testfile.go")
	have, err := comments.FileContentWithoutHeader("testfile.go", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_otherCommentFirst(t *testing.T) {
	give := strings.Trim(`
// another comment block

// copyright Ory
// all rights reserved

file content`, "\n")
	want := strings.Trim(`
// another comment block

file content`, "\n")
	createTestFile(t, "testfile.go", give)
	defer os.Remove("testfile.go")
	have, err := comments.FileContentWithoutHeader("testfile.go", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_unknownFile(t *testing.T) {
	give := "file content"
	want := "file content"
	createTestFile(t, "testfile.txt", give)
	defer os.Remove("testfile.txt")
	have, err := comments.FileContentWithoutHeader("testfile.txt", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func createTestFile(t *testing.T, name, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0744)
	assert.NoError(t, err)
}
