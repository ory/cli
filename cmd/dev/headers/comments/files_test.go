package comments_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/stretchr/testify/assert"
)

func TestFileContentWithoutHeader_knownFile(t *testing.T) {
	give := `
<!-- copyright Ory -->
<!-- all rights reserved -->

file content`
	want := `
file content`
	createTestFile(t, "testfile.md", trim(give))
	defer os.Remove("testfile.md")
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, trim(want), have)
}

func TestFileContentWithoutHeader_otherCommentFirst(t *testing.T) {
	give := `
<!-- another comment block -->

<!-- copyright Ory -->
<!-- all rights reserved -->

file content`
	want := `
<!-- another comment block -->

file content`
	createTestFile(t, "testfile.md", trim(give))
	defer os.Remove("testfile.md")
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, trim(want), have)
}

func TestFileContentWithoutHeader_unknownFile(t *testing.T) {
	give := "file content"
	want := "file content"
	createTestFile(t, "testfile.txt", trim(give))
	defer os.Remove("testfile.txt")
	have, err := comments.FileContentWithoutHeader("testfile.txt", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, trim(want), have)
}

// Helper function that indicates within a test that we create a test file
// and encapsulates incidental complexity.
func createTestFile(t *testing.T, name, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0744)
	assert.NoError(t, err)
}

func trim(text string) string {
	return strings.Trim(text, "\n")
}
