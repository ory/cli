// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package relationtuples_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestIsAllowedNoUnauthenticated(t *testing.T) {
	t.Parallel()

	cmd := testhelpers.Cmd(testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t))

	// with quiet flag
	_, _, err := cmd.Exec(nil,
		"is", "allowed", "user", "relation", "namespace", "object",
		"--quiet", "--project", defaultProject.Id)
	require.ErrorIsf(t, err, client.ErrNoConfigQuiet, "got error: %v", err)

	// without quiet flag
	_, _, err = cmd.Exec(nil,
		"is", "allowed", "user", "relation", "namespace", "object",
		"--project", defaultProject.Id)
	require.ErrorIsf(t, err, testhelpers.ErrAuthFlowTriggered, "got error: %v", err)
}
