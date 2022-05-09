package project_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/assertx"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
)

//go:embed fixtures/update-kratos/json/config.json
var fixtureKratosConfig []byte

func TestProjectIdentityConfig(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	cmd := testhelpers.ConfigAwareCmd(configDir)
	email, password := testhelpers.RegisterAccount(t, configDir)

	project := testhelpers.CreateProject(t, configDir)
	t.Run("is able to update a project", func(t *testing.T) {
		stdout, _, err := cmd.ExecDebug(t, nil, "update", "kratos-config", project, "--format", "json", "--file", "./fixtures/update-kratos/json/config.json")
		require.NoError(t, err)

		assertx.EqualAsJSONExcept(t, json.RawMessage(fixtureKratosConfig), json.RawMessage(stdout), []string{
			"serve",
			"cookies",
			"identity.default_schema_id",
			"identity.schemas",
			"session.cookie",
			"courier.smtp.from_name",
		})

		snapshotx.SnapshotTExcept(t, json.RawMessage(stdout), []string{
			"serve.public.base_url",
			"serve.admin.base_url",
			"session.cookie.domain",
			"session.cookie.name",
			"cookies.domain",
			"courier.smtp.from_name",
		})
	})

	t.Run("prints good error messages for failing schemas", func(t *testing.T) {
		stdout, stderr, err := cmd.ExecDebug(t, nil, "update", "identity-config", project, "--format", "json", "--file", "./fixtures/update-kratos/fail/config.json")
		require.ErrorIs(t, err, cmdx.ErrNoPrintButFail)

		t.Run("stdout", func(t *testing.T) {
			snapshotx.SnapshotTExcept(t, stdout, nil)
		})

		t.Run("stderr", func(t *testing.T) {
			assert.Contains(t, stderr, "oneOf failed")
		})
	})

	t.Run("is able to update a project after authenticating", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigPasswordAwareCmd(configDir, password)
		// Create the account
		var r bytes.Buffer
		r.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
		r.WriteString(email + "\n") // Email fakeEmail()
		_, _, err := cmd.ExecDebug(t, &r, "update", "ic", project, "--format", "json", "--file", "./fixtures/update-kratos/json/config.json")
		require.NoError(t, err)
	})
}
