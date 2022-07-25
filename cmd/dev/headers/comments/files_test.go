package comments_test

import (
	"os"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/stretchr/testify/assert"
)

func TestFileContentWithoutHeader_knownFile(t *testing.T) {
	give := `
<!-- copyright Ory -->
<!-- all rights reserved -->

hello world`
	want := `
hello world`
	createTestFile(t, "testfile.md", give)
	defer os.Remove("testfile.md")
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_otherCommentFirst(t *testing.T) {
	give := `
<!-- another comment -->

<!-- copyright Ory -->
<!-- all rights reserved -->

hello world`
	want := `
<!-- another comment -->

hello world`
	createTestFile(t, "testfile.md", give)
	defer os.Remove("testfile.md")
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_unknownFile(t *testing.T) {
	give := "hello world"
	want := "hello world"
	createTestFile(t, "testfile.txt", give)
	defer os.Remove("testfile.txt")
	have, err := comments.FileContentWithoutHeader("testfile.txt", "copyright")
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

// Helper function that indicates within a test that we create a test file
// and encapsulates incidental complexity.
func createTestFile(t *testing.T, name, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0744)
	assert.NoError(t, err)
}
