package comments_test

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/stretchr/testify/assert"
)

func TestGetFileType(t *testing.T) {
	t.Parallel()
	tests := map[string]comments.FileType{
		"one.yml":  "yml",
		"one.yaml": "yaml",
		"one.md":   "md",
		"one":      "",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := comments.GetFileType(give)
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
			have := comments.PrependPound(give)
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
			have := comments.PrependDoubleSlash(give)
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
			have := comments.WrapInHtmlComment(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	give := "# Copyright © 1997 Ory Corp Inc.\n\nname: test\nhello: world\n"
	want := "name: test\nhello: world\n"
	have := comments.Remove(give, comments.PrependPound, "Copyright ©")
	assert.Equal(t, want, have)
}
