// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetProject(t *testing.T) {
	runWithProject(t, func(t *testing.T, exec execFunc, projectID string) {
		stdout, _, err := exec(nil, "get", "project", "--format", "json")
		require.NoError(t, err)
		assert.Equal(t, projectID, gjson.Get(stdout, "id").String())
		assert.NotEmpty(t, gjson.Get(stdout, "slug").String())
	}, DefaultProject|PositionalProject)
}

func TestGetServiceConfig(t *testing.T) {
	t.Run("service=kratos", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "get", "kratos-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "selfservice.flows.error.ui_url").Exists())
		}, DefaultProject|PositionalProject)
	})

	t.Run("service=keto", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "get", "keto-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "namespaces").Exists(), stdout)
		}, DefaultProject|PositionalProject)
	})

	t.Run("service=hydra", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "get", "oauth2-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "oauth2").Exists(), stdout)
		}, DefaultProject|PositionalProject)
	})
}
