// Copyright © 2022 Ory Corp

package project_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/assertx"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
)

var (
	//go:embed fixtures/update/json/config.json
	fixtureProject json.RawMessage
	//go:embed fixtures/update-kratos/json/config.json
	fixtureKratosConfig json.RawMessage
	//go:embed fixtures/update-keto/json/config.json
	fixtureKetoConfig json.RawMessage
	//go:embed fixtures/update-hydra/json/config.json
	fixtureHydraConfig json.RawMessage
)

func TestUpdateProject(t *testing.T) {
	project := testhelpers.CreateProject(t, defaultConfig)

	for _, tc := range []struct {
		subcommand,
		pathSuccess,
		pathFailure,
		failureContains string
		fixture json.RawMessage
	}{
		{
			subcommand:      "project",
			pathSuccess:     "fixtures/update/json/config.json",
			pathFailure:     "fixtures/update/fail/config.json",
			failureContains: "oneOf failed",
			fixture:         fixtureProject,
		},
		{
			subcommand:      "identity-config",
			pathSuccess:     "fixtures/update-kratos/json/config.json",
			pathFailure:     "fixtures/update-kratos/fail/config.json",
			failureContains: "oneOf failed",
			fixture:         fixtureKratosConfig,
		},
		{
			subcommand:      "permission-config",
			pathSuccess:     "fixtures/update-keto/json/config.json",
			pathFailure:     "fixtures/update-keto/fail/config.json",
			failureContains: "cannot unmarshal string into Go struct field",
			fixture:         fixtureKetoConfig,
		},
		{
			subcommand:      "oauth2-config",
			pathSuccess:     "fixtures/update-hydra/json/config.json",
			pathFailure:     "fixtures/update-hydra/fail/config.json",
			failureContains: "cannot unmarshal number into Go struct field",
			fixture:         fixtureHydraConfig,
		},
	} {
		t.Run("target="+tc.subcommand, func(t *testing.T) {
			t.Run("is able to update a project", func(t *testing.T) {
				stdout, _, err := defaultCmd.Exec(nil, "update", tc.subcommand, project, "--format", "json", "--file", tc.pathSuccess)
				require.NoError(t, err)

				assertx.EqualAsJSONExcept(t, tc.fixture, json.RawMessage(stdout), []string{
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
					"services.identity.config.session.cookie",
					"services.oauth2.config.serve.cookies.names",
					"services.oauth2.config.serve.cookies.domain",
					"services.oauth2.config.urls.self",
					"services.oauth2.config.oauth2.session",
					// for kratos cmd
					"serve",
					"cookies",
					"identity.default_schema_id",
					"identity.schemas",
					"session.cookie",
					"courier.smtp.from_name",
					// for keto cmd
					// for hydra cmd
					"serve.cookies.names",
					"serve.cookies.domain",
					"urls.self",
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
					"services.oauth2.config.serve.cookies.names",
					"services.oauth2.config.serve.cookies.domain",
					"services.oauth2.config.urls.self",
					// for kratos cmd
					"serve.public.base_url",
					"serve.admin.base_url",
					"session.cookie.domain",
					"session.cookie.name",
					"cookies.domain",
					"courier.smtp.from_name",
					// for keto cmd
					// for hydra cmd
					"serve.cookies.names",
					"serve.cookies.domain",
					"urls.self",
					// bucket changes locally vs staging
					"services.identity.config.identity.schemas.1.url",
					"identity.schemas.1.url",
				))
			})

			t.Run("prints good error messages for failing schemas", func(t *testing.T) {
				stdout, stderr, err := defaultCmd.Exec(nil, "update", tc.subcommand, project, "--format", "json", "--file", tc.pathFailure)
				require.ErrorIs(t, err, cmdx.ErrNoPrintButFail)

				t.Run("stdout", func(t *testing.T) {
					snapshotx.SnapshotT(t, stdout)
				})
				t.Run("stderr", func(t *testing.T) {
					assert.Contains(t, stderr, tc.failureContains)
				})
			})
		})
	}
}
