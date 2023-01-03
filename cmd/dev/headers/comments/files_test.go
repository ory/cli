// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package comments_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/cli/cmd/dev/headers/comments"
)

func TestFileContentWithoutHeader(t *testing.T) {
	t.Run("known file, copyright header, empty line", func(t *testing.T) {
		give := "// Copyright © 2021 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nfile content"
		want := "file content"
		createTestFile(t, "testfile.go", give)
		defer os.Remove("testfile.go")
		have, err := comments.FileContentWithoutHeader("testfile.go", regexp.MustCompile(`Copyright © \d{4} Ory Corp`))
		assert.NoError(t, err)
		assert.Equal(t, want, have)
	})

	t.Run("known file, copyright header, no empty line", func(t *testing.T) {
		give := "// Copyright © 2021 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\nfile content"
		want := "file content"
		createTestFile(t, "testfile.go", give)
		defer os.Remove("testfile.go")
		have, err := comments.FileContentWithoutHeader("testfile.go", regexp.MustCompile(`Copyright © \d{4} Ory Corp`))
		assert.NoError(t, err)
		assert.Equal(t, want, have)
	})

	t.Run("known file, other comment first", func(t *testing.T) {
		give := "// another comment block\n\n// Copyright © 2021 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nfile content"
		want := "// another comment block\n\nfile content"
		createTestFile(t, "testfile.go", give)
		defer os.Remove("testfile.go")
		have, err := comments.FileContentWithoutHeader("testfile.go", regexp.MustCompile(`Copyright © \d{4} Ory Corp`))
		assert.NoError(t, err)
		assert.Equal(t, want, have)
	})

	t.Run("unknown file", func(t *testing.T) {
		give := "// Copyright © 2021 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nfile content"
		want := give
		createTestFile(t, "testfile.txt", give)
		defer os.Remove("testfile.txt")
		have, err := comments.FileContentWithoutHeader("testfile.txt", regexp.MustCompile(`Copyright © \d{4} Ory Corp`))
		assert.NoError(t, err)
		assert.Equal(t, want, have)
	})
}

func createTestFile(t *testing.T, name, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0744)
	assert.NoError(t, err)
}
