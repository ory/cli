// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofrs/uuid/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

var (
	ctx            context.Context
	defaultProject *cloud.Project
	defaultCmd     *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	ctx, _, _, _, defaultProject, defaultCmd = testhelpers.CreateDefaultAssetsBrowser()
	m.Run()
}

func TestCreateClient(t *testing.T) {
	t.Parallel()

	t.Run("is not able to create client if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "client", "--quiet", "--project", defaultProject.Id)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "client", "--project", defaultProject.Id)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	t.Run("is able to create client", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "create", "client", "--format", "json", "--project", defaultProject.Id)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		t.Logf("Created client: %s", stdout)
	})
}

func TestDeleteClient(t *testing.T) {
	t.Parallel()

	t.Run("is not able to delete oauth2 client if not authenticated and quiet flag", func(t *testing.T) {
		userID := testhelpers.CreateClient(ctx, t, defaultProject.Id).Get("client_id").String()

		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "delete", "oauth2-client", "--quiet", "--project", defaultProject.Id, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow if not authenticated", func(t *testing.T) {
		userID := testhelpers.CreateClient(ctx, t, defaultProject.Id).Get("client_id").String()

		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "delete", "oauth2-client", "--project", defaultProject.Id, userID)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	t.Run("is able to delete oauth2 client", func(t *testing.T) {
		clientID := testhelpers.CreateClient(ctx, t, defaultProject.Id).Get("client_id").String()
		stdout, stderr, err := defaultCmd.Exec(nil, "delete", "oauth2-client", "--format", "json", "--project", defaultProject.Id, clientID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, clientID, out.String(), "stdout: %s", stdout)
	})
}

func TestGetClient(t *testing.T) {
	t.Parallel()

	clientID := testhelpers.CreateClient(ctx, t, defaultProject.Id).Get("client_id").String()

	t.Run("is not able to get oauth2 if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "get", "oauth2-client", "--quiet", "--project", defaultProject.Id, clientID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "get", "oauth2-client", "--project", defaultProject.Id, clientID)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	t.Run("is able to get oauth2", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "get", "oauth2-client", "--format", "json", "--project", defaultProject.Id, clientID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, clientID, out.Array()[0].Get("client_id").String())
	})
}

func TestImportClient(t *testing.T) {
	t.Parallel()

	t.Run("is not able to import oauth2-client if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "import", "oauth2-client", "--quiet", "--project", defaultProject.Id)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "import", "oauth2-client", "--project", defaultProject.Id)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	t.Run("is able to import oauth2-client", func(t *testing.T) {
		name := uuid.Must(uuid.NewV4()).String()
		stdout, stderr, err := defaultCmd.Exec(nil, "import", "oauth2-client", "--format", "json", "--project", defaultProject.Id, testhelpers.MakeRandomClient(t, name))
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, name, out.Get("client_name").String())
	})
}

func TestListClients(t *testing.T) {
	t.Parallel()

	workspace := testhelpers.CreateWorkspace(ctx, t)
	project := testhelpers.CreateProject(ctx, t, workspace)
	clientID := testhelpers.CreateClient(ctx, t, project.Id).Get("client_id").String()

	t.Run("is not able to list oauth2 clients if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "list", "oauth2-clients", "--quiet", "--project", project.Id)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "list", "oauth2-clients", "--project", project.Id)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s oauth2 clients", proc), func(t *testing.T) {
			stdout, stderr, err := defaultCmd.Exec(nil, proc, "oauth2-clients", "--format", "json", "--project", project.Id)
			require.NoError(t, err, stderr)
			out := gjson.Parse(stdout).Get("items")
			assert.True(t, gjson.Valid(stdout))
			assert.Len(t, out.Array(), 1)
			assert.Equal(t, clientID, out.Get("0.client_id").String(), "%s", out)
		})
	}
}

func TestUpdateOAuth2(t *testing.T) {
	t.Parallel()

	clientID := testhelpers.CreateClient(ctx, t, defaultProject.Id).Get("client_id").String()

	t.Run("is not able to update oauth2 if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "update", "oauth2-client", "--quiet", "--project", defaultProject.Id, clientID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "update", "oauth2-client", "--project", defaultProject.Id, clientID)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	t.Run("is able to update oauth2", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "update", "oauth2-client", "--format", "json", "--project", defaultProject.Id, clientID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, clientID, out.Array()[0].Get("client_id").String())
	})
}
