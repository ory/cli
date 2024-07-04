// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestAuthLogout(t *testing.T) {
	configDir := testhelpers.NewConfigFile(t)
	testhelpers.RegisterAccount(t, configDir)

	exec := testhelpers.Cmd(configDir)
	_, _, err := exec.Exec(nil, "auth", "logout")
	require.NoError(t, err)

	ac := testhelpers.ReadConfig(t, configDir)
	assert.Empty(t, ac.SessionToken)
}
