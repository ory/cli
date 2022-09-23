package oauth2_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestCreateClient(t *testing.T) {
	t.Run("is not able to create client if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "create", "client", "--quiet", "--project", defaultProject)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to create client", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "create", "client", "--format", "json", "--project", defaultProject)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		t.Logf("Created client: %s", stdout)
	})
}
