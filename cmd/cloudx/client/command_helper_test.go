// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/randx"
	"github.com/ory/x/stringslice"

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
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	}); err != nil {
		panic(err)
	}
	testhelpers.UseStaging()
	m.Run()
}

func TestCommandHelper(t *testing.T) {
	ctx := client.ContextWithOptions(
		context.Background(),
		client.WithConfigLocation(testhelpers.NewConfigFile(t)),
		client.WithNoConfirm(true),
		client.WithQuiet(true),
		client.WithVerboseErrWriter(io.Discard),
		client.WithOpenBrowserHook(func(uri string) error {
			return errors.WithStack(fmt.Errorf("open browser hook not expected: %s", uri))
		}))

	email, password, name, sessionToken := testhelpers.RegisterAccount(ctx, t)
	defaultConfigFile := testhelpers.NewConfigFile(t)
	authenticated, err := client.NewCommandHelper(ctx, client.WithConfigLocation(defaultConfigFile), client.WithSessionToken(t, sessionToken))
	require.NoError(t, err)

	defaultWorkspace, err := authenticated.CreateWorkspace(ctx, randx.MustString(6, randx.AlphaNum))
	require.NoError(t, err)

	defaultWorkspaceAPIKey, err := authenticated.CreateWorkspaceAPIKey(ctx, defaultWorkspace.Id, randx.MustString(6, randx.AlphaNum))
	require.NoError(t, err)
	_ = defaultWorkspaceAPIKey

	defaultProject, err := authenticated.CreateProject(ctx, randx.MustString(6, randx.AlphaNum), "dev", nil, true)
	require.NoError(t, err)
	defaultWorkspaceProject, err := authenticated.CreateProject(ctx, randx.MustString(6, randx.AlphaNum), "dev", &defaultWorkspace.Id, false)
	require.NoError(t, err)

	assertValidProject := func(t *testing.T, actual *cloud.Project) {
		assert.NotZero(t, actual.Slug)
		assert.NotZero(t, actual.Services.Identity.Config)
		assert.NotZero(t, actual.Services.Permission.Config)
	}

	t.Run("func=SelectProjectWorkspace", func(t *testing.T) {
		t.Parallel()
		h, err := client.NewCommandHelper(ctx, client.WithSessionToken(t, sessionToken), client.WithConfigLocation(defaultConfigFile))
		require.NoError(t, err)
		otherProject, err := h.CreateProject(ctx, "other project", "dev", &defaultWorkspace.Id, false)
		require.NoError(t, err)

		t.Run("can change the selected project and workspace", func(t *testing.T) {
			// create new helper to ensure clean internal state
			h, err := client.NewCommandHelper(ctx, client.WithSessionToken(t, sessionToken), client.WithConfigLocation(defaultConfigFile))
			require.NoError(t, err)

			current, err := h.ProjectID()
			require.NoError(t, err)
			require.Equal(t, current, defaultProject.Id)

			require.NoError(t, h.SelectProject(otherProject.Id))
			require.NoError(t, h.SelectWorkspace(defaultWorkspace.Id))

			actualProject, err := h.ProjectID()
			require.NoError(t, err)
			assert.Equal(t, otherProject.Id, actualProject)

			actualWorkspace := h.WorkspaceID()
			require.NotNil(t, actualWorkspace)
			assert.Equal(t, defaultWorkspace.Id, *actualWorkspace)

			// check if persistent across instances
			h, err = client.NewCommandHelper(ctx, client.WithSessionToken(t, sessionToken), client.WithConfigLocation(defaultConfigFile))
			require.NoError(t, err)

			current, err = h.ProjectID()
			require.NoError(t, err)
			assert.Equal(t, current, otherProject.Id)
		})
	})

	t.Run("func=ListProjects", func(t *testing.T) {
		t.Parallel()

		configFile := testhelpers.NewConfigFile(t)
		_, _, _, sessionToken := testhelpers.RegisterAccount(ctx, t)

		h, err := client.NewCommandHelper(ctx, client.WithSessionToken(t, sessionToken), client.WithConfigLocation(configFile))
		require.NoError(t, err)

		t.Run("empty list", func(t *testing.T) {
			projects, err := h.ListProjects(ctx, nil)

			require.NoError(t, err)
			require.Empty(t, projects)
		})

		t.Run("list of projects", func(t *testing.T) {
			p0, err := h.CreateProject(ctx, "p0", "dev", nil, false)
			require.NoError(t, err)
			p1, err := h.CreateProject(ctx, "p1", "dev", nil, false)
			require.NoError(t, err)

			projects, err := h.ListProjects(ctx, nil)
			require.NoError(t, err)

			require.Len(t, projects, 2)
			assert.ElementsMatch(t, []string{p0.Id, p1.Id}, []string{projects[0].Id, projects[1].Id})
		})

		t.Run("list of workspace projects", func(t *testing.T) {
			workspace, err := h.CreateWorkspace(ctx, t.Name())
			require.NoError(t, err)
			p0, err := h.CreateProject(ctx, "p0", "dev", &workspace.Id, false)
			require.NoError(t, err)
			p1, err := h.CreateProject(ctx, "p1", "dev", &workspace.Id, false)
			require.NoError(t, err)

			projects, err := h.ListProjects(ctx, &workspace.Id)
			require.NoError(t, err)

			require.Len(t, projects, 2)
			assert.ElementsMatch(t, []string{p0.Id, p1.Id}, []string{projects[0].Id, projects[1].Id})
		})
	})

	t.Run("func=CreateProject", func(t *testing.T) {
		t.Parallel()
		configPath := testhelpers.NewConfigFile(t)

		h, err := client.NewCommandHelper(ctx, client.WithSessionToken(t, sessionToken), client.WithConfigLocation(configPath))
		require.NoError(t, err)
		workspace, err := h.CreateWorkspace(ctx, t.Name())
		require.NoError(t, err)

		name0 := "new project name0"
		name1 := "new project name1"
		name2 := "new project name2"

		project0, err := h.CreateProject(ctx, name0, "dev", &workspace.Id, true)
		require.NoError(t, err)
		project1, err := h.CreateProject(ctx, name1, "dev", nil, false)
		require.NoError(t, err)
		project2, err := h.CreateProject(ctx, name2, "dev", nil, false)
		require.NoError(t, err)

		assert.Len(t, stringslice.Unique([]string{project0.Id, project1.Id, project2.Id}), 3)
		assert.Len(t, stringslice.Unique([]string{project0.Slug, project1.Slug, project2.Slug}), 3)
		assert.Equal(t, []string{name0, name1, name2}, []string{project0.Name, project1.Name, project2.Name})

		assert.Equal(t, &workspace.Id, project0.WorkspaceId.Get())
		assert.Nil(t, project1.WorkspaceId.Get())
		assert.Nil(t, project2.WorkspaceId.Get())

		selectedID, err := h.ProjectID()
		require.NoError(t, err)
		assert.Equal(t, project0.Id, selectedID)
		assert.Equal(t, &workspace.Id, h.WorkspaceID())
	})

	t.Run("func=Authenticate", func(t *testing.T) {
		t.Parallel()

		_, page, cleanup := testhelpers.SetupPlaywright(t)
		t.Cleanup(cleanup)

		// ensure the browser has a valid session cookie
		testhelpers.BrowserLogin(t, page, email, password)
		t.Logf("browser login successful")

		// set up the command helper
		ctx := client.ContextWithOptions(ctx, client.WithConfigLocation(testhelpers.NewConfigFile(t)))
		h, err := client.NewCommandHelper(
			ctx,
			client.WithQuiet(false),
			client.WithOpenBrowserHook(testhelpers.PlaywrightAcceptConsentBrowserHook(t, page, password)),
		)
		require.NoError(t, err)

		// authenticate
		require.NoError(t, h.Authenticate(ctx))
		t.Logf("authentication successful")

		// we don't need playwright anymore
		cleanup()

		// assert success
		config, err := h.GetAuthenticatedConfig(ctx)
		require.NoError(t, err)
		require.NotNil(t, config)
		assert.Equal(t, email, config.IdentityTraits.Email)
		assert.Equal(t, name, config.IdentityTraits.Name)
		require.NotNil(t, config.AccessToken)
		assert.NotEmpty(t, config.AccessToken.AccessToken)

		// simple requests against all services to see if the token is valid and gets used
		p, err := h.GetProject(ctx, defaultProject.Id, nil)
		require.NoError(t, err)
		assert.Equal(t, defaultProject.Id, p.Id)
		t.Logf("project request successful")

		assert.JSONEq(t, "[]", testhelpers.ListIdentities(ctx, t, defaultProject.Id).Get("identities").Raw)
		t.Logf("list identities request successful")

		assert.JSONEq(t, "[]", testhelpers.ListClients(ctx, t, defaultProject.Id).Get("items").Raw)
		t.Logf("list clients request successful")

		assert.JSONEq(t, "[]", testhelpers.ListRelationTuples(ctx, t, defaultProject.Id).Get("relation_tuples").Raw)
		t.Logf("list relation tuples request successful")

		t.Run("refreshes and stores the refreshed oauth2 access token", func(t *testing.T) {
			ctx := client.ContextWithOptions(ctx, client.WithOpenBrowserHook(func(uri string) error {
				return fmt.Errorf("open browser hook not expected: %s", uri)
			}))

			oldToken := *config.AccessToken
			oldToken.Expiry = time.Unix(0, 0)
			oldConfig := *config
			oldConfig.AccessToken = &oldToken
			require.NoError(t, h.UpdateConfig(&oldConfig))

			actual, err := h.GetProject(ctx, defaultProject.Id, nil)
			require.NoError(t, err)
			assert.Equal(t, defaultProject.Id, actual.Id)

			newConfig, err := h.GetAuthenticatedConfig(ctx)
			require.NoError(t, err)
			require.NotNil(t, newConfig.AccessToken)
			assert.NotEmpty(t, newConfig.AccessToken.AccessToken)
			assert.NotEqual(t, oldToken.AccessToken, newConfig.AccessToken.AccessToken)
		})
	})

	t.Run("func=CreateProjectAPIKey and DeleteApiKey", func(t *testing.T) {
		t.Parallel()

		keyName := "a test key"

		key, err := authenticated.CreateProjectAPIKey(ctx, defaultProject.Id, keyName)
		require.NoError(t, err)
		assert.Equal(t, keyName, key.Name)
		assert.NotNil(t, keyName, key.Value)

		// check that the key works
		ctxWithKey := client.ContextWithOptions(ctx,
			client.WithWorkspaceAPIKey(sessionToken), // TODO this key should not be required, currently it is though to look up the slug
			client.WithProjectAPIKey(*key.Value))
		list := testhelpers.ListIdentities(ctxWithKey, t, defaultProject.Id)
		assert.True(t, list.Get("identities").Exists(), list.Raw)
		assert.True(t, list.Get("identities").IsArray(), list.Raw)

		require.NoError(t, authenticated.DeleteProjectAPIKey(ctx, defaultProject.Id, key.Id))

		_, stdErr, err := testhelpers.Cmd(ctxWithKey).Exec(nil, "list", "identities", "--project", defaultProject.Id)
		assert.ErrorIs(t, err, cmdx.ErrNoPrintButFail)
		assert.Contains(t, stdErr, "Access credentials are invalid")
	})

	t.Run("func=GetProject", func(t *testing.T) {
		for name, p := range map[string]*cloud.Project{
			"without workspace": defaultProject,
			"with workspace":    defaultWorkspaceProject,
		} {
			t.Run("is able to get project "+name, func(t *testing.T) {
				t.Parallel()

				actual, err := authenticated.GetProject(ctx, p.Id, p.WorkspaceId.Get())
				require.NoError(t, err)
				assert.Equal(t, p.Id, actual.Id)
				assertValidProject(t, p)

				actual, err = authenticated.GetProject(ctx, p.Slug[0:4], p.WorkspaceId.Get())
				require.NoError(t, err)
				assert.Equal(t, p.Id, actual.Id)
			})

			t.Run("is not able to get project if not authenticated and quiet flag "+name, func(t *testing.T) {
				t.Parallel()

				h, err := client.NewCommandHelper(ctx, client.WithQuiet(true))
				require.NoError(t, err)
				_, err = h.GetProject(ctx, p.Id, p.WorkspaceId.Get())
				assert.ErrorIs(t, err, client.ErrNoConfigQuiet)
			})
		}
	})

	t.Run("func=UpdateProject", func(t *testing.T) {
		t.Parallel()

		t.Run("is able to update a project", func(t *testing.T) {
			t.Skip("TODO")

			res, err := authenticated.UpdateProject(ctx, defaultProject.Id, "", []json.RawMessage{updatedProjectConfig})
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
			res, err := authenticated.UpdateProject(ctx, defaultProject.Id, name, []json.RawMessage{updatedProjectConfig})
			require.NoError(t, err)
			assert.Equal(t, name, res.Project.Name)
		})
	})
}
