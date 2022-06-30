package headers_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ory/cli/cmd/dev/headers"
	"github.com/stretchr/testify/assert"
)

func TestFileExt(t *testing.T) {
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
	dir := createTmpDir()
	dir.createFile("c-sharp.cs", "using System;\n\nnamespace Foo.Bar {\n")
	dir.createFile("dart.dart", "int a = 1;\nint b = 2;")
	dir.createFile("golang.go", "package test\n\nimport foo\n")
	dir.createFile("java.java", "import java.io.File;\n\nFile myFile = new File();")
	dir.createFile("javascript.js", "const a = 1\nconst b = 2\n")
	dir.createFile("php.php", "$a = 1;\n$b = 2;\n")
	dir.createFile("python.py", "a = 1\nb = 2\n")
	dir.createFile("ruby.rb", "a = 1\nb = 2\n")
	dir.createFile("rust.rs", "let a = 1;\nlet mut b = 2;\n")
	dir.createFile("typescript.ts", "const a = 1\nconst b = 2\n")
	dir.createFile("vue.vue", "<template>\n<Header />")
	dir.createFile("yaml.yml", "one: two\nalpha: beta")
	err := headers.AddLicenses(dir.path, 2022)
	assert.NoError(t, err)
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nusing System;\n\nnamespace Foo.Bar {\n", dir.content("c-sharp.cs"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nint a = 1;\nint b = 2;", dir.content("dart.dart"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\npackage test\n\nimport foo\n", dir.content("golang.go"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nimport java.io.File;\n\nFile myFile = new File();", dir.content("java.java"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nconst a = 1\nconst b = 2\n", dir.content("javascript.js"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\n$a = 1;\n$b = 2;\n", dir.content("php.php"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\na = 1\nb = 2\n", dir.content("python.py"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\na = 1\nb = 2\n", dir.content("ruby.rb"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nlet a = 1;\nlet mut b = 2;\n", dir.content("rust.rs"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nconst a = 1\nconst b = 2\n", dir.content("typescript.ts"))
	assert.Equal(t, "<!-- Copyright © 2022 Ory Corp Inc. -->\n\n<template>\n<Header />", dir.content("vue.vue")) // don't add license headers to YML files
	assert.Equal(t, "one: two\nalpha: beta", dir.content("yaml.yml"))                                            // don't add license headers to YML files
}

// HELPERS

// a directory used for testing, no need to clean up
type testDir struct {
	path string
}

func createTmpDir() testDir {
	path, err := ioutil.TempDir("", "ory-license")
	if err != nil {
		panic(err)
	}
	return testDir{path}
}

func (t testDir) createFile(name, content string) {
	err := os.WriteFile(filepath.Join(t.path, name), []byte(content), 0744)
	if err != nil {
		panic(err)
	}
}

func (t testDir) content(path string) string {
	content, err := os.ReadFile(filepath.Join(t.path, path))
	if err != nil {
		panic(err)
	}
	return string(content)
}
