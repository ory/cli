// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestAuthenticator(t *testing.T) {
	t.Parallel()

	t.Run("errors without config and --quiet flag", func(t *testing.T) {
		ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "auth", "--quiet")
		assert.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("triggers auth flow when not authenticated", func(t *testing.T) {
		ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
		_, _, err := testhelpers.Cmd(ctx).Exec(nil, "auth")
		assert.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	})

	// the full e2e flow is tested on the internal helper function instead of the full CLI wrapper
}
