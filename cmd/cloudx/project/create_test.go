// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestCreateProject(t *testing.T) {
	t.Parallel()

	parseOutput := func(stdout string) (id string, slug string, name string) {
		id = gjson.Get(stdout, "id").String()
		slug = gjson.Get(stdout, "slug").String()
		name = gjson.Get(stdout, "name").String()
		return
	}
	assertResult := func(t *testing.T, configDir string, stdout string, expectedName string) (id, slug, name string) {
		id, slug, name = parseOutput(stdout)
		assert.NotEmpty(t, id, stdout)
		assert.NotEmpty(t, slug, stdout)
		assert.Equal(t, expectedName, name, stdout)
		return
	}

	t.Run("is able to create a project", func(t *testing.T) {
		t.Parallel()

		ctx := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		testhelpers.SetDefaultProject(ctx, t, defaultProject.Id)

		name := testhelpers.TestName()
		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", name, "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)

		assert.Equal(t, defaultProject.Id, testhelpers.GetDefaultProjectID(ctx, t))
	})

	t.Run("is able to create a project and use the project as default", func(t *testing.T) {
		t.Parallel()

		ctx := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)

		name := testhelpers.TestName()
		stdout, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", name, "--use-project", "--format", "json")
		require.NoError(t, err)
		id, _, _ := assertResult(t, defaultConfig, stdout, name)

		assert.Equal(t, id, testhelpers.GetDefaultProjectID(ctx, t))
	})

	t.Run("is able to create a project and use name from stdin", func(t *testing.T) {
		t.Parallel()

		name := testhelpers.TestName()
		stdin := strings.NewReader(name + "\n")
		stdout, _, err := defaultCmd.Exec(stdin, "create", "project", "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)
	})

	t.Run("is not able to create a project if no name flag and quiet flag", func(t *testing.T) {
		t.Parallel()

		name := testhelpers.TestName()
		stdin := strings.NewReader(name)
		_, stderr, err := defaultCmd.Exec(stdin, "create", "project", "--quiet")
		require.Error(t, err)
		assert.Contains(t, stderr, "you must specify the --name and --environment flags when using --quiet")
	})

	t.Run("is not able to create a project if not authenticated and quiet flag", func(t *testing.T) {
		t.Parallel()

		ctx := client.ContextWithOptions(context.Background(), client.WithConfigLocation(testhelpers.NewConfigFile(t)))

		name := testhelpers.TestName()
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", name, "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		t.Parallel()

		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)

		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "create", "project", "--name", testhelpers.TestName(), "--format", "json")
		assert.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})
}
