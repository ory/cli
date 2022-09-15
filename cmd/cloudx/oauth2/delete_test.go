package oauth2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestDeleteOAuth2(t *testing.T) {
	t.Run("is not able to delete oauth2 client if not authenticated and quiet flag", func(t *testing.T) {
		userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "delete", "oauth2-client", "--quiet", "--project", defaultProject, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to delete oauth2 client", func(t *testing.T) {
		userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()
		stdout, stderr, err := defaultCmd.Exec(nil, "delete", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, userID, out.String(), "stdout: %s", stdout)
	})

	t.Run("is able to delete oauth2 client after authenticating", func(t *testing.T) {
		userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "delete", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Equal(t, userID, out.String(), stdout)
	})
}
