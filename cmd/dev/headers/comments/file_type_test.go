package comments_test

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/stretchr/testify/assert"
)

func TestContainsFileType(t *testing.T) {
	t.Parallel()
	fileTypes := []comments.FileType{"ts", "md", "go"}
	assert.True(t, comments.ContainsFileType(fileTypes, "ts"))
	assert.True(t, comments.ContainsFileType(fileTypes, "go"))
	assert.False(t, comments.ContainsFileType(fileTypes, "rs"))
}

func TestGetFileType(t *testing.T) {
	t.Parallel()
	tests := map[string]comments.FileType{
		"one.yml":  "yml",
		"one.yaml": "yaml",
		"one.md":   "md",
		"one.xx":   "xx",
		"one":      "",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := comments.GetFileType(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestSupports(t *testing.T) {
	assert.True(t, comments.SupportsFile("foo.ts"))
	assert.True(t, comments.SupportsFile("foo.md"))
	assert.False(t, comments.SupportsFile("foo.xx"))
}
