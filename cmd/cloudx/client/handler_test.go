package client_test

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"testing"

	"github.com/ory/x/assertx"
	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	cloud "github.com/ory/client-go"
)

//go:embed fixtures/update_project/config.json
var config json.RawMessage

func TestCommandHelper(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	email, password := testhelpers.RegisterAccount(t, configDir)
	project := testhelpers.CreateProject(t, configDir)

	loggedIn := &client.CommandHelper{
		ConfigLocation:   configDir,
		NoConfirm:        true,
		IsQuiet:          true,
		VerboseWriter:    io.Discard,
		VerboseErrWriter: io.Discard,
		Ctx:              context.Background(),
	}
	assertValidProject := func(t *testing.T, actual *cloud.Project) {
		assert.Equal(t, project, actual.Id)
		assert.NotZero(t, actual.Slug)
		assert.NotZero(t, actual.Services.Identity.Config)
		assert.NotZero(t, actual.Services.Permission.Config)
	}
	reauth := func() *client.CommandHelper {
		notYetLoggedIn := *loggedIn
		notYetLoggedIn.ConfigLocation = testhelpers.NewConfigDir(t)
		notYetLoggedIn.IsQuiet = false
		notYetLoggedIn.PwReader = func() ([]byte, error) {
			return []byte(password), nil
		}
		notYetLoggedIn.Stdin = bufio.NewReader(bytes.NewBufferString(
			"y\n" + // Do you already have an Ory Console account you wish to use? [y/n]: y
				email + "\n")) // Email fakeEmail()
		return &notYetLoggedIn
	}

	t.Run("func=GetProject", func(t *testing.T) {
		t.Run("is able to get project", func(t *testing.T) {
			p, err := loggedIn.GetProject(project)
			require.NoError(t, err)
			assertValidProject(t, p)
		})

		t.Run("is not able to list projects if not authenticated and quiet flag", func(t *testing.T) {
			notLoggedIn := *loggedIn
			notLoggedIn.ConfigLocation = testhelpers.NewConfigDir(t)
			_, err := notLoggedIn.GetProject(project)
			assert.ErrorIs(t, err, client.ErrNoConfigQuiet)
		})

		t.Run("is able to get project after authenticating", func(t *testing.T) {
			notYetLoggedIn := reauth()
			p, err := notYetLoggedIn.GetProject(project)
			require.NoError(t, err)
			assertValidProject(t, p)
		})
	})

	t.Run("func=UpdateProject", func(t *testing.T) {
		t.Run("is able to update a project", func(t *testing.T) {
			res, err := loggedIn.UpdateProject(project, "", []json.RawMessage{config})
			require.NoError(t, err)

			assertx.EqualAsJSONExcept(t, config, res.Project, []string{
				"id",
				"revision_id",
				"state",
				"slug",
				"services.identity.config.serve",
				"services.identity.config.cookies",
				"services.identity.config.identity.default_schema_id",
				"services.identity.config.identity.schemas",
				"services.identity.config.session.cookie",
				"services.oauth2.config.urls.self",
				"services.oauth2.config.serve.public.tls",
				"services.oauth2.config.serve.tls",
				"services.oauth2.config.serve.admin.tls",
				"services.oauth2.config.serve.cookies.domain",
				"services.oauth2.config.serve.cookies.names",
				"services.oauth2.config.oauth2.session.encrypt_at_rest",
				"services.oauth2.config.oauth2.expose_internal_errors",
				"services.oauth2.config.oauth2.hashers",
				"services.oauth2.config.hsm",
				"services.oauth2.config.clients",
			})

			snapshotx.SnapshotT(t, res, snapshotx.ExceptPaths(
				"project.id",
				"project.revision_id",
				"project.slug",
				"project.services.identity.config.serve.public.base_url",
				"project.services.identity.config.serve.admin.base_url",
				"project.services.identity.config.session.cookie.domain",
				"project.services.identity.config.session.cookie.name",
				"project.services.identity.config.cookies.domain",
				"project.services.oauth2.config.urls.self",
				"project.services.oauth2.config.serve.cookies.domain",
				"project.services.oauth2.config.serve.cookies.names",
			))
		})

		t.Run("is able to update a projects name", func(t *testing.T) {
			name := testhelpers.FakeName()
			res, err := loggedIn.UpdateProject(project, name, []json.RawMessage{config})
			require.NoError(t, err)
			assert.Equal(t, name, res.Project.Name)
		})

		t.Run("is able to update a project after authenticating", func(t *testing.T) {
			notYetLoggedIn := reauth()
			res, err := notYetLoggedIn.UpdateProject(project, "", []json.RawMessage{config})
			require.NoError(t, err)
			assertValidProject(t, &res.Project)

			for _, w := range res.Warnings {
				t.Logf("Warning: %s", *w.Message)
			}
			assert.Len(t, res.Warnings, 0)
		})
	})
}
