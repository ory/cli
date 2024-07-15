package client_test

import (
	"bytes"
	"context"
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
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
