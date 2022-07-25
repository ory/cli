package comments

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	have := remove(give, poundComments, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestRemoveHtmlStyle(t *testing.T) {
	t.Parallel()
	give := "<!-- Copyright © 1997 Ory Corp Inc. -->\n<!-- All rights reserved -->\n\nname: test\nhello: world\n"
	want := "name: test\nhello: world\n"
	have := remove(give, htmlComments, "Copyright ©")
	assert.Equal(t, want, have)
}
