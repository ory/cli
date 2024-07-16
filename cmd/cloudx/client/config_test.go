// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestLegacyConfigHandling(t *testing.T) {
	ctx := context.Background()
	legacyConfigFile := testhelpers.NewConfigFile(t)
	require.NoError(t, os.WriteFile(legacyConfigFile, []byte(`{"version": "v0alpha0"}`), 0600))

	out := bytes.Buffer{}
	h, err := client.NewCommandHelper(
		ctx,
		client.WithConfigLocation(legacyConfigFile),
		client.WithOpenBrowserHook(func(string) error {
			return testhelpers.ErrAuthFlowTriggered
		}),
		client.WithVerboseErrWriter(&out),
		client.WithStdin(strings.NewReader("\n")),
	)
	require.NoError(t, err)

	_, err = h.GetAuthenticatedConfig(ctx)
	assert.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered)
	assert.Contains(t, out.String(), "Thanks for upgrading!")
}
