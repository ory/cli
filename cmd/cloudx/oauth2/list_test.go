package oauth2_test

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestListIdentities(t *testing.T) {
	project := testhelpers.CreateProject(t, defaultConfig)

	userID := testhelpers.CreateClient(t, defaultCmd, project).Get("client_id").String()

	t.Run("is not able to list oauth2 clients if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "list", "oauth2-clients", "--quiet", "--project", project)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s oauth2 clients", proc), func(t *testing.T) {
			stdout, stderr, err := defaultCmd.Exec(nil, proc, "oauth2-clients", "--format", "json", "--project", project)
			require.NoError(t, err, stderr)
			out := gjson.Parse(stdout).Get("items")
			assert.True(t, gjson.Valid(stdout))
			assert.Len(t, out.Array(), 1)
			assert.Equal(t, userID, out.Array()[0].Get("client_id").String(), "%s", out)
		})
	}

	t.Run("is able to list oauth2 clients after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "ls", "oauth2-clients", "--format", "json", "--project", project)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout).Get("items")
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String(), "%s", out)
	})
}
