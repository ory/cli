package cloudx

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestListIdentities(t *testing.T) {
	configDir := newConfigDir(t)

	cmd := configAwareCmd(configDir)

	email, password := registerAccount(t, configDir)
	project := createProject(t, configDir)

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s identities", proc), func(t *testing.T) {
			stdout, _, err := cmd.Exec(nil, proc, "identities", "--format", "json", "--project", project)
			require.NoError(t, err)
			out := gjson.Parse(stdout)
			assert.Len(t, out.Array(), 2)
		})
	}

	t.Run("is not able to list identities if not authenticated and quiet flag", func(t *testing.T) {
		configDir := newConfigDir(t)
		cmd := configAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "list", "identities", "--quiet", "--project", project)
		require.ErrorIs(t, err, ErrNoConfigQuiet)
	})

	t.Run("is able to list identities after authenticating", func(t *testing.T) {
		configDir := newConfigDir(t)
		cmd := configPasswordAwareCmd(configDir, password)
		// Create the account
		var r bytes.Buffer
		r.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
		r.WriteString(email + "\n") // Email fakeEmail()
		stdout, _, err := cmd.Exec(&r, "ls", "identities", "--format", "json", "--project", project)
		require.NoError(t, err)

		for _, project := range gjson.Parse(stdout).Array() {
			assert.Contains(t, "projects", project.Get("id").String())
		}
	})
}
