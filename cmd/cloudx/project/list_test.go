// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
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
	configDir := testhelpers.NewConfigFile(t)
	cmd := testhelpers.CmdWithConfig(configDir)
	email, password, _ := testhelpers.RegisterAccount(t, configDir)

	projects := make([]*cloud.Project, 3)
	projectIDs := make([]string, len(projects))
	for k := range projects {
		projects[k] = testhelpers.CreateProject(t, configDir, nil)
		projectIDs[k] = projects[k].Id
	}
	t.Logf("Created projects %+v", projects)

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
			stdout, _, err := cmd.Exec(nil, proc, "projects", "--format", "json")
			require.NoError(t, err)
			assertHasProjects(t, stdout)
		})
	}

	t.Run("is not able to list projects if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigFile(t)
		cmd := testhelpers.CmdWithConfig(configDir)
		_, _, err := cmd.Exec(nil, "list", "projects", "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to list projects after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, email, password)
		stdout, _, err := cmd.Exec(r, "ls", "projects", "--format", "json")
		require.NoError(t, err)
		assertHasProjects(t, stdout)
	})
}
