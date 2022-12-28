// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestUseProject(t *testing.T) {
	t.Run("is able to use project", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)

		stdout, _, err := defaultCmd.Exec(nil, "use", "project", extraProject, "--format", "json")
		require.NoError(t, err)
		assert.Equal(t, extraProject, gjson.Get(stdout, "id").String())
	})
	t.Run("is able to print default project", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)

		stdout, _, err := defaultCmd.Exec(nil, "use", "project", "--format", "json")
		require.NoError(t, err)
		assert.Equal(t, defaultProject, gjson.Get(stdout, "id").String())
	})
}
