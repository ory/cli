package cloudx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestDeleteIdentity(t *testing.T) {
	configDir := newConfigDir(t)
	cmd := configAwareCmd(configDir)

	email, password := registerAccount(t, configDir)
	project := createProject(t, configDir)

	t.Run("is not able to delete identities if not authenticated and quiet flag", func(t *testing.T) {
		userID := importIdentity(t, cmd, project, nil)
		configDir := newConfigDir(t)
		cmd := configAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "delete", "identity", "--quiet", "--project", project, userID)
		require.ErrorIs(t, err, ErrNoConfigQuiet)
	})

	t.Run("is able to delete identities", func(t *testing.T) {
		userID := importIdentity(t, cmd, project, nil)
		stdout, stderr, err := cmd.Exec(nil, "delete", "identity", "--format", "json", "--project", project, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, userID, out.String(), "stdout: %s", stdout)
	})

	t.Run("is able to delete identities after authenticating", func(t *testing.T) {
		userID := importIdentity(t, cmd, project, nil)
		cmd, r := withReAuth(t, email, password)
		stdout, stderr, err := cmd.Exec(r, "delete", "identity", "--format", "json", "--project", project, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Equal(t, userID, out.String(), stdout)
	})
}
