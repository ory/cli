package headers

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/dev/headers/tests"
	"github.com/stretchr/testify/assert"
)

func TestAddLicenses(t *testing.T) {
	t.Parallel()
	dir := tests.CreateTmpDir()
	dir.CreateFile(".gitignore", "git-ignored.go")
	dir.CreateFile("c-sharp.cs", "using System;\n\nnamespace Foo.Bar {\n")
	dir.CreateFile("dart.dart", "int a = 1;\nint b = 2;")
	dir.CreateFile("golang.go", "package test\n\nimport foo\n")
	dir.CreateFile("java.java", "import java.io.File;\n\nFile myFile = new File();")
	dir.CreateFile("javascript.js", "const a = 1\nconst b = 2\n")
	dir.CreateFile("git-ignored.go", "package ignore_this_file")
	dir.CreateFile("markdown.md", "# hello world")
	dir.CreateFile("php.php", "$a = 1;\n$b = 2;\n")
	dir.CreateFile("python.py", "a = 1\nb = 2\n")
	dir.CreateFile("ruby.rb", "a = 1\nb = 2\n")
	dir.CreateFile("rust.rs", "let a = 1;\nlet mut b = 2;\n")
	dir.CreateFile("typescript.ts", "const a = 1\nconst b = 2\n")
	dir.CreateFile("vue.vue", "<template>\n<Header />")
	dir.CreateFile("yaml.yml", "one: two\nalpha: beta")
	err := AddLicenses(dir.Path, 2022)
	assert.NoError(t, err)
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nusing System;\n\nnamespace Foo.Bar {\n", dir.Content("c-sharp.cs"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nint a = 1;\nint b = 2;", dir.Content("dart.dart"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\npackage test\n\nimport foo\n", dir.Content("golang.go"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nimport java.io.File;\n\nFile myFile = new File();", dir.Content("java.java"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nconst a = 1\nconst b = 2\n", dir.Content("javascript.js"))
	assert.Equal(t, "package ignore_this_file", dir.Content("git-ignored.go"))
	assert.Equal(t, "# hello world", dir.Content("markdown.md"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\n$a = 1;\n$b = 2;\n", dir.Content("php.php"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\na = 1\nb = 2\n", dir.Content("python.py"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\na = 1\nb = 2\n", dir.Content("ruby.rb"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nlet a = 1;\nlet mut b = 2;\n", dir.Content("rust.rs"))
	assert.Equal(t, "// Copyright © 2022 Ory Corp Inc.\n\nconst a = 1\nconst b = 2\n", dir.Content("typescript.ts"))
	assert.Equal(t, "<!-- Copyright © 2022 Ory Corp Inc. -->\n\n<template>\n<Header />", dir.Content("vue.vue"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\none: two\nalpha: beta", dir.Content("yaml.yml"))
}

func TestShouldAddLicense(t *testing.T) {
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
		"x.yml":  true,
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %t", give, want), func(t *testing.T) {
			assert.Equal(t, want, shouldAddLicense(give))
		})
	}
}
