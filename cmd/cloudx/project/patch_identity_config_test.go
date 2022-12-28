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

func TestPatchKratosConfig(t *testing.T) {
	t.Run("is able to replace a key", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "kratos-config", extraProject, "--format", "json", "--replace", `/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
		})
		t.Run("default project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "kratos-config", "--format", "json", "--replace", `/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
		})
	})

	t.Run("is able to add a key", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "identity-config", extraProject, "--format", "json", "--add", `/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
		})
		t.Run("default project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "identity-config", "--format", "json", "--add", `/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
		})
	})

	t.Run("is able to add a key with string", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "ic", extraProject, "--format", "json", "--replace", "/selfservice/flows/error/ui_url=\"https://example.com/error-ui\"")
			require.NoError(t, err)
			assert.Equal(t, "https://example.com/error-ui", gjson.Get(stdout, "selfservice.flows.error.ui_url").String())
		})
		t.Run("default project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "ic", "--format", "json", "--replace", "/selfservice/flows/error/ui_url=\"https://example.com/error-ui\"")
			require.NoError(t, err)
			assert.Equal(t, "https://example.com/error-ui", gjson.Get(stdout, "selfservice.flows.error.ui_url").String())
		})
	})

	t.Run("fails if no opts are given", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "ic", extraProject, "--format", "json")
			require.Error(t, err, stdout)
		})
		t.Run("default project", func(t *testing.T) {
			assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))

			stdout, _, err := defaultCmd.Exec(nil, "patch", "ic", "--format", "json")
			require.Error(t, err, stdout)
		})
	})
}
