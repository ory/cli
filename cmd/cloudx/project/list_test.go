// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestListProject(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	cmd := testhelpers.ConfigAwareCmd(configDir)
	email, password := testhelpers.RegisterAccount(t, configDir)

	projects := make([]string, 3)
	for k := range projects {
		projects[k] = testhelpers.CreateProject(t, configDir)
	}
	t.Logf("Creating projects %v", projects)

	assertHasProjects := func(t *testing.T, stdout string) {
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), len(projects))
		out.ForEach(func(_, project gjson.Result) bool {
			assert.Contains(t, projects, project.Get("id").String())
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
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
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
