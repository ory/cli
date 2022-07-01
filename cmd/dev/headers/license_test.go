package headers_test

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/dev/headers"
	"github.com/ory/cli/cmd/dev/headers/tests"
	"github.com/stretchr/testify/assert"
)

func TestFileExt(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"one.yml":  "yml",
		"one.yaml": "yaml",
		"one.md":   "md",
		"one":      "",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := headers.FileExt(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestPrependPound(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "# Hello",            // single line text
		"Hello\n":        "# Hello\n",          // single line text
		"Hello\nWorld":   "# Hello\n# World",   // multi-line text
		"Hello\nWorld\n": "# Hello\n# World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := headers.PrependPound(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestPrependDoubleSlash(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "// Hello",             // single line text
		"Hello\n":        "// Hello\n",           // single line text
		"Hello\nWorld":   "// Hello\n// World",   // multi-line text
		"Hello\nWorld\n": "// Hello\n// World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := headers.PrependDoubleSlash(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestWrapInHtmlComment(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Hello":          "<!-- Hello -->",                   // single line text
		"Hello\n":        "<!-- Hello -->\n",                 // single line text
		"Hello\nWorld":   "<!-- Hello -->\n<!-- World -->",   // multi-line text
		"Hello\nWorld\n": "<!-- Hello -->\n<!-- World -->\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := headers.WrapInHtmlComment(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	give := "# Copyright © 1997 Ory Corp Inc.\n\nname: test\nhello: world\n"
	want := "name: test\nhello: world\n"
	have := headers.Remove(give, headers.PrependPound, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestAddLicenses(t *testing.T) {
	t.Parallel()
	dir := tests.CreateTmpDir()
	dir.CreateFile(".gitignore", "legacy.go")
	dir.CreateFile("c-sharp.cs", "using System;\n\nnamespace Foo.Bar {\n")
	dir.CreateFile("dart.dart", "int a = 1;\nint b = 2;")
	dir.CreateFile("golang.go", "package test\n\nimport foo\n")
	dir.CreateFile("java.java", "import java.io.File;\n\nFile myFile = new File();")
	dir.CreateFile("javascript.js", "const a = 1\nconst b = 2\n")
	dir.CreateFile("legacy.go", "package ignore_this_file")
	dir.CreateFile("php.php", "$a = 1;\n$b = 2;\n")
	dir.CreateFile("python.py", "a = 1\nb = 2\n")
	dir.CreateFile("ruby.rb", "a = 1\nb = 2\n")
	dir.CreateFile("rust.rs", "let a = 1;\nlet mut b = 2;\n")
	dir.CreateFile("typescript.ts", "const a = 1\nconst b = 2\n")
	dir.CreateFile("vue.vue", "<template>\n<Header />")
	dir.CreateFile("yaml.yml", "one: two\nalpha: beta")
	err := headers.AddLicenses(dir.Path, 2022)
	assert.NoError(t, err)
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nusing System;\n\nnamespace Foo.Bar {\n", dir.Content("c-sharp.cs"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nint a = 1;\nint b = 2;", dir.Content("dart.dart"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\npackage test\n\nimport foo\n", dir.Content("golang.go"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nimport java.io.File;\n\nFile myFile = new File();", dir.Content("java.java"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nconst a = 1\nconst b = 2\n", dir.Content("javascript.js"))
	assert.Equal(t, "package ignore_this_file", dir.Content("legacy.go"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\n$a = 1;\n$b = 2;\n", dir.Content("php.php"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\na = 1\nb = 2\n", dir.Content("python.py"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\na = 1\nb = 2\n", dir.Content("ruby.rb"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nlet a = 1;\nlet mut b = 2;\n", dir.Content("rust.rs"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nconst a = 1\nconst b = 2\n", dir.Content("typescript.ts"))
	assert.Equal(t, "<!-- Copyright © 2022 Ory Corp Inc. -->\n\n<template>\n<Header />", dir.Content("vue.vue")) // don't add license headers to YML files
	assert.Equal(t, "one: two\nalpha: beta", dir.Content("yaml.yml"))                                            // don't add license headers to YML files
}
