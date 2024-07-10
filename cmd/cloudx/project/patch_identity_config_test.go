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
	t.Parallel()

	for _, tc := range []struct {
		name    string
		doPatch func(t *testing.T, exec execFunc)
	}{
		{
			name: "is able to replace a key",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "kratos-config", "--format", "json", "--replace", `/selfservice/methods/password/enabled=false`)
				require.NoError(t, err)
				assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
			},
		},
		{
			name: "is able to add a key",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "identity-config", "--format", "json", "--add", `/selfservice/methods/password/enabled=false`)
				require.NoError(t, err)
				assert.False(t, gjson.Get(stdout, "selfservice.methods.password.enabled").Bool())
			},
		},
		{
			name: "is able to add a key with string",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "ic", "--format", "json", "--replace", "/selfservice/flows/error/ui_url=\"https://example.com/error-ui\"")
				require.NoError(t, err)
				assert.Equal(t, "https://example.com/error-ui", gjson.Get(stdout, "selfservice.flows.error.ui_url").String())
			},
		},
		{
			name: "fails if no opts are given",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "ic", "--format", "json")
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
