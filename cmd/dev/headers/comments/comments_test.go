package comments

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsFileType(t *testing.T) {
	t.Parallel()
	fileTypes := []FileType{"ts", "md", "go"}
	assert.True(t, ContainsFileType(fileTypes, "ts"))
	assert.True(t, ContainsFileType(fileTypes, "go"))
	assert.False(t, ContainsFileType(fileTypes, "rs"))
}

func TestFileContentWithoutHeader_knownFile(t *testing.T) {
	err := os.WriteFile("testfile.md", []byte("<!-- copyright Ory -->\n<!-- all rights reserved -->\n\nhello world"), 0744)
	defer os.Remove("testfile.md")
	assert.NoError(t, err)
	have, err := FileContentWithoutHeader("testfile.md", "copyright")
	want := "hello world"
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_otherCommentFirst(t *testing.T) {
	err := os.WriteFile("testfile.md", []byte("<!-- another comment -->\n\n<!-- copyright Ory -->\n<!-- all rights reserved -->\n\nhello world"), 0744)
	defer os.Remove("testfile.md")
	assert.NoError(t, err)
	have, err := FileContentWithoutHeader("testfile.md", "copyright")
	want := "<!-- another comment -->\n\nhello world"
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestFileContentWithoutHeader_unknownFile(t *testing.T) {
	err := os.WriteFile("testfile.txt", []byte("hello world"), 0744)
	defer os.Remove("testfile.txt")
	assert.NoError(t, err)
	have, err := FileContentWithoutHeader("testfile.txt", "copyright")
	want := "hello world"
	assert.NoError(t, err)
	assert.Equal(t, want, have)
}

func TestGetFileType(t *testing.T) {
	t.Parallel()
	tests := map[string]FileType{
		"one.yml":  "yml",
		"one.yaml": "yaml",
		"one.md":   "md",
		"one.xx":   "xx",
		"one":      "",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := GetFileType(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestDoubleSlashFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "// Hello",             // single line text
		"Hello\n":        "// Hello\n",           // single line text
		"Hello\nWorld":   "// Hello\n// World",   // multi-line text
		"Hello\nWorld\n": "// Hello\n// World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := doubleSlashComments.render(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestPoundFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "# Hello",            // single line text
		"Hello\n":        "# Hello\n",          // single line text
		"Hello\nWorld":   "# Hello\n# World",   // multi-line text
		"Hello\nWorld\n": "# Hello\n# World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := poundComments.render(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestHtmlFormat_render(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "<!-- Hello -->",                   // single line text
		"Hello\n":        "<!-- Hello -->\n",                 // single line text
		"Hello\nWorld":   "<!-- Hello -->\n<!-- World -->",   // multi-line text
		"Hello\nWorld\n": "<!-- Hello -->\n<!-- World -->\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := htmlComments.render(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestHtmlFormat_renderStart(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "<!-- Hello",               // single line text
		"Hello\n":        "<!-- Hello\n",             // single line text
		"Hello\nWorld":   "<!-- Hello\n<!-- World",   // multi-line text
		"Hello\nWorld\n": "<!-- Hello\n<!-- World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := htmlComments.renderStart(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestRemove_pound_beginning(t *testing.T) {
	t.Parallel()
	give := "# Copyright © 1997 Ory Corp Inc.\n\n# another comment\n\nname: test\nhello: world\n"
	want := "# another comment\n\nname: test\nhello: world\n"
	have := remove(give, prependPound, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestRemoveHtmlStyle(t *testing.T) {
	t.Parallel()
	give := "<!-- Copyright © 1997 Ory Corp Inc. -->\n<!-- All rights reserved -->\n\nname: test\nhello: world\n"
	want := "name: test\nhello: world\n"
	have := remove(give, prependHtmlComment, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestSupports(t *testing.T) {
	assert.True(t, Supports("foo.ts"))
	assert.True(t, Supports("foo.md"))
	assert.False(t, Supports("foo.xx"))
}
