// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package comments_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/cli/cmd/dev/headers/comments"
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
		"foo.yml":  "yml",
		"foo.yaml": "yml",
		"foo.md":   "md",
		"foo.xxx":  "xxx",
		"foo":      "",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			have := comments.GetFileType(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestSupports(t *testing.T) {
	t.Parallel()
	assert.True(t, comments.SupportsFile("foo.ts"))
	assert.True(t, comments.SupportsFile("foo.md"))
	assert.False(t, comments.SupportsFile("foo.xxx"))
}
