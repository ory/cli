// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestPatchHydraConfig(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		doPatch func(t *testing.T, exec execFunc)
	}{
		{
			name: "is able to replace a key",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "hydra-config", "--format", "json", "--replace", `/strategies/access_token="jwt"`)
				require.NoError(t, err)
				assert.Equal(t, "jwt", gjson.Get(stdout, "strategies.access_token").String())
			},
		},
		{
			name: "is able to add a key",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "oauth2-config", "--format", "json", "--add", `/ttl/login_consent_request="1h"`)
				require.NoError(t, err)
				assert.Equal(t, "1h0m0s", gjson.Get(stdout, "ttl.login_consent_request").String())
			},
		},
		{
			name: "is able to add a key with string",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "oc", "--format", "json", "--replace", `/ttl/refresh_token="2h"`)
				require.NoError(t, err)
				assert.Equal(t, "2h0m0s", gjson.Get(stdout, "ttl.refresh_token").String())
			},
		},
		{
			name: "fails if no opts are given",
			doPatch: func(t *testing.T, exec execFunc) {
				stdout, _, err := exec(nil, "patch", "oc", "--format", "json")
				require.Error(t, err, stdout)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runWithProjectAsDefault(ctx, t, defaultProject.Id, tc.doPatch)
			runWithProjectAsFlag(ctx, t, extraProject.Id, tc.doPatch)
		})
	}
}
