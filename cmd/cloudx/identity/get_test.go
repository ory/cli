// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"context"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetIdentity(t *testing.T) {
	t.Parallel()

	userID := testhelpers.ImportIdentity(ctx, t, defaultProject.Id, nil)

	t.Run("is not able to get identity if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "get", "identity", "--quiet", "--project", defaultProject.Id, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "get", "identity", "--project", defaultProject.Id, userID)
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	t.Run("is able to get identity", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "get", "identity", "--format", "json", "--project", defaultProject.Id, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("id").String())
	})
}
