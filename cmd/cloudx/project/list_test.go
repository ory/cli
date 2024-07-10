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

	// this test needs a separate account to properly list projects
	_, _, _, sessionToken := testhelpers.RegisterAccount(context.Background(), t)
	ctx := client.ContextWithOptions(ctx,
		client.WithSessionToken(t, sessionToken),
		client.WithConfigLocation(testhelpers.NewConfigFile(t)))
	cmd := testhelpers.Cmd(ctx)

	projects := make([]*cloud.Project, 3)
	projectIDs := make([]string, len(projects))
	for k := range projects {
		projects[k] = testhelpers.CreateProject(ctx, t, nil)
		projectIDs[k] = projects[k].Id
	}

	assertHasProjects := func(t *testing.T, stdout string) {
		out := gjson.Parse(stdout)
		assert.EqualValues(t, out.Get("#").Int(), len(projects))
		out.ForEach(func(_, project gjson.Result) bool {
			assert.Contains(t, projectIDs, project.Get("id").String())
			return true
		})
	}

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s projects", proc), func(t *testing.T) {
			t.Parallel()

			stdout, _, err := cmd.Exec(nil, proc, "projects", "--format", "json")
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
