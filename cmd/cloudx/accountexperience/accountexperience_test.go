// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestMain(m *testing.M) {
	testhelpers.RunAgainstStaging(m)
}

func TestOpenAXPages(t *testing.T) {
	cfg := testhelpers.NewConfigFile(t)
	testhelpers.RegisterAccount(t, cfg)
	project := testhelpers.CreateProject(t, cfg, nil)
	cmd := testhelpers.CmdWithConfig(cfg)

	t.Run("is able to open all pages", func(t *testing.T) {
		for _, flowType := range []string{"login", "registration", "recovery", "verification", "settings"} {
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
