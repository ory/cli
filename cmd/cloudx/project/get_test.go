package project_test

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/ghodss/yaml"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetProject(t *testing.T) {
	t.Run(fmt.Sprintf("is able to get project"), func(t *testing.T) {
		stdout, _, err := defaultCmd.Exec(nil, "get", "project", defaultProject, "--format", "json")
		require.NoError(t, err)
		assert.Contains(t, defaultProject, gjson.Parse(stdout).Get("id").String())
		assert.NotEmpty(t, defaultProject, gjson.Parse(stdout).Get("slug").String())
	})

	t.Run(fmt.Sprintf("is able to get project"), func(t *testing.T) {
		stdout, _, err := defaultCmd.Exec(nil, "get", "project", defaultProject, "--format", "yaml")
		require.NoError(t, err)
		actual, err := yaml.YAMLToJSON([]byte(stdout))
		require.NoError(t, err)
		assert.Contains(t, defaultProject, gjson.ParseBytes(actual).Get("id").String())
	})

	t.Run("is not able to get project if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "get", "project", defaultProject, "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to get project after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, _, err := cmd.Exec(r, "get", "project", defaultProject, "--format", "json")
		require.NoError(t, err)
		assert.Contains(t, defaultProject, gjson.Parse(stdout).Get("id").String())
	})
}
