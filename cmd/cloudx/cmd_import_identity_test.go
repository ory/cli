package cloudx

import (
	"bytes"
	"github.com/ory/x/cmdx"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func makeRandomIdentity(t *testing.T, email string) string {
	homeDir, err := os.MkdirTemp(os.TempDir(), "cloudx-*")
	require.NoError(t, err)
	path := filepath.Join(homeDir, "import.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
  "schema_id": "preset://username",
  "traits": {
    "username": "`+email+`"
  }
}`), 0600))
	return path
}

func TestImportIdentity(t *testing.T) {
	configDir := newConfigDir(t)
	cmd := configAwareCmd(configDir)

	email, password := registerAccount(t, configDir)
	project := createProject(t, configDir)

	t.Run("is not able to import identities if not authenticated and quiet flag", func(t *testing.T) {
		configDir := newConfigDir(t)
		cmd := configAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "import", "identities", "--quiet", "--project", project)
		require.ErrorIs(t, err, ErrNoConfigQuiet)
	})

	success := func(t *testing.T, cmd *cmdx.CommandExecuter, stdin *bytes.Buffer) {
		email := fakeEmail()
		stdout, stderr, err := cmd.Exec(stdin, "import", "identities", "--format", "json", "--project", project, makeRandomIdentity(t, email))
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, email, out.Get("traits.username").String())
	}

	t.Run("is able to import identities", func(t *testing.T) {
		success(t, cmd, nil)
	})

	t.Run("is able to import identities after authenticating", func(t *testing.T) {
		cmd, r := withReAuth(t, email, password)
		success(t, cmd, r)
	})
}
