package comments

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/tests"
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

func TestHtmlFormat_renderBlock(t *testing.T) {
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

func TestHtmlFormat_renderLineStart(t *testing.T) {
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

func TestRemove_pound(t *testing.T) {
	t.Parallel()
	give := tests.Trim(`
# Copyright © 1997 Ory Corp Inc.

# another comment

name: test
hello: world`)
	want := tests.Trim(`
# another comment

name: test
hello: world`)
	have := poundComments.remove(give, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestRemoveHtmlStyle(t *testing.T) {
	t.Parallel()
	give := "<!-- Copyright © 1997 Ory Corp Inc. -->\n<!-- All rights reserved -->\n\nname: test\nhello: world\n"
	want := "name: test\nhello: world\n"
	have := htmlComments.remove(give, "Copyright ©")
	assert.Equal(t, want, have)
}
