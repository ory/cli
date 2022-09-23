package oauth2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/ory/kratos/x"

	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
)

func TestImportIdentity(t *testing.T) {
	t.Run("is not able to import oauth2-client if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "import", "oauth2-client", "--quiet", "--project", defaultProject)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to import oauth2-client", func(t *testing.T) {
		name := x.NewUUID().String()
		stdout, stderr, err := defaultCmd.Exec(nil, "import", "oauth2-client", "--format", "json", "--project", defaultProject, testhelpers.MakeRandomClient(t, name))
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, name, out.Get("client_name").String())
	})

	t.Run("is able to import oauth2-client after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		name := x.NewUUID().String()
		stdout, stderr, err := cmd.Exec(r, "import", "oauth2-client", "--format", "json", "--project", defaultProject, testhelpers.MakeRandomClient(t, name))
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, name, out.Get("client_name").String())
	})
}
