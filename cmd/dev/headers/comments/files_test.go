package comments_test

import (
	"os"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/stretchr/testify/assert"
)

func TestFileContentWithoutHeader_knownFile(t *testing.T) {
	err := os.WriteFile("testfile.md", []byte("<!-- copyright Ory -->\n<!-- all rights reserved -->\n\nhello world"), 0744)
	defer os.Remove("testfile.md")
	assert.NoError(t, err)
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	want := "hello world"
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_otherCommentFirst(t *testing.T) {
	err := os.WriteFile("testfile.md", []byte("<!-- another comment -->\n\n<!-- copyright Ory -->\n<!-- all rights reserved -->\n\nhello world"), 0744)
	defer os.Remove("testfile.md")
	assert.NoError(t, err)
	have, err := comments.FileContentWithoutHeader("testfile.md", "copyright")
	want := "<!-- another comment -->\n\nhello world"
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_unknownFile(t *testing.T) {
	err := os.WriteFile("testfile.txt", []byte("hello world"), 0744)
	defer os.Remove("testfile.txt")
	assert.NoError(t, err)
	have, err := comments.FileContentWithoutHeader("testfile.txt", "copyright")
	want := "hello world"
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}
