package project_test

import (
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
	project := testhelpers.CreateProject(t, defaultConfig)
	t.Run("is able to update a project", func(t *testing.T) {
		stdout, _, err := defaultCmd.Exec(nil, "update", "kratos-config", project, "--format", "json", "--file", "./fixtures/update-kratos/json/config.json")
		require.NoError(t, err)

		assertx.EqualAsJSONExcept(t, json.RawMessage(fixtureKratosConfig), json.RawMessage(stdout), []string{
			"serve",
			"cookies",
			"identity.default_schema_id",
			"identity.schemas",
			"session.cookie",
			"courier.smtp.from_name",
		})

		snapshotx.SnapshotT(t, json.RawMessage(stdout), snapshotx.ExceptPaths(
			"serve.public.base_url",
			"serve.admin.base_url",
			"session.cookie.domain",
			"session.cookie.name",
			"cookies.domain",
			"courier.smtp.from_name",
		))
	})

	t.Run("prints good error messages for failing schemas", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "update", "identity-config", project, "--format", "json", "--file", "./fixtures/update-kratos/fail/config.json")
		require.ErrorIs(t, err, cmdx.ErrNoPrintButFail)

		t.Run("stdout", func(t *testing.T) {
			snapshotx.SnapshotT(t, stdout)
		})

		t.Run("stderr", func(t *testing.T) {
			assert.Contains(t, stderr, "oneOf failed")
		})
	})

	t.Run("is able to update a project after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		_, _, err := cmd.Exec(r, "update", "ic", project, "--format", "json", "--file", "./fixtures/update-kratos/json/config.json")
		require.NoError(t, err)
	})
}
