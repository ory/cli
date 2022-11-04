// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHeaders(t *testing.T) {
	root := CreateTmpDir()
	defer root.Delete()

	t.Run("adds the given comment in a format matching the file type", func(t *testing.T) {
		type test struct {
			ext  string // extension of the file type to test
			give string // file content before `AddHeaders` runs
			want string // expected file content after `AddHeaders` runs
		}
		tests := []test{
			{ext: "cs", give: "using System;\n\nnamespace Foo.Bar {\n", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nusing System;\n\nnamespace Foo.Bar {\n"},
			{ext: "dart", give: "int a = 1;\nint b = 2;", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nint a = 1;\nint b = 2;"},
			{ext: "go", give: "package test\n\nimport foo\n", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\npackage test\n\nimport foo\n"},
			{ext: "java", give: "import java.io.File;\n\nFile myFile = new File();", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nimport java.io.File;\n\nFile myFile = new File();"},
			{ext: "js", give: "const a = 1\nconst b = 2\n", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nconst a = 1\nconst b = 2\n"},
			{ext: "md", give: "# hello world", want: "# hello world"},
			{ext: "php", give: "$a = 1;\n$b = 2;\n", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\n$a = 1;\n$b = 2;\n"},
			{ext: "py", give: "a = 1\nb = 2\n", want: "# Copyright © 2022 Ory Corp\n# SPDX-License-Identifier: Apache-2.0\n\na = 1\nb = 2\n"},
			{ext: "rb", give: "a = 1\nb = 2\n", want: "# Copyright © 2022 Ory Corp\n# SPDX-License-Identifier: Apache-2.0\n\na = 1\nb = 2\n"},
			{ext: "rs", give: "let a = 1;\nlet mut b = 2;\n", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nlet a = 1;\nlet mut b = 2;\n"},
			{ext: "ts", give: "const a = 1\nconst b = 2\n", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nconst a = 1\nconst b = 2\n"},
			{ext: "vue", give: "<template>\n<Header />", want: "<!-- Copyright © 2022 Ory Corp -->\n<!-- SPDX-License-Identifier: Apache-2.0 -->\n\n<template>\n<Header />"},
			{ext: "yml", give: "one: two\nalpha: beta", want: "one: two\nalpha: beta"},
			{ext: "yaml", give: "one: two\nalpha: beta", want: "one: two\nalpha: beta"},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("%q file type", test.ext), func(t *testing.T) {
				dir := root.CreateDir(test.ext)
				filename := fmt.Sprintf("file.%s", test.ext)
				dir.CreateFile(filename, test.give)
				err := AddHeaders(dir.Path, 2022, []string{})
				assert.NoError(t, err)
				assert.Equal(t, test.want, dir.Content(filename))
			})
		}
	})

	t.Run("ignores files in .gitignore", func(t *testing.T) {
		dir := root.CreateDir("gitignored")
		dir.CreateFile(".gitignore", "git-ignored.go")
		dir.CreateFile("git-ignored.go", "package ignore_this_file")
		err := AddHeaders(dir.Path, 2022, []string{})
		assert.NoError(t, err)
		assert.Equal(t, "package ignore_this_file", dir.Content("git-ignored.go"))
	})

	t.Run("ignores files in the given `exclude` parameter", func(t *testing.T) {
		//
	})
}

func TestIsInFolders(t *testing.T) {
	exclude := []string{"internal/httpclient", "generated/"}
	tests := map[string]bool{
		"foo.md":                                false,
		"foo/bar/baz.md":                        false,
		"internal/httpclient/README.md":         true,
		"internal/httpclient/foo/bar/README.md": true,
		"generated/README.md":                   true,
		"generated/foo/bar/README.md":           true,
	}
	for give, want := range tests {
		assert.Equal(t, want, isInFolders(give, exclude), "%q -> %t", give, want)
	}
}

func TestFileTypeNeedsCopyrightHeader(t *testing.T) {
	tests := map[string]bool{
		"x.cs":   true,
		"x.dart": true,
		"x.go":   true,
		"x.java": true,
		"x.js":   true,
		"x.md":   false,
		"x.php":  true,
		"x.py":   true,
		"x.rb":   true,
		"x.rs":   true,
		"x.ts":   true,
		"x.vue":  true,
		"x.yml":  false, // data is not protected by copyright law
		"x.yaml": false, // data is not protected by copyright law
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %t", give, want), func(t *testing.T) {
			assert.Equal(t, want, fileTypeNeedsCopyrightHeader(give))
		})
	}
}
