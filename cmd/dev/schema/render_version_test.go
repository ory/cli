package schema

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/containerd/continuity/fs"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddVersionToSchema(t *testing.T) {
	t.Run("case=simple snapshot", func(t *testing.T) {
		testDir, err := ioutil.TempDir("", "version-schema-test-")
		require.NoError(t, err)

		require.NoError(t, fs.CopyDir("fixtures/render_version_test", testDir))

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
	})

	t.Run("case=skips pre release", func(t *testing.T) {
		cmd := &cobra.Command{}

		testDir := t.TempDir()
		require.NoError(t, fs.CopyDir(testDir, "fixtures/render_version_test"))

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
