// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestPatchPermissionConfig(t *testing.T) {
	t.Run("is able to replace a key using keto-config", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "keto-config", "--format", "json", "--replace", `/namespaces=[{"name":"files", "id": 2}]`)
			require.NoError(t, err)
			assert.Equal(t, "files", gjson.Get(stdout, "namespaces.0.name").String(), stdout)
		}, DefaultProject|PositionalProject)
	})

	t.Run("is able to add a key using permission-config", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			_, _, err := exec(nil, "patch", "permission-config", "--format", "json", "--replace", `/namespaces=[]`)
			require.NoError(t, err)

			stdout, _, err := exec(nil, "patch", "permission-config", "--format", "json", "--add", `/namespaces/0={"name":"docs", "id": 3}`)
			require.NoError(t, err)
			assert.Equal(t, "docs", gjson.Get(stdout, "namespaces.0.name").String(), stdout)
		}, DefaultProject|PositionalProject)
	})

	t.Run("is able to replace a key using pc", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "pc", "--format", "json", "--replace", `/namespaces=[{"name":"people", "id": 4}]`)
			require.NoError(t, err)
			assert.Equal(t, "people", gjson.Get(stdout, "namespaces.0.name").String(), stdout)
		}, DefaultProject|PositionalProject)
	})

	t.Run("fails if no opts are given", func(t *testing.T) {
		runWithProject(t, func(t *testing.T, exec execFunc, _ string) {
			stdout, _, err := exec(nil, "patch", "pc", "--format", "json")
			require.Error(t, err, stdout)
		}, DefaultProject|PositionalProject|FlagProject)
	})
}
