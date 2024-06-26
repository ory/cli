// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestUseProject(t *testing.T) {
	t.Run("is able to use project", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject.Id)

		stdout, _, err := defaultCmd.Exec(nil, "use", "project", extraProject.Id, "--quiet")
		require.NoError(t, err)
		assert.Equal(t, extraProject.Id, strings.TrimSpace(stdout))
	})
	t.Run("is able to print default project", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject.Id)

		stdout, _, err := defaultCmd.Exec(nil, "use", "project", "--quiet")
		require.NoError(t, err)
		assert.Equal(t, defaultProject.Id, strings.TrimSpace(stdout))
	})
}
