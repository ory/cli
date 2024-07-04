// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
)

func TestImportIdentity(t *testing.T) {
	t.Run("is not able to import identities if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigFile(t)
		cmd := testhelpers.Cmd(configDir)
		_, _, err := cmd.Exec(nil, "import", "identities", "--quiet", "--project", defaultProject.Id)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to import identities", func(t *testing.T) {
		testhelpers.ImportIdentity(t, defaultCmd, defaultProject.Id, nil)
	})

	t.Run("is able to import identities after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		testhelpers.ImportIdentity(t, cmd, defaultProject.Id, r)
	})
}
