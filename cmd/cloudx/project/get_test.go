// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetProject(t *testing.T) {
	testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
	t.Run("is able to get project", func(t *testing.T) {
		assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

		stdout, _, err := defaultCmd.Exec(nil, "get", "project", extraProject, "--format", "json")
		require.NoError(t, err)
		assert.Equal(t, extraProject, gjson.Get(stdout, "id").String())
		assert.NotEmpty(t, gjson.Get(stdout, "slug").String())
	})
	t.Run("is able to get default project", func(t *testing.T) {
		assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

		stdout, _, err := defaultCmd.Exec(nil, "get", "project", "--format", "json")
		require.NoError(t, err)
		assert.Equal(t, defaultProject, gjson.Get(stdout, "id").String())
		assert.NotEmpty(t, gjson.Get(stdout, "slug").String())
	})
}

func TestGetServiceConfig(t *testing.T) {
	t.Run("service=kratos", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "get", "kratos-config", defaultProject, "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "selfservice.flows.error.ui_url").Exists())
		})
		t.Run("default project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "get", "kratos-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "selfservice.flows.error.ui_url").Exists())
		})
	})

	t.Run("service=keto", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "get", "keto-config", defaultProject, "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "namespaces").Exists(), stdout)
		})
		t.Run("default project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "get", "keto-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "namespaces").Exists(), stdout)
		})
	})

	t.Run("service=hydra", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "get", "oauth2-config", defaultProject, "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "oauth2").Exists(), stdout)
		})
		t.Run("default project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "get", "oauth2-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "oauth2").Exists(), stdout)
		})
	})
}
