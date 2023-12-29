// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestListIdentities(t *testing.T) {
	project := testhelpers.CreateProject(t, defaultConfig)

	userID := testhelpers.ImportIdentity(t, defaultCmd, project, nil)

	t.Run("is not able to list identities if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "list", "identities", "--quiet", "--project", project, "--page-size=5")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s identities", proc), func(t *testing.T) {
			stdout, stderr, err := defaultCmd.Exec(nil, proc, "identities", "--format", "json", "--project", project, "--page-size=5")
			require.NoError(t, err, stderr)
			out := gjson.Parse(stdout)
			assert.True(t, gjson.Valid(stdout))
			assert.Len(t, out.Array(), 1)
			assert.Equal(t, userID, out.Array()[0].Get("id").String())
		})
	}

	t.Run("is able to list identities after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "ls", "identities", "--format", "json", "--project", project, "--page-size=5")
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("id").String())
	})
}
