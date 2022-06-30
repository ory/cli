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

func TestYmlComment(t *testing.T) {
	tests := map[string]string{
		"Hello":          "# Hello",            // single line text
		"Hello\n":        "# Hello\n",          // single line text
		"Hello\nWorld":   "# Hello\n# World",   // multi-line text
		"Hello\nWorld\n": "# Hello\n# World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := headers.YmlComment(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestRemoveHeader(t *testing.T) {
	t.Parallel()
	give := "# Copyright © 1997 Ory Corp Inc.\n\nname: test\nhello: world\n"
	want := "name: test\nhello: world\n"
	have := headers.Remove(give, headers.YmlComment, "Copyright ©")
	assert.Equal(t, want, have)
}

func TestAddLicenses(t *testing.T) {
	dir := createTmpDir()
	dir.createFile("one.yml", "one: two\nalpha: beta")
	dir.createFile("two.yml", "three: four\ngamma: delta")
	err := headers.AddLicenses(dir.path, 2022)
	assert.NoError(t, err)
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\none: two\nalpha: beta", dir.content("one.yml"))
	assert.Equal(t, "# Copyright © 2022 Ory Corp Inc.\n\nthree: four\ngamma: delta", dir.content("two.yml"))
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
