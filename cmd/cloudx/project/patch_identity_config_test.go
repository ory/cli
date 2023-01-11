// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestPatchKratosConfig(t *testing.T) {
	t.Run("is able to replace a key", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "kratos-config", "--format", "json", "--replace", `/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
		}, DefaultProject|PositionalProject)
	})

	t.Run("is able to add a key", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "identity-config", "--format", "json", "--add", `/selfservice/methods/password/enabled=false`)
			require.NoError(t, err)
			assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
		}, DefaultProject|PositionalProject)
	})

	t.Run("is able to add a key with string", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "ic", "--format", "json", "--replace", "/selfservice/flows/error/ui_url=\"https://example.com/error-ui\"")
			require.NoError(t, err)
			assert.Equal(t, "https://example.com/error-ui", gjson.Get(stdout, "selfservice.flows.error.ui_url").String())
		}, DefaultProject|PositionalProject)
	})

	t.Run("fails if no opts are given", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "ic", "--format", "json")
			require.Error(t, err, stdout)
		}, DefaultProject|PositionalProject|FlagProject)
	})
}
