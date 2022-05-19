package project_test

import (
	_ "embed"
	"encoding/json"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/ory/x/assertx"
	"strings"
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/update/json/config.json
var fixtureProject json.RawMessage

//go:embed fixtures/update-kratos/json/config.json
var fixtureKratosConfig json.RawMessage

func TestUpdateProject(t *testing.T) {
	project := testhelpers.CreateProject(t, defaultConfig)

	for _, tc := range []struct {
		subcommand, pathSuccess, pathFailure string
		fixture                              json.RawMessage
	}{
		{
			subcommand:  "project",
			pathSuccess: "fixtures/update/json/config.json",
			pathFailure: "fixtures/update/fail/config.json",
			fixture:     fixtureProject,
		},
		{
			subcommand:  "identity-config",
			pathSuccess: "fixtures/update-kratos/json/config.json",
			pathFailure: "fixtures/update-kratos/fail/config.json",
			fixture:     fixtureKratosConfig,
		},
	} {
		t.Run("target="+tc.subcommand, func(t *testing.T) {
			t.Run("is able to update a project", func(t *testing.T) {
				stdout, _, err := defaultCmd.Exec(nil, "update", tc.subcommand, project, "--format", "json", "--file", tc.pathSuccess)
				require.NoError(t, err)

				assertx.EqualAsJSONExcept(t, fixtureProject, json.RawMessage(stdout), []string{
					// for project cmd
					"id",
					"revision_id",
					"state",
					"slug",
					"services.identity.config.serve",
					"services.identity.config.cookies",
					"services.identity.config.identity.default_schema_id",
					"services.identity.config.identity.schemas",
					"services.identity.config.session.cookie",
					// for kratos cmd
					"serve",
					"cookies",
					"identity.default_schema_id",
					"identity.schemas",
					"session.cookie",
					"courier.smtp.from_name",
					// for keto cmd
				})

				snapshotx.SnapshotT(t, json.RawMessage(stdout), snapshotx.ExceptPaths(
					// for project cmd
					"id",
					"revision_id",
					"slug",
					"services.identity.config.serve.public.base_url",
					"services.identity.config.serve.admin.base_url",
					"services.identity.config.session.cookie.domain",
					"services.identity.config.session.cookie.name",
					"services.identity.config.cookies.domain",
					// for kratos cmd
					"serve.public.base_url",
					"serve.admin.base_url",
					"session.cookie.domain",
					"session.cookie.name",
					"cookies.domain",
					"courier.smtp.from_name",
					// for keto cmd
				))
			})

			t.Run("prints good error messages for failing schemas", func(t *testing.T) {
				_, stderr, err := defaultCmd.Exec(nil, "update", tc.subcommand, project, "--format", "json", "--file", tc.pathFailure)
				require.ErrorIs(t, err, cmdx.ErrNoPrintButFail)

				t.Run("stdout", func(t *testing.T) {
					cupaloy.New(
						cupaloy.CreateNewAutomatically(true),
						cupaloy.FailOnUpdate(true),
						cupaloy.SnapshotFileExtension(".txt"),
					).SnapshotT(t, strings.SplitN(stderr, "\n", 2)[1])
				})
				t.Run("stderr", func(t *testing.T) {
					assert.Contains(t, stderr, "oneOf failed")
				})
			})
		})
	}
}
