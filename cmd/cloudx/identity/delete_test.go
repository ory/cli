// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestDeleteIdentity(t *testing.T) {
	t.Parallel()

	userID := testhelpers.ImportIdentity(ctx, t, defaultProject.Id, nil)

	t.Run("is not able to delete identities if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "delete", "identity", "--quiet", "--project", defaultProject.Id, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "delete", "identity", "--project", defaultProject.Id, userID)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	t.Run("is able to delete identities", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "delete", "identity", "--format", "json", "--project", defaultProject.Id, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, userID, out.String(), "stdout: %s", stdout)
	})
}
