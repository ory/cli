// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestPatchProject(t *testing.T) {
	t.Run("is able to replace a key", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "project", "--format", "json", "--replace", `/services/identity/config/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.password.enabled").Bool())
		}, WithDefaultProject, WithPositionalProject)
	})

	t.Run("is able to add a key", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "project", "--format", "json", "--add", `/services/identity/config/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.password.enabled").Bool())
		}, WithDefaultProject, WithPositionalProject)
	})

	t.Run("is able to add a key with string", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "project", "--format", "json", "--replace", "/services/identity/config/selfservice/flows/error/ui_url=\"https://example.com/error-ui\"")
			require.NoError(t, err)
			assert.Equal(t, "https://example.com/error-ui", gjson.Get(stdout, "services.identity.config.selfservice.flows.error.ui_url").String())
		}, WithDefaultProject, WithPositionalProject)
	})

	t.Run("is able to add a key with raw json", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "project", "--format", "json", "--replace", `/services/identity/config/selfservice/flows/error={"ui_url":"https://example.org/error-ui"}`)
			require.NoError(t, err)
			assert.Equal(t, "https://example.org/error-ui", gjson.Get(stdout, "services.identity.config.selfservice.flows.error.ui_url").String())
		}, WithDefaultProject, WithPositionalProject)
	})

	t.Run("is able to remove a key", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "project", "--format", "json", "--remove", `/services/identity/config/selfservice/methods/password/enabled`)
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.password.enabled").Bool())
		}, WithDefaultProject, WithPositionalProject)
	})

	t.Run("fails if no opts are given", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "project", "--format", "json")
			require.Error(t, err, stdout)
		}, WithDefaultProject, WithPositionalProject)
	})

	t.Run("is able to update several keys", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "project", "--format", "json",
				"--replace", `/services/identity/config/selfservice/methods/link/enabled=true`,
				"--replace", `/services/identity/config/selfservice/methods/oidc/enabled=true`,
				"--remove", `/services/identity/config/selfservice/methods/profile/enabled`,
				"--remove", `/services/identity/config/selfservice/methods/password/enabled`,
				"--add", `/services/identity/config/selfservice/methods/totp/enabled=true`,
				"--add", `/services/identity/config/selfservice/methods/lookup_secret/enabled=true`,
				"-f", "fixtures/patch/1.json",
				"-f", "fixtures/patch/2.json",
			)
			require.NoError(t, err)
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.password.enabled").Bool())
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.profile.enabled").Bool())
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.link.enabled").Bool())
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.oidc.enabled").Bool())
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.totp.enabled").Bool())
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.lookup_secret.enabled").Bool())
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.webauthn.enabled").Bool())
			assert.True(t, gjson.Get(stdout, "services.identity.config.selfservice.methods.webauthn.config.passwordless").Bool())
			assert.Equal(t, "some value", gjson.Get(stdout, "services.identity.config.selfservice.methods.webauthn.config.rp.display_name").String())
		}, WithDefaultProject, WithPositionalProject)
	})
}
