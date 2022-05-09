package project_test

import (
	"bytes"
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

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s projects", proc), func(t *testing.T) {
			stdout, _, err := cmd.Exec(nil, proc, "projects", "--format", "json")
			require.NoError(t, err)
			out := gjson.Parse(stdout)
			assert.Len(t, out.Array(), 3)
			for _, project := range out.Array() {
				assert.Contains(t, projects, project.Get("id").String())
			}
		})
	}

	t.Run("is not able to list projects if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "list", "projects", "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to list projects after authenticating", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigPasswordAwareCmd(configDir, password)
		// Create the account
		var r bytes.Buffer
		r.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
		r.WriteString(email + "\n") // Email fakeEmail()
		stdout, _, err := cmd.Exec(&r, "ls", "projects", "--format", "json")
		require.NoError(t, err)
		for _, project := range gjson.Parse(stdout).Array() {
			assert.Contains(t, projects, project.Get("id").String())
		}
	})
}
