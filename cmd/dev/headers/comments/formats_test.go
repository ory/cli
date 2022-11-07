// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package comments

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleSlashCommentsRender(t *testing.T) {
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

func TestDoubleSlashCommentsRemove(t *testing.T) {
	t.Parallel()
	give := "// Copyright © 1997 Ory Corp Inc.\n// SPDX-License-Identifier: Apache-2.0\n\n// another comment\n\nname: test\nhello: world"
	want := "// another comment\n\nname: test\nhello: world"
	have := doubleSlashComments.remove(give, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestPoundCommentsRender(t *testing.T) {
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

func TestPoundCommentsRemove(t *testing.T) {
	t.Parallel()
	give := strings.Trim(`
# Copyright © 1997 Ory Corp Inc.

# another comment

name: test
hello: world`, "\n")
	want := strings.Trim(`
# another comment

name: test
hello: world`, "\n")
	have := poundComments.remove(give, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestHtmlCommentsRenderBlock(t *testing.T) {
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

func TestHtmlCommentsRenderLineStart(t *testing.T) {
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

func TestHtmlCommentsRemove(t *testing.T) {
	t.Parallel()
	give := strings.Trim(`
<!-- another comment -->

<!-- Copyright © 1997 Ory Corp Inc. -->
<!-- All rights reserved -->

name: test
hello: world`, "\n")
	want := strings.Trim(`
<!-- another comment -->

name: test
hello: world`, "\n")
	have := htmlComments.remove(give, "Copyright ©")
	assert.Equal(t, want, have)
}
