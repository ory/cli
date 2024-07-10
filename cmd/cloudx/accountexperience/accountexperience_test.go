// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestMain(m *testing.M) {
	testhelpers.UseStaging()
	m.Run()
}

func TestOpenAXPages(t *testing.T) {
	_, _, _, sessionToken := testhelpers.RegisterAccount(context.Background(), t)
	ctx := client.ContextWithOptions(context.Background(),
		client.WithConfigLocation(testhelpers.NewConfigFile(t)),
		client.WithSessionToken(t, sessionToken))
	project := testhelpers.CreateProject(ctx, t, nil)
	cmd := testhelpers.Cmd(ctx)

	t.Run("is able to open all pages", func(t *testing.T) {
		for _, flowType := range []string{"login", "registration", "recovery", "verification", "settings"} {
			testhelpers.Cmd(client.ContextWithOptions(ctx, client.WithOpenBrowserHook(func(uri string) error {
				assert.Truef(t, strings.HasPrefix(uri, "https://"+project.Slug), "expected %q to have prefix %q", uri, "https://"+project.Slug)
				assert.Contains(t, uri, flowType)
				return nil
			})))

			stdout, stderr, err := cmd.Exec(nil, "open", "account-experience", flowType, "--quiet")
			require.NoError(t, err, stderr)
			assert.Contains(t, stdout, "https://"+project.Slug)
			assert.Contains(t, stdout, flowType)
		}
	})

	t.Run("errors on unknown flow type", func(t *testing.T) {
		stdout, stderr, err := cmd.Exec(nil, "open", "account-experience", "unknown", "--quiet")
		require.Error(t, err)
		assert.Contains(t, stderr, "unknown flow type", stdout)
	})
}
