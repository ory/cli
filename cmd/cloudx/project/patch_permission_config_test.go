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
	t.Parallel()

	for _, tc := range []struct {
		name string
		// doPatch will use the same project in parallel, so it is important to only do one operation per test
		doPatch func(t *testing.T, exec execFunc)
	}{
		{
			name: "is able to replace a key using keto-config",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "keto-config", "--format", "json", "--replace", `/namespaces=[{"name":"files", "id": 2}]`)
				require.NoError(t, err)
				assert.Equal(t, "files", gjson.Get(stdout, "namespaces.0.name").String(), stdout)
			},
		},
		{
			name: "is able to add a key using permission-config",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "permission-config", "--format", "json", "--replace", `/namespaces=[{"name":"docs", "id": 3}]`)
				require.NoError(t, err)
				assert.Equal(t, "docs", gjson.Get(stdout, "namespaces.0.name").String(), stdout)
			},
		},
		{
			name: "is able to replace a key using pc",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "pc", "--format", "json", "--replace", `/namespaces=[{"name":"people", "id": 4}]`)
				require.NoError(t, err)
				assert.Equal(t, "people", gjson.Get(stdout, "namespaces.0.name").String(), stdout)
			},
		},
		{
			name: "fails if no opts are given",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "pc", "--format", "json")
				require.Error(t, err, stdout)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			runWithProjectAsDefault(ctx, t, defaultProject.Id, tc.doPatch)
			runWithProjectAsFlag(ctx, t, extraProject.Id, tc.doPatch)
		})
	}
}
