// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package comments_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/cli/cmd/dev/headers/comments"
)

func TestFileContentWithoutHeader(t *testing.T) {
	t.Run("known file, copyright header", func(t *testing.T) {
		give := strings.Trim(`
file content`, "\n")
		want := strings.Trim(`
file content`, "\n")
		createTestFile(t, "testfile.go", give)
		defer os.Remove("testfile.go")
		have, err := comments.FileContentWithoutHeader("testfile.go", "Copyright ©")
		assert.NoError(t, err)
		assert.Equal(t, want, have)
	})

	t.Run("known file, other comment first", func(t *testing.T) {
		give := strings.Trim(`
// another comment block

file content`, "\n")
		want := strings.Trim(`
// another comment block

file content`, "\n")
		createTestFile(t, "testfile.go", give)
		defer os.Remove("testfile.go")
		have, err := comments.FileContentWithoutHeader("testfile.go", "Copyright ©")
		assert.NoError(t, err)
		assert.Equal(t, want, have)
	})

	t.Run("unknown file", func(t *testing.T) {
		give := strings.Trim(`
file content`, "\n")
		want := give
		createTestFile(t, "testfile.txt", give)
		defer os.Remove("testfile.txt")
		have, err := comments.FileContentWithoutHeader("testfile.txt", "Copyright ©")
		assert.NoError(t, err)
		assert.Equal(t, want, have)
	})
}

func createTestFile(t *testing.T, name, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0744)
	assert.NoError(t, err)
}
