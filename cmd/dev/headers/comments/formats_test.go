package comments

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleSlashFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "// Hello",
		"Hello\n":        "// Hello\n",
		"Hello\nWorld":   "// Hello\n// World",
		"Hello\nWorld\n": "// Hello\n// World\n",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := doubleSlashComments.renderBlock(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestPoundFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "# Hello",
		"Hello\n":        "# Hello\n",
		"Hello\nWorld":   "# Hello\n# World",
		"Hello\nWorld\n": "# Hello\n# World\n",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := poundComments.renderBlock(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestHtmlFormat_render(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "<!-- Hello -->",
		"Hello\n":        "<!-- Hello -->\n",
		"Hello\nWorld":   "<!-- Hello -->\n<!-- World -->",
		"Hello\nWorld\n": "<!-- Hello -->\n<!-- World -->\n",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := htmlComments.renderBlock(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestHtmlFormat_renderStart(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":   "<!-- Hello",
		"Hello\n": "<!-- Hello\n",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := htmlComments.renderLineStart(give)
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
