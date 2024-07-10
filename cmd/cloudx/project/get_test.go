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
	t.Parallel()

	getProject := func(projectID string) func(t *testing.T, exec execFunc) {
		return func(t *testing.T, exec execFunc) {
			stdout, _, err := exec(nil, "get", "project", "--format", "json")
			require.NoError(t, err)
			assert.Equal(t, projectID, gjson.Get(stdout, "id").String())
			assert.NotEmpty(t, gjson.Get(stdout, "slug").String())
		}
	}

	runWithProjectAsDefault(ctx, t, defaultProject.Id, getProject(defaultProject.Id))
	runWithProjectAsArgument(ctx, t, extraProject.Id, getProject(extraProject.Id))
}

func TestGetServiceConfig(t *testing.T) {
	t.Parallel()

	t.Run("service=identity", func(t *testing.T) {
		t.Parallel()

		getIdentityConfig := func(t *testing.T, exec execFunc) {
			stdout, _, err := exec(nil, "get", "identity-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "selfservice.flows.error.ui_url").Exists())
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, getIdentityConfig)
		runWithProjectAsFlag(ctx, t, extraProject.Id, getIdentityConfig)
	})

	t.Run("service=permissions", func(t *testing.T) {
		t.Parallel()

		getPermissionsConfig := func(t *testing.T, exec execFunc) {
			stdout, _, err := exec(nil, "get", "permission-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "namespaces").Exists(), stdout)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, getPermissionsConfig)
		runWithProjectAsFlag(ctx, t, extraProject.Id, getPermissionsConfig)
	})

	t.Run("service=oauth2", func(t *testing.T) {
		t.Parallel()

		getOAuth2Config := func(t *testing.T, exec execFunc) {
			stdout, _, err := exec(nil, "get", "oauth2-config", "--format", "json")
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "oauth2").Exists(), stdout)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, getOAuth2Config)
		runWithProjectAsFlag(ctx, t, extraProject.Id, getOAuth2Config)
	})
}
