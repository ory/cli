// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestCreateProject(t *testing.T) {
	t.Parallel()

	parseOutput := func(stdout string) (id string, slug string, name string) {
		id = gjson.Get(stdout, "id").Str
		slug = gjson.Get(stdout, "slug").Str
		name = gjson.Get(stdout, "name").Str
		return
	}
	assertResult := func(t *testing.T, configDir string, stdout string, expectedName string) (id, slug, name string) {
		id, slug, name = parseOutput(stdout)
		assert.NotEmpty(t, id, stdout)
		assert.NotEmpty(t, slug, stdout)
		assert.Equal(t, expectedName, name, stdout)
		return
	}

	t.Run("requires workspace", func(t *testing.T) {
		t.Parallel()

		ctx, newConfig := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		conf := testhelpers.ReadConfig(t, newConfig)
		conf.SelectedWorkspace = uuid.Nil
		rawConfig, err := json.Marshal(conf)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(newConfig, rawConfig, 0600))

		for _, tc := range []struct {
			expectedErr string
			extraArgs   []string
		}{
			{
				expectedErr: "no workspace found",
				extraArgs:   nil,
			},
			{
				expectedErr: "The requested action was forbidden", // because the workspace does not exist we get a permission check error
				extraArgs:   []string{"--workspace", uuid.Must(uuid.NewV4()).String()},
			},
		} {
			t.Run("args="+strings.Join(tc.extraArgs, " "), func(t *testing.T) {
				_, stderr, err := testhelpers.Cmd(ctx).Exec(nil, append([]string{"create", "project", "--name", testhelpers.TestName(), "--format", "json", "--quiet"}, tc.extraArgs...)...)
				require.Error(t, err)
				assert.Contains(t, stderr, tc.expectedErr)
			})
		}
	})

	t.Run("is able to create a project with a workspace", func(t *testing.T) {
		t.Parallel()

		ctx, _ := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		testhelpers.SetDefaultProject(ctx, t, defaultProject.Id)

		name := testhelpers.TestName()
		pjRaw, stderr, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", name, "--create-workspace", "My new workspace", "--environment", "dev", "--format", "json")
		require.NoErrorf(t, err, "%s", stderr)
		assertResult(t, defaultConfig, pjRaw, name)

		wsID := gjson.Get(pjRaw, "workspace_id").Str
		require.NotZerof(t, wsID, "%s", pjRaw)

		wsRaw := testhelpers.Cmd(ctx).ExecNoErr(t, "get", "workspace", wsID, "--format", "json")
		assert.Equalf(t, wsID, gjson.Get(wsRaw, "id").Str, "%s", wsRaw)
		assert.Equalf(t, "My new workspace", gjson.Get(wsRaw, "name").Str, "%s", wsRaw)
	})

	t.Run("is able to create a project", func(t *testing.T) {
		t.Parallel()

		ctx, _ := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		testhelpers.SetDefaultProject(ctx, t, defaultProject.Id)

		wsID := testhelpers.CreateWorkspace(ctx, t)

		name := testhelpers.TestName()
		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", name, "--format", "json", "--workspace", wsID)
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)

		assert.Equal(t, defaultProject.Id, testhelpers.GetDefaultProjectID(ctx, t))
	})

	t.Run("is able to create a project and use the project as default", func(t *testing.T) {
		t.Parallel()

		ctx, _ := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)

		workspace := testhelpers.CreateWorkspace(ctx, t)
		name := testhelpers.TestName()
		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", name, "--use-project", "--workspace", workspace, "--environment", "dev", "--format", "json")
		require.NoError(t, err)
		id, _, _ := assertResult(t, defaultConfig, stdout, name)

		assert.Equal(t, id, testhelpers.GetDefaultProjectID(ctx, t))
	})

	t.Run("is able to create a project and use name from stdin", func(t *testing.T) {
		t.Parallel()

		name := testhelpers.TestName()
		stdin := strings.NewReader(name + "\n")
		workspace := testhelpers.CreateWorkspace(ctx, t)
		stdout, _, err := defaultCmd.Exec(stdin, "create", "project", "--workspace", workspace, "--environment", "dev", "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)
	})

	t.Run("is not able to create a project if no name flag and quiet flag", func(t *testing.T) {
		t.Parallel()

		workspace := testhelpers.CreateWorkspace(ctx, t)
		name := testhelpers.TestName()
		stdin := strings.NewReader(name)
		_, stderr, err := defaultCmd.Exec(stdin, "create", "project", "--quiet", "--workspace", workspace, "--environment", "dev")
		require.Error(t, err)
		assert.Contains(t, stderr, "you must specify the --name flag when using --quiet")
	})

	t.Run("is not able to create a project if not authenticated and quiet flag", func(t *testing.T) {
		t.Parallel()

		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)

		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", testhelpers.TestName(), "--workspace", uuid.Must(uuid.NewV4()).String(), "--environment", "dev", "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		t.Parallel()

		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)

		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", testhelpers.TestName(), "--workspace", uuid.Must(uuid.NewV4()).String(), "--environment", "dev", "--format", "json")
		assert.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})
}
