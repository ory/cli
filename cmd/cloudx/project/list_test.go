// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"context"
	"fmt"
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestListProject(t *testing.T) {
	t.Parallel()

	workspace := testhelpers.CreateWorkspace(ctx, t)
	projects := make([]*cloud.Project, 2)
	projectIDs := make([]string, len(projects))
	for k := range projects {
		projects[k] = testhelpers.CreateProject(ctx, t, workspace)
		projectIDs[k] = projects[k].Id
	}

	assertHasProjects := func(t *testing.T, stdout string) {
		out := gjson.Parse(stdout)
		require.EqualValues(t, out.Get("#").Int(), len(projects))
		actualIDs := make([]string, 0, len(projects))
		for _, id := range out.Get("#.id").Array() {
			actualIDs = append(actualIDs, id.Str)
		}
		assert.ElementsMatch(t, projectIDs, actualIDs)
	}

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s projects", proc), func(t *testing.T) {
			t.Parallel()

			stdout, _, err := defaultCmd.Exec(nil, proc, "projects", "--format", "json", "--workspace", workspace)
			require.NoError(t, err)
			assertHasProjects(t, stdout)
		})
	}

	t.Run("is not able to list projects if not authenticated and quiet flag", func(t *testing.T) {
		t.Parallel()

		cmd := testhelpers.Cmd(testhelpers.WithCleanConfigFile(context.Background(), t))
		_, _, err := cmd.Exec(nil, "list", "projects", "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers authentication flow", func(t *testing.T) {
		t.Parallel()

		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)

		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "ls", "projects")
		assert.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})
}
