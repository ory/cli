// Copyright © 2022 Ory Corp

package schema

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAllFiles(t *testing.T, dir string) (files []string) {
	entries, err := os.ReadDir(dir)
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
	t.Run("case=simple snapshot", func(t *testing.T) {
		testDir, err := os.MkdirTemp("", "version-schema-test-")
		require.NoError(t, err)

		copyDir(t, "fixtures/render_version_test", testDir)

		wd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.Chdir(testDir))

		cmd := new(cobra.Command)
		err = cmd.ExecuteContext(context.Background())
		require.NoError(t, err)
		addVersionToSchema(cmd, []string{"hydra", "v1.0.0", ".schema/config.schema.json"})

		require.NoError(t, os.Chdir(wd))

		expected, err := os.ReadFile("expected/render_version_test/.schema/version.schema.json")
		require.NoError(t, err)
		actual, err := os.ReadFile(path.Join(testDir, ".schema/version.schema.json"))
		require.NoError(t, err)
		// converting to string to have nice output in case they are not equal
		assert.Equal(t, string(expected), string(actual), "To accept the new render output run:\ncp %s %s", path.Join(testDir, ".schema/version.schema.json"), path.Join(wd, "expected/render_version_test/.schema/version.schema.json"))
	})

	t.Run("case=skips pre release", func(t *testing.T) {
		cmd := &cobra.Command{}

		testDir := t.TempDir()
		copyDir(t, "fixtures/render_version_test", testDir)

		expected, err := os.ReadFile(filepath.Join("fixtures", "render_version_test", ".schema", "version.schema.json"))
		require.NoError(t, err)

		for _, version := range []string{"v1.10.5-pre.1", "v0.5.3-alpha.1.pre.0"} {
			t.Run("version="+version, func(t *testing.T) {
				stdOut := &bytes.Buffer{}
				cmd.SetOut(stdOut)

				addVersionToSchema(cmd, []string{"project-name", version})
				require.Contains(t, stdOut.String(), "is a pre release")

				actual, err := os.ReadFile(filepath.Join(testDir, ".schema", "version.schema.json"))
				require.NoError(t, err)
				assert.Equal(t, expected, actual)
			})
		}
	})
}
