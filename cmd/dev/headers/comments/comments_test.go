package comments

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanComment(t *testing.T) {
	assert.True(t, Supports("foo.ts"))
	assert.True(t, Supports("foo.md"))
	assert.False(t, Supports("foo.xx"))
}

func TestGetFileType(t *testing.T) {
	t.Parallel()
	tests := map[string]FileType{
		"one.yml":  "yml",
		"one.yaml": "yaml",
		"one.md":   "md",
		"one":      "",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := GetFileType(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestPrependPound(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "# Hello",            // single line text
		"Hello\n":        "# Hello\n",          // single line text
		"Hello\nWorld":   "# Hello\n# World",   // multi-line text
		"Hello\nWorld\n": "# Hello\n# World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := PrependPound(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestPrependDoubleSlash(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "// Hello",             // single line text
		"Hello\n":        "// Hello\n",           // single line text
		"Hello\nWorld":   "// Hello\n// World",   // multi-line text
		"Hello\nWorld\n": "// Hello\n// World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := PrependDoubleSlash(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestWrapInHtmlComment(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "<!-- Hello -->",                   // single line text
		"Hello\n":        "<!-- Hello -->\n",                 // single line text
		"Hello\nWorld":   "<!-- Hello -->\n<!-- World -->",   // multi-line text
		"Hello\nWorld\n": "<!-- Hello -->\n<!-- World -->\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := WrapInHtmlComment(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	give := "# Copyright © 1997 Ory Corp Inc.\n\nname: test\nhello: world\n"
	want := "name: test\nhello: world\n"
	have := Remove(give, PrependPound, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestContainsFileType(t *testing.T) {
	t.Parallel()
	fileTypes := []FileType{"ts", "md", "go"}
	assert.True(t, ContainsFileType(fileTypes, "ts"))
	assert.True(t, ContainsFileType(fileTypes, "go"))
	assert.False(t, ContainsFileType(fileTypes, "rs"))
}