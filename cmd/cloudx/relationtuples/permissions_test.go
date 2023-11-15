// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package relationtuples_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestIsAllowedNoUnauthenticated(t *testing.T) {
	t.Parallel()

	configDir := testhelpers.NewConfigDir(t)
	cmd := testhelpers.ConfigAwareCmd(configDir)
	_, _, err := cmd.Exec(nil,
		"is", "allowed", "user", "relation", "namespace", "object",
		"--quiet", "--project", defaultProject)
	require.ErrorIsf(t, err, client.ErrNoConfigQuiet, "got error: %v", err)
}
