// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestDeleteIdentity(t *testing.T) {
	t.Run("is not able to delete identities if not authenticated and quiet flag", func(t *testing.T) {
		userID := testhelpers.ImportIdentity(t, defaultCmd, defaultProject, nil)
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "delete", "identity", "--quiet", "--project", defaultProject, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to delete identities", func(t *testing.T) {
		userID := testhelpers.ImportIdentity(t, defaultCmd, defaultProject, nil)
		stdout, stderr, err := defaultCmd.Exec(nil, "delete", "identity", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, userID, out.String(), "stdout: %s", stdout)
	})

	t.Run("is able to delete identities after authenticating", func(t *testing.T) {
		userID := testhelpers.ImportIdentity(t, defaultCmd, defaultProject, nil)
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "delete", "identity", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Equal(t, userID, out.String(), stdout)
	})
}
