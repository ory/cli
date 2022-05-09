package cloudx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestAuthLogout(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	testhelpers.RegisterAccount(t, configDir)

	exec := testhelpers.ConfigAwareCmd(configDir)
	_, _, err := exec.Exec(nil, "auth", "logout")
	require.NoError(t, err)

	ac := testhelpers.ReadConfig(t, configDir)
	assert.Empty(t, ac.SessionToken)
}
