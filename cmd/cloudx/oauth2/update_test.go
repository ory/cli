package oauth2_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestUpdateOAuth2(t *testing.T) {
	userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()

	t.Run("is not able to update oauth2 if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "update", "oauth2-client", "--quiet", "--project", defaultProject, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to update oauth2", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "update", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String())
	})

	t.Run("is able to update oauth2 after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "update", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String())
	})
}
