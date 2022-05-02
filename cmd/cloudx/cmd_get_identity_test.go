package cloudx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetIdentity(t *testing.T) {
	configDir := newConfigDir(t)
	cmd := configAwareCmd(configDir)

	email, password := registerAccount(t, configDir)
	project := createProject(t, configDir)

	userID := importIdentity(t, cmd, project, nil)

	t.Run("is not able to get identities if not authenticated and quiet flag", func(t *testing.T) {
		configDir := newConfigDir(t)
		cmd := configAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "get", "identity", "--quiet", "--project", project, userID)
		require.ErrorIs(t, err, ErrNoConfigQuiet)
	})

	t.Run("is able to get identities", func(t *testing.T) {
		stdout, stderr, err := cmd.Exec(nil, "get", "identity", "--format", "json", "--project", project, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("id").String())
	})

	t.Run("is able to get identities after authenticating", func(t *testing.T) {
		cmd, r := withReAuth(t, email, password)
		stdout, stderr, err := cmd.Exec(r, "get", "identity", "--format", "json", "--project", project, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("id").String())
	})
}
