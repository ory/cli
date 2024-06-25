// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"testing"

	"github.com/containerd/continuity/fs"

	"github.com/ory/x/assertx"
	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	cloud "github.com/ory/client-go"
)

//go:embed fixtures/update_project/config.json
var updatedProjectConfig json.RawMessage

func TestMain(m *testing.M) {
	testhelpers.RunAgainstStaging(m)
}

func TestCommandHelper(t *testing.T) {
	ctx := context.Background()
	configPath := testhelpers.NewConfigFile(t)
	email, password, name := testhelpers.RegisterAccount(t, configPath)
	project := testhelpers.CreateProject(t, configPath, nil)

	defaultOpts := func() []client.CommandHelperOption {
		return []client.CommandHelperOption{
			client.WithNoConfirm(true),
			client.WithQuiet(true),
			client.WithVerboseErrWriter(io.Discard),
		}
	}

	loggedIn, err := client.NewCommandHelper(
		ctx,
		append(defaultOpts(), client.WithConfigLocation(configPath))...,
	)
	require.NoError(t, err)
	assertValidProject := func(t *testing.T, actual *cloud.Project) {
		assert.NotZero(t, actual.Slug)
		assert.NotZero(t, actual.Services.Identity.Config)
		assert.NotZero(t, actual.Services.Permission.Config)
	}

	t.Run("func=SelectProject", func(t *testing.T) {
		t.Parallel()
		configDir := testhelpers.NewConfigFile(t)
		testhelpers.RegisterAccount(t, configDir)
		firstProject := testhelpers.CreateProject(t, configDir, nil)
		secondProject := testhelpers.CreateProject(t, configDir, nil)
		testhelpers.SetDefaultProject(t, configDir, secondProject.Id)

		t.Run("can change the selected project", func(t *testing.T) {
			h, err := client.NewCommandHelper(ctx, append(defaultOpts(), client.WithConfigLocation(configDir))...)
			require.NoError(t, err)

			current, err := h.ProjectID()
			require.NoError(t, err)
			require.Equal(t, current, secondProject.Id)

			require.NoError(t, h.SelectProject(firstProject.Id))

			selected, err := h.ProjectID()
			require.NoError(t, err)
			assert.Equal(t, selected, firstProject.Id)
		})
	})

	t.Run("func=ListProjects", func(t *testing.T) {
		t.Parallel()
		configFile := testhelpers.NewConfigFile(t)
		testhelpers.RegisterAccount(t, configFile)

		h, err := client.NewCommandHelper(ctx, append(defaultOpts(), client.WithConfigLocation(configFile))...)
		require.NoError(t, err)

		t.Run("empty list", func(t *testing.T) {
			projects, err := h.ListProjects(ctx, nil)

			require.NoError(t, err)
			require.Empty(t, projects)
		})

		t.Run("list of projects", func(t *testing.T) {
			p0, p1 := testhelpers.CreateProject(t, configFile, nil), testhelpers.CreateProject(t, configFile, nil)

			projects, err := h.ListProjects(ctx, nil)

			require.NoError(t, err)
			require.Len(t, projects, 2)
			assert.ElementsMatch(t, []string{p0.Id, p1.Id}, []string{projects[0].Id, projects[1].Id})
		})

		t.Run("list of workspace projects", func(t *testing.T) {
			workspace := testhelpers.CreateWorkspace(t, configFile)
			p0, p1 := testhelpers.CreateProject(t, configFile, &workspace), testhelpers.CreateProject(t, configFile, &workspace)

			projects, err := h.ListProjects(ctx, &workspace)

			require.NoError(t, err)
			require.Len(t, projects, 2)
			assert.ElementsMatch(t, []string{p0.Id, p1.Id}, []string{projects[0].Id, projects[1].Id})
		})
	})

	t.Run("func=CreateProject", func(t *testing.T) {
		t.Parallel()
		configPath := testhelpers.NewConfigFile(t)
		testhelpers.RegisterAccount(t, configPath)

		h, err := client.NewCommandHelper(ctx, append(defaultOpts(), client.WithConfigLocation(configPath))...)
		require.NoError(t, err)

		t.Run("creates project and sets default project", func(t *testing.T) {
			newName := "new project name"

			project, err := h.CreateProject(ctx, newName, "dev", nil, true)
			require.NoError(t, err)
			assert.Equal(t, project.Name, newName)

			defaultID, err := h.ProjectID()
			require.NoError(t, err)
			assert.Equal(t, project.Id, defaultID)
		})

		t.Run("creates two projects with different names", func(t *testing.T) {
			name1 := "new project name1"
			name2 := "new project name2"

			project1, err := h.CreateProject(ctx, name1, "dev", nil, true)
			require.NoError(t, err)

			project2, err := h.CreateProject(ctx, name2, "dev", nil, false)
			require.NoError(t, err)

			assert.NotEqual(t, project1.Id, project2.Id)
			assert.NotEqual(t, project1.Name, project2.Name)
			assert.NotEqual(t, project1.Slug, project2.Slug)

			selectedID, err := h.ProjectID()
			require.NoError(t, err)
			assert.Equal(t, project1.Id, selectedID)
		})
	})

	t.Run("func=Authenticate", func(t *testing.T) {
		t.Parallel()
		configPath := testhelpers.NewConfigFile(t)
		email2, password2, name2 := testhelpers.FakeAccount()

		t.Run("create new account", func(t *testing.T) {
			h, err := client.NewCommandHelper(ctx,
				client.WithConfigLocation(configPath),
				client.WithStdin(testhelpers.RegistrationBuffer(name2, email2)),
				client.WithPasswordReader(func() ([]byte, error) { return []byte(password2), nil }),
			)
			require.NoError(t, err)

			require.NoError(t, h.Authenticate(ctx))

			config, err := h.GetAuthenticatedConfig(ctx)
			require.NoError(t, err)
			require.NotNil(t, config)
			assert.Equal(t, email2, config.IdentityTraits.Email)
			assert.Equal(t, name2, config.IdentityTraits.Name)
		})

		t.Run("log into existing account", func(t *testing.T) {
			var r bytes.Buffer
			_, _ = r.WriteString("y\n")         // Do you want to sign in to an existing Ory Network account? [y/n]: y
			_, _ = r.WriteString(email2 + "\n") // Email: FakeEmail()
			h, err := client.NewCommandHelper(ctx,
				client.WithConfigLocation(testhelpers.NewConfigFile(t)),
				client.WithStdin(&r),
				client.WithPasswordReader(func() ([]byte, error) { return []byte(password2), nil }),
			)
			require.NoError(t, err)

			require.NoError(t, h.Authenticate(ctx))

			config, err := h.GetAuthenticatedConfig(ctx)
			require.NoError(t, err)
			require.NotNil(t, config)
			assert.Equal(t, email2, config.IdentityTraits.Email)
			assert.Equal(t, name2, config.IdentityTraits.Name)
		})

		t.Run("retry login after wrong password", func(t *testing.T) {
			var r bytes.Buffer
			_, _ = r.WriteString("y\n")         // Do you want to sign in to an existing Ory Network account? [y/n]: y
			_, _ = r.WriteString(email2 + "\n") // Email: FakeEmail()
			_, _ = r.WriteString(email2 + "\n") // Email: FakeEmail() [RETRY]

			retry := false
			pwReader := func() ([]byte, error) {
				if retry {
					return []byte(password2), nil
				}
				retry = true
				return []byte("wrong"), nil
			}

			h, err := client.NewCommandHelper(ctx,
				client.WithConfigLocation(testhelpers.NewConfigFile(t)),
				client.WithStdin(&r),
				client.WithPasswordReader(pwReader),
			)

			require.NoError(t, h.Authenticate(ctx))

			config, err := h.GetAuthenticatedConfig(ctx)
			require.NoError(t, err)
			require.NotNil(t, config)
			assert.Equal(t, email2, config.IdentityTraits.Email)
			assert.Equal(t, name2, config.IdentityTraits.Name)
		})

		t.Run("switch logged in account", func(t *testing.T) {
			newConfigPath := testhelpers.NewConfigFile(t)
			require.NoError(t, fs.CopyFile(newConfigPath, configPath))

			var r bytes.Buffer
			_, _ = r.WriteString("y\n")        // You are signed in as \"%s\" already. Do you wish to authenticate with another account?: y
			_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
			_, _ = r.WriteString(email + "\n") // Email: FakeEmail()

			h, err := client.NewCommandHelper(ctx,
				client.WithConfigLocation(newConfigPath),
				client.WithStdin(&r),
				client.WithPasswordReader(func() ([]byte, error) { return []byte(password), nil }),
			)

			require.NoError(t, h.Authenticate(ctx))

			config, err := h.GetAuthenticatedConfig(ctx)
			require.NoError(t, err)
			require.NotNil(t, config)
			assert.Equal(t, config.IdentityTraits.Email, email)
			assert.Equal(t, config.IdentityTraits.Name, name)
		})
	})

	t.Run("func=CreateAPIKey and DeleteApiKey", func(t *testing.T) {
		t.Run("is able to get project", func(t *testing.T) {
			name := "a test key"
			token, err := loggedIn.CreateAPIKey(ctx, project.Id, name)
			require.NoError(t, err)
			assert.Equal(t, name, token.Name)
			assert.NotEmpty(t, name, token.Value)

			require.NoError(t, loggedIn.DeleteAPIKey(ctx, project.Id, token.Id))
		})
	})

	t.Run("func=GetProject", func(t *testing.T) {
		t.Run("is able to get project", func(t *testing.T) {
			p, err := loggedIn.GetProject(ctx, project.Id, nil)
			require.NoError(t, err)
			assert.Equal(t, project.Id, p.Id)
			assertValidProject(t, p)

			actual, err := loggedIn.GetProject(ctx, p.Slug[0:4], nil)
			require.NoError(t, err)
			assert.Equal(t, p, actual)
		})

		t.Run("is able to get workspace project", func(t *testing.T) {
			workspace := testhelpers.CreateWorkspace(t, configPath)
			project := testhelpers.CreateProject(t, configPath, &workspace)
			p, err := loggedIn.GetProject(ctx, project.Id, &workspace)
			require.NoError(t, err)
			assert.Equal(t, project, p)
			assertValidProject(t, p)

			actual, err := loggedIn.GetProject(ctx, p.Slug[0:4], &workspace)
			require.NoError(t, err)
			assert.Equal(t, project, actual)
		})

		t.Run("is not able to get project if not authenticated and quiet flag", func(t *testing.T) {
			h, err := client.NewCommandHelper(ctx, append(
				defaultOpts(),
				client.WithConfigLocation(testhelpers.NewConfigFile(t)),
				client.WithQuiet(true),
			)...)
			require.NoError(t, err)
			_, err = h.GetProject(ctx, project.Id, nil)
			assert.ErrorIs(t, err, client.ErrNoConfigQuiet)
		})
	})

	t.Run("func=UpdateProject", func(t *testing.T) {
		t.Run("is able to update a project", func(t *testing.T) {
			t.Skip("TODO")

			res, err := loggedIn.UpdateProject(ctx, project.Id, "", []json.RawMessage{updatedProjectConfig})
			require.NoErrorf(t, err, "%+v", err)

			assertx.EqualAsJSONExcept(t, updatedProjectConfig, res.Project, []string{
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
			res, err := loggedIn.UpdateProject(ctx, project.Id, name, []json.RawMessage{updatedProjectConfig})
			require.NoError(t, err)
			assert.Equal(t, name, res.Project.Name)
		})
	})
}
