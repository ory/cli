// helper functions for testing file operations, panic on error
package tests

import (
	"os"
	"path/filepath"
)

// HELPERS

// a filesystem directory used for testing
type Dir struct {
	Path string
}

func CreateTmpDir() Dir {
	path, err := os.MkdirTemp("", "ory-license")
	if err != nil {
		panic(err)
	}
	return Dir{path}
}

func (t Dir) Content(path string) string {
	content, err := os.ReadFile(filepath.Join(t.Path, path))
	if err != nil {
		panic(err)
	}
	return string(content)
}

func (t Dir) CreateDir(name string) Dir {
	t.RemoveDir(name)
	path := filepath.Join(t.Path, name)
	err := os.Mkdir(path, 0744)
	if err != nil {
		panic(err)
	}
	return Dir{path}
}

func (t Dir) CreateFile(name, content string) string {
	path := filepath.Join(t.Path, name)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, []byte(content), 0744)
	if err != nil {
		panic(err)
	}
	return path
}

func (t Dir) Filename(base string) string {
	return filepath.Join(t.Path, base)
}

func (t Dir) RemoveDir(name string) {
	os.RemoveAll(filepath.Join(t.Path, name))
}

func (t Dir) Cleanup() {
	err := os.RemoveAll(t.Path)
	if err != nil {
		panic(err)
	}
}
