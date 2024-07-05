// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"context"
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
)

func TestImportIdentity(t *testing.T) {
	t.Parallel()

	t.Run("is not able to import identities if not authenticated and quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "import", "identities", "--quiet", "--project", defaultProject.Id)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "import", "identities", "--quiet", "--project", defaultProject.Id)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to import identities", func(t *testing.T) {
		testhelpers.ImportIdentity(ctx, t, defaultProject.Id, nil)
	})
}
