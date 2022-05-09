package project_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/ghodss/yaml"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetKratosConfig(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	cmd := testhelpers.ConfigAwareCmd(configDir)
	email, password := testhelpers.RegisterAccount(t, configDir)

	project := testhelpers.CreateProject(t, configDir)

	t.Run(fmt.Sprintf("is able to get project"), func(t *testing.T) {
		stdout, _, err := cmd.Exec(nil, "get", "kratos-config", project, "--format", "json")
		require.NoError(t, err)
		actual, err := yaml.YAMLToJSON([]byte(stdout))
		require.NoError(t, err)
		assert.Equal(t, "/ui/error", gjson.GetBytes(actual, "selfservice.flows.error.ui_url").String())
	})

	t.Run("is not able to list projects if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "get", "identity-config", project, "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to get project after authenticating", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigPasswordAwareCmd(configDir, password)
		// Create the account
		var r bytes.Buffer
		r.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
		r.WriteString(email + "\n") // Email fakeEmail()
		stdout, _, err := cmd.Exec(&r, "get", "ic", project, "--format", "json")
		require.NoError(t, err)
		actual, err := yaml.YAMLToJSON([]byte(stdout))
		require.NoError(t, err)
		assert.Equal(t, "/ui/error", gjson.GetBytes(actual, "selfservice.flows.error.ui_url").String())
	})
}
