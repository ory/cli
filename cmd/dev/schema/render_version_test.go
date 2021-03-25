package schema

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAllFiles(t *testing.T, dir string) (files []string) {
	entries, err := ioutil.ReadDir(dir)
	require.NoError(t, err)

	for _, e := range entries {
		switch {
		case e.IsDir():
			files = append(files, getAllFiles(t, path.Join(dir, e.Name()))...)
		default:
			files = append(files, path.Join(dir, e.Name()))
		}
	}
	return
}

func copyDir(t *testing.T, src, dst string) {
	files := getAllFiles(t, src)
	for _, fn := range files {
		srcF, err := os.Open(fn)
		require.NoError(t, err)
		dstFn := path.Join(dst, strings.TrimLeft(fn, src))
		require.NoError(t, os.MkdirAll(path.Dir(dstFn), 0777))
		dstF, err := os.OpenFile(dstFn, os.O_WRONLY|os.O_CREATE, 0666)
		require.NoError(t, err)
		_, err = io.Copy(dstF, srcF)
		require.NoError(t, err)
		require.NoError(t, srcF.Close())
		require.NoError(t, dstF.Close())
	}
}

func TestAddVersionToSchema(t *testing.T) {
	testDir, err := ioutil.TempDir("", "version-schema-test-")
	require.NoError(t, err)

	copyDir(t, "fixtures/render_version_test", testDir)

	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(testDir))

	addVersionToSchema(nil, []string{"hydra", "v1.0.0", ".schema/config.schema.json"})

	require.NoError(t, os.Chdir(wd))

	expected, err := ioutil.ReadFile("expected/render_version_test/.schema/version.schema.json")
	require.NoError(t, err)
	actual, err := ioutil.ReadFile(path.Join(testDir, ".schema/version.schema.json"))
	require.NoError(t, err)
	// converting to string to have nice output in case they are not equal
	assert.Equal(t, string(expected), string(actual), "To accept the new render output run:\ncp %s %s", path.Join(testDir, ".schema/version.schema.json"), path.Join(wd, "expected/render_version_test/.schema/version.schema.json"))
}
