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

func TestPatchHydraConfig(t *testing.T) {
	t.Run("is able to replace a key", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "hydra-config", extraProject, "--format", "json", "--replace", `/strategies/access_token="jwt"`)
			require.NoError(t, err)
			assert.Equal(t, "jwt", gjson.Get(stdout, "strategies.access_token").String())
		})
		t.Run("default project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "hydra-config", "--format", "json", "--replace", `/strategies/access_token="jwt"`)
			require.NoError(t, err)
			assert.Equal(t, "jwt", gjson.Get(stdout, "strategies.access_token").String())
		})
	})

	t.Run("is able to add a key", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "oauth2-config", extraProject, "--format", "json", "--add", `/ttl/login_consent_request="1h"`)
			require.NoError(t, err)
			assert.Equal(t, "1h0m0s", gjson.Get(stdout, "ttl.login_consent_request").String())
		})
		t.Run("default project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "oauth2-config", "--format", "json", "--add", `/ttl/login_consent_request="1h"`)
			require.NoError(t, err)
			assert.Equal(t, "1h0m0s", gjson.Get(stdout, "ttl.login_consent_request").String())
		})
	})

	t.Run("is able to add a key with string", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "oc", extraProject, "--format", "json", "--replace", `/ttl/refresh_token="2h"`)
			require.NoError(t, err)
			assert.Equal(t, "2h0m0s", gjson.Get(stdout, "ttl.refresh_token").String())
		})
		t.Run("default project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "oc", "--format", "json", "--replace", `/ttl/refresh_token="2h"`)
			require.NoError(t, err)
			assert.Equal(t, "2h0m0s", gjson.Get(stdout, "ttl.refresh_token").String())
		})
	})

	t.Run("fails if no opts are given", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)
		t.Run("explicit project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "oc", extraProject, "--format", "json")
			require.Error(t, err, stdout)
		})
		t.Run("default project", func(t *testing.T) {
			stdout, _, err := defaultCmd.Exec(nil, "patch", "oc", "--format", "json")
			require.Error(t, err, stdout)
		})
	})
}
