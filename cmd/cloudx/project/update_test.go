package project_test

import (
	_ "embed"
	"encoding/json"
	"github.com/ory/x/assertx"
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
)

//go:embed fixtures/update/json/config.json
var fixture []byte

func TestUpdateProject(t *testing.T) {
	project := testhelpers.CreateProject(t, defaultConfig)

	t.Run("is able to update a project", func(t *testing.T) {
		stdout, _, err := defaultCmd.Exec(nil, "update", "project", project, "--format", "json", "--file", "./fixtures/update/json/config.json")
		require.NoError(t, err)

		assertx.EqualAsJSONExcept(t, json.RawMessage(fixture), json.RawMessage(stdout), []string{
			"id",
			"revision_id",
			"state",
			"slug",
			"services.identity.config.serve",
			"services.identity.config.cookies",
			"services.identity.config.identity.default_schema_id",
			"services.identity.config.identity.schemas",
			"services.identity.config.session.cookie",
		})

		snapshotx.SnapshotT(t, json.RawMessage(stdout), snapshotx.ExceptPaths(
			"id",
			"revision_id",
			"slug",
			"services.identity.config.serve.public.base_url",
			"services.identity.config.serve.admin.base_url",
			"services.identity.config.session.cookie.domain",
			"services.identity.config.session.cookie.name",
			"services.identity.config.cookies.domain",
		))
	})
	t.Run("is able to update a projects name", func(t *testing.T) {
		name := testhelpers.FakeName()
		stdout, _, err := defaultCmd.Exec(nil, "update", "project", project, "--name", name, "--format", "json", "--file", "./fixtures/update/json/config.json")
		require.NoError(t, err)
		assert.Equal(t, name, gjson.Get(stdout, "name").String())
	})

	t.Run("prints good error messages for failing schemas", func(t *testing.T) {
		updatedName := testhelpers.TestProjectName()
		stdout, stderr, err := defaultCmd.Exec(nil, "update", "project", project, "--name", updatedName, "--format", "json", "--file", "./fixtures/update/fail/config.json")
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
		stdout, stderr, err := cmd.Exec(r, "update", "project", project, "--format", "json", "--file", "./fixtures/update/json/config.json")
		require.NoError(t, err, "stdout: %s\nstderr: %s", stdout, stderr)
	})
}
