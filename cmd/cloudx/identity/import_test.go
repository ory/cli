package identity_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
)

func TestImportIdentity(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	cmd := testhelpers.ConfigAwareCmd(configDir)

	email, password := testhelpers.RegisterAccount(t, configDir)
	project := testhelpers.CreateProject(t, configDir)

	t.Run("is not able to import identities if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "import", "identities", "--quiet", "--project", project)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to import identities", func(t *testing.T) {
		testhelpers.ImportIdentity(t, cmd, project, nil)
	})

	t.Run("is able to import identities after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, email, password)
		testhelpers.ImportIdentity(t, cmd, project, r)
	})
}
