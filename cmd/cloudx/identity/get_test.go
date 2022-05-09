package identity_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetIdentity(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	cmd := testhelpers.ConfigAwareCmd(configDir)

	email, password := testhelpers.RegisterAccount(t, configDir)
	project := testhelpers.CreateProject(t, configDir)

	userID := testhelpers.ImportIdentity(t, cmd, project, nil)

	t.Run("is not able to get identities if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "get", "identity", "--quiet", "--project", project, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
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
		cmd, r := testhelpers.WithReAuth(t, email, password)
		stdout, stderr, err := cmd.Exec(r, "get", "identity", "--format", "json", "--project", project, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("id").String())
	})
}
