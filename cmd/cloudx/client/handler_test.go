// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

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
			"y\n" + // Do you want to sign in to an existing Ory Network account? [y/n]: y
				email + "\n")) // Email fakeEmail()
		return &notYetLoggedIn
	}

	t.Run("func=SetDefaultProject", func(t *testing.T) {
		t.Parallel()
		configDir := testhelpers.NewConfigDir(t)
		testhelpers.RegisterAccount(t, configDir)
		otherId := testhelpers.CreateProject(t, configDir)
		defaultId := testhelpers.CreateProject(t, configDir)
		testhelpers.SetDefaultProject(t, configDir, defaultId)

		cmdBase := client.CommandHelper{
			ConfigLocation: configDir,
		}

		t.Run("can change the selected project", func(t *testing.T) {
			cmd := cmdBase
			current := cmd.GetDefaultProjectID()
			assert.Equal(t, current, defaultId)

			err := cmd.SetDefaultProject(otherId)
			assert.NoError(t, err)

			selected := cmd.GetDefaultProjectID()
			assert.Equal(t, selected, otherId)
		})
	})

	t.Run("func=ListProjects", func(t *testing.T) {
		t.Parallel()
		configDir := testhelpers.NewConfigDir(t)
		testhelpers.RegisterAccount(t, configDir)

		cmdBase := client.CommandHelper{
			ConfigLocation: configDir,
		}

		t.Run("With no projects returns empty list", func(t *testing.T) {
			cmd := cmdBase

			projects, err := cmd.ListProjects()

			require.NoError(t, err)
			require.Empty(t, projects)
		})

		t.Run("With some projects returns list of projects", func(t *testing.T) {
			cmd := cmdBase
			project_name1 := "new_project_name1"
			project_name2 := "new_project_name2"

			project1, err := cmd.CreateProject(project_name1, false)
			require.NoError(t, err)
			project2, err := cmd.CreateProject(project_name2, false)
			require.NoError(t, err)

			projects, err := cmd.ListProjects()

			require.NoError(t, err)
			assert.Len(t, projects, 2)
			assert.ElementsMatch(t, []string{projects[0].Id, projects[1].Id}, []string{project1.Id, project2.Id})
		})
	})

	t.Run("func=CreateProject", func(t *testing.T) {
		t.Parallel()
		configDir := testhelpers.NewConfigDir(t)
		testhelpers.RegisterAccount(t, configDir)

		cmdBase := client.CommandHelper{
			ConfigLocation: configDir,
		}

		t.Run("creates project and sets default project", func(t *testing.T) {
			cmd := cmdBase
			project_name := "new_project_name"

			project, err := cmd.CreateProject(project_name, true)
			require.NoError(t, err)
			assert.Equal(t, project.Name, project_name)

			defaultId := cmd.GetDefaultProjectID()
			assert.Equal(t, project.Id, defaultId)
		})

		t.Run("creates two projects with different names", func(t *testing.T) {
			cmd := cmdBase
			project_name1 := "new_project_name1"
			project_name2 := "new_project_name2"

			project1, err := cmd.CreateProject(project_name1, true)
			require.NoError(t, err)

			project2, err := cmd.CreateProject(project_name2, false)
			require.NoError(t, err)

			assert.NotEqual(t, project1.Id, project2.Id)
			assert.NotEqual(t, project1.Name, project2.Name)
			assert.NotEqual(t, project1.Slug, project2.Slug)

			defaultId := cmd.GetDefaultProjectID()
			assert.Equal(t, project1.Id, defaultId)
		})
	})

	t.Run("func=Authenticate", func(t *testing.T) {
		t.Parallel()
		cmdBase := client.CommandHelper{
			ConfigLocation:   testhelpers.NewConfigDir(t),
			NoConfirm:        true,
			IsQuiet:          false,
			VerboseWriter:    io.Discard,
			VerboseErrWriter: io.Discard,
		}

		t.Run("create new account", func(t *testing.T) {
			cmd := cmdBase

			name := testhelpers.FakeName()
			email := testhelpers.FakeEmail()
			var r bytes.Buffer
			_, _ = r.WriteString("n\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: n
			_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
			_, _ = r.WriteString(name + "\n")  // Name: FakeName()
			_, _ = r.WriteString("n\n")        // Subscribe to the Ory Security Newsletter to get platform and security updates? [y/n]: n
			_, _ = r.WriteString("y\n")        // I accept the Terms of Service [y/n]: y
			cmd.Stdin = bufio.NewReader(&r)

			password := testhelpers.FakePassword()
			cmd.PwReader = func() ([]byte, error) { return []byte(password), nil }

			authCtx, err := cmd.Authenticate()

			require.NoError(t, err)
			require.NotNil(t, authCtx)
			require.Equal(t, authCtx.IdentityTraits.Email, email)
		})

		t.Run("log into existing account", func(t *testing.T) {
			cmd := cmdBase

			var r bytes.Buffer
			_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
			_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
			cmd.Stdin = bufio.NewReader(&r)

			cmd.PwReader = func() ([]byte, error) { return []byte(password), nil }

			auth_ctx, err := cmd.Authenticate()

			require.NoError(t, err)
			require.NotNil(t, auth_ctx)
			require.Equal(t, auth_ctx.IdentityTraits.Email, email)
		})

		t.Run("retry login after wrong password", func(t *testing.T) {
			cmd := cmdBase

			var r bytes.Buffer
			_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
			_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
			_, _ = r.WriteString(email + "\n") // Email: FakeEmail() [RETRY]
			cmd.Stdin = bufio.NewReader(&r)

			var retry = false
			cmd.PwReader = func() ([]byte, error) {
				if retry {
					return []byte(password), nil
				}
				retry = true
				return []byte("wrong"), nil
			}

			auth_ctx, err := cmd.Authenticate()

			require.NoError(t, err)
			require.NotNil(t, auth_ctx)
			require.Equal(t, auth_ctx.IdentityTraits.Email, email)
		})

		t.Run("switch logged in account", func(t *testing.T) {
			cmd := *loggedIn

			cmd.NoConfirm = false
			cmd.IsQuiet = false

			var r bytes.Buffer
			_, _ = r.WriteString("y\n")        // You are signed in as \"%s\" already. Do you wish to authenticate with another account?: y
			_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
			_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
			cmd.Stdin = bufio.NewReader(&r)

			cmd.PwReader = func() ([]byte, error) { return []byte(password), nil }

			auth_ctx, err := cmd.Authenticate()

			require.NoError(t, err)
			require.NotNil(t, auth_ctx)
			require.Equal(t, auth_ctx.IdentityTraits.Email, email)
		})
	})

	t.Run("func=CreateAPIKey and DeleteApiKey", func(t *testing.T) {
		t.Run("is able to get project", func(t *testing.T) {
			name := "a test key"
			token, err := loggedIn.CreateAPIKey(project, name)
			require.NoError(t, err)
			assert.Equal(t, name, token.Name)
			assert.NotEmpty(t, name, token.Value)

			require.NoError(t, loggedIn.DeleteAPIKey(project, token.Id))
		})
	})

	t.Run("func=GetProject", func(t *testing.T) {
		t.Run("is able to get project", func(t *testing.T) {
			p, err := loggedIn.GetProject(project)
			require.NoError(t, err)
			assertValidProject(t, p)

			actual, err := loggedIn.GetProject(p.Slug[0:4])
			require.NoError(t, err)
			assert.Equal(t, p, actual)
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
			require.NoErrorf(t, err, "%+v", err)

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
				"services.identity.config.selfservice.allowed_return_urls.0",
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
				"services.oauth2.config.oauth2.session",
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
				"project.services.identity.config.selfservice.allowed_return_urls.0",
				"project.services.oauth2.config.urls.self",
				"project.services.oauth2.config.serve.cookies.domain",
				"project.services.oauth2.config.serve.cookies.names",
				"project.services.identity.config.identity.schemas.1.url", // bucket changes locally vs staging
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
