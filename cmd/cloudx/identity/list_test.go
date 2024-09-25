// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestListIdentities(t *testing.T) {
	t.Parallel()

	workspace := testhelpers.CreateWorkspace(ctx, t)
	project := testhelpers.CreateProject(ctx, t, workspace)
	userID := testhelpers.ImportIdentity(ctx, t, project.Id, nil)

	t.Run("is not able to list identities if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "list", "identities", "--quiet", "--project", project.Id, "--consistency", "strong")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "list", "identities", "--project", project.Id, "--consistency", "strong")
		require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s identities", proc), func(t *testing.T) {
			stdout, stderr, err := defaultCmd.Exec(nil, proc, "identities", "--format", "json", "--project", project.Id, "--consistency", "strong")
			require.NoError(t, err, stderr)
			out := gjson.Parse(stdout)
			assert.True(t, gjson.Valid(stdout))
			assert.Len(t, out.Get("identities").Array(), 1)
			assert.Equal(t, userID, out.Get("identities").Array()[0].Get("id").String(), out.Raw)
		})
	}
}
