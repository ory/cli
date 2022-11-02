// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetIdentity(t *testing.T) {
	userID := testhelpers.ImportIdentity(t, defaultCmd, defaultProject, nil)

	t.Run("is not able to get identity if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "get", "identity", "--quiet", "--project", defaultProject, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to get identity", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "get", "identity", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("id").String())
	})

	t.Run("is able to get identity after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "get", "identity", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("id").String())
	})
}
