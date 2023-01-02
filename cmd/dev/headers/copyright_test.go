// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package headers

import (
	"fmt"
	"regexp"
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
				err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
				assert.NoError(t, err)
				assert.Equal(t, test.want, dir.Content(filename))
			})
		}
	})

	t.Run("does not add a header to files in .prettierignore", func(t *testing.T) {
		dir := root.CreateDir("prettierignored")
		dir.CreateFile(".prettierignore", "prettier-ignored.go")
		content := "package ignore_this_file"
		dir.CreateFile("prettier-ignored.go", content)
		err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, content, dir.Content("prettier-ignored.go"))
	})

	t.Run("does not add a header to files ignored by path in .prettierignore", func(t *testing.T) {
		dir := root.CreateDir("prettierignored")
		subdir := dir.CreateDir("subdir")
		dir.CreateFile(".prettierignore", "subdir/")
		content := "package ignore_this_file"
		subdir.CreateFile("prettier-ignored.go", content)
		err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, content, subdir.Content("prettier-ignored.go"))
	})

	t.Run("does not add a header to files ignored by path in .gitignore", func(t *testing.T) {
		dir := root.CreateDir("gitignored")
		subdir := dir.CreateDir("subdir")
		dir.CreateFile(".gitignore", "subdir/")
		content := "package ignore_this_file"
		subdir.CreateFile("git-ignored.go", content)
		err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, content, subdir.Content("git-ignored.go"))
	})

	t.Run("does not add a header to files in .gitignore", func(t *testing.T) {
		dir := root.CreateDir("gitignored")
		dir.CreateFile(".gitignore", "git-ignored.go")
		dir.CreateFile("git-ignored.go", "package ignore_this_file")
		err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, "package ignore_this_file", dir.Content("git-ignored.go"))
	})

	t.Run("does not add a header to files in node_modules", func(t *testing.T) {
		dir := root.CreateDir("node_modules").CreateDir(".bin")
		content := "#!/usr/bin/env bash\necho hello"
		dir.CreateFile("nodemon.ts", content)
		err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, content, dir.Content("nodemon.ts"))
	})

	t.Run("does not add a header to files in the given `exclude` argument", func(t *testing.T) {
		dir1 := root.CreateDir("excluded")
		dir2 := dir1.CreateDir("generated")
		content := "package this_file_is_excluded"
		dir2.CreateFile("excluded.go", content)
		err := AddHeaders(dir1.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{"generated"}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, content, dir2.Content("excluded.go"))
	})

	t.Run("does not change copyright year for existing header", func(t *testing.T) {
		type test struct {
			ext  string
			want string
		}
		tests := []test{
			{ext: "cs", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nusing System;\n\nnamespace Foo.Bar {\n"},
			{ext: "dart", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nint a = 1;\nint b = 2;"},
			{ext: "go", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\npackage test\n\nimport foo\n"},
			{ext: "java", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nimport java.io.File;\n\nFile myFile = new File();"},
			{ext: "js", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nconst a = 1\nconst b = 2\n"},
			{ext: "md", want: "# hello world"},
			{ext: "php", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\n$a = 1;\n$b = 2;\n"},
			{ext: "py", want: "# Copyright © 2022 Ory Corp\n# SPDX-License-Identifier: Apache-2.0\n\na = 1\nb = 2\n"},
			{ext: "rb", want: "# Copyright © 2022 Ory Corp\n# SPDX-License-Identifier: Apache-2.0\n\na = 1\nb = 2\n"},
			{ext: "rs", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nlet a = 1;\nlet mut b = 2;\n"},
			{ext: "ts", want: "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\nconst a = 1\nconst b = 2\n"},
			{ext: "vue", want: "<!-- Copyright © 2022 Ory Corp -->\n<!-- SPDX-License-Identifier: Apache-2.0 -->\n\n<template>\n<Header />"},
			{ext: "yml", want: "one: two\nalpha: beta"},
			{ext: "yaml", want: "one: two\nalpha: beta"},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("%q file type", test.ext), func(t *testing.T) {
				dir := root.CreateDir(fmt.Sprintf("no-update-%s", test.ext))
				filename := fmt.Sprintf("file.%s", test.ext)
				dir.CreateFile(filename, test.want)
				err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2023), []string{}, regexp.MustCompile(HEADER_REGEXP))
				assert.NoError(t, err)
				assert.Equal(t, test.want, dir.Content(filename))
			})
		}
	})

	t.Run("open-source copyright headers", func(t *testing.T) {
		dir := root.CreateDir("open-source")
		dir.CreateFile("file.go", "package open_source")
		err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_OPEN_SOURCE, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, "// Copyright © 2022 Ory Corp\n// SPDX-License-Identifier: Apache-2.0\n\npackage open_source", dir.Content("file.go"))
	})

	t.Run("proprietary copyright headers", func(t *testing.T) {
		dir := root.CreateDir("open-source")
		dir.CreateFile("file.go", "package open_source")
		err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_PROPRIETARY, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
		assert.NoError(t, err)
		assert.Equal(t, "// Copyright © 2022 Ory Corp\n// Proprietary and confidential.\n// Unauthorized copying of this file is prohibited.\n\npackage open_source", dir.Content("file.go"))
	})

	t.Run("correctly handles files with BOM", func(t *testing.T) {
		type test struct {
			bom  string
			name string
		}
		for _, test := range []test{
			{bom: "\xef\xbb\xbf", name: "UTF-8"},
			{bom: "\ufffe", name: "UTF-16 (LE)"},
			{bom: "\ufeff", name: "UTF-16 (BE)"},
			{bom: "\ufffe\x00\x00", name: "UTF-32 (LE)"},
			{bom: "\x00\x00\ufeff", name: "UTF-32 (BE)"},
		} {
			t.Run(fmt.Sprintf("%s BOM", test.name), func(t *testing.T) {
				dir := root.CreateDir("bom-test")

				content := "package open_source"
				dir.CreateFile("file.go", fmt.Sprintf("%s%s", test.bom, content))
				err := AddHeaders(dir.Path, fmt.Sprintf(HEADER_TEMPLATE_PROPRIETARY, 2022), []string{}, regexp.MustCompile(HEADER_REGEXP))
				assert.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("%s// Copyright © 2022 Ory Corp\n// Proprietary and confidential.\n// Unauthorized copying of this file is prohibited.\n\n%s", test.bom, content), dir.Content("file.go"))
			})
		}
	})
}

func TestPathContainsFolders(t *testing.T) {
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
		assert.Equal(t, want, pathContainsFolders(give, exclude), "%q -> %t", give, want)
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
