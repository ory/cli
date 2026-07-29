// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package workspace_test

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

// TestUseWorkspace covers https://github.com/ory/cli/issues/406. The command is
// exercised against a config file only, so it needs no Ory Network access:
// selecting an already-known workspace ID never hits the API.
func TestUseWorkspace(t *testing.T) {
	t.Parallel()

	const (
		initial     = "11111111-1111-1111-1111-111111111111"
		other       = "33333333-3333-3333-3333-333333333333"
		project     = "22222222-2222-2222-2222-222222222222"
		accessToken = "the-access-token"
	)

	newContext := func(t *testing.T, workspace string) (context.Context, string) {
		t.Helper()

		conf := map[string]any{
			"version":          client.ConfigVersion,
			"selected_project": project,
			"access_token":     map[string]any{"access_token": accessToken, "token_type": "bearer"},
		}
		if workspace != "" {
			conf["selected_workspace"] = workspace
		}
		raw, err := json.Marshal(conf)
		require.NoError(t, err)

		location := testhelpers.NewConfigFile(t)
		require.NoError(t, os.WriteFile(location, raw, 0600))

		return client.ContextWithOptions(context.Background(), client.WithConfigLocation(location)), location
	}

	t.Run("case=prints the default workspace when no id is given", func(t *testing.T) {
		t.Parallel()

		ctx, _ := newContext(t, initial)

		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "use", "workspace", "--quiet")
		require.NoError(t, err)
		assert.Equal(t, initial, strings.TrimSpace(stdout))
	})

	t.Run("case=prints the default workspace as JSON", func(t *testing.T) {
		t.Parallel()

		ctx, _ := newContext(t, initial)

		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "use", "workspace", "--format", "json")
		require.NoError(t, err)
		assert.JSONEq(t, `{"id":"`+initial+`"}`, stdout)
	})

	t.Run("case=does not persist an overridden workspace when no id is given", func(t *testing.T) {
		t.Parallel()

		ctx, location := newContext(t, initial)
		// Simulates ORY_WORKSPACE (or a workspace API key) pointing somewhere
		// else: printing the default must not rewrite it.
		ctx = client.ContextWithOptions(ctx, client.WithWorkspaceOverride(other))

		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "use", "workspace", "--quiet")
		require.NoError(t, err)
		assert.Equal(t, other, strings.TrimSpace(stdout))

		conf := testhelpers.ReadConfig(t, location)
		assert.Equal(t, initial, conf.SelectedWorkspace.String())
	})

	t.Run("case=sets the default workspace and persists it", func(t *testing.T) {
		t.Parallel()

		ctx, location := newContext(t, initial)

		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "use", "workspace", other, "--quiet")
		require.NoError(t, err)
		assert.Equal(t, other, strings.TrimSpace(stdout))

		conf := testhelpers.ReadConfig(t, location)
		assert.Equal(t, other, conf.SelectedWorkspace.String())
		assert.Equal(t, project, conf.SelectedProject.String(), "selecting a workspace must not clear the project")
		require.NotNil(t, conf.AccessToken, "selecting a workspace must not log the user out")
		assert.Equal(t, accessToken, conf.AccessToken.AccessToken)
	})

	t.Run("case=errors when no workspace is set", func(t *testing.T) {
		t.Parallel()

		ctx, _ := newContext(t, "")

		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "use", "workspace", "--quiet")
		assert.ErrorIs(t, err, client.ErrWorkspaceNotSet)
	})
}
