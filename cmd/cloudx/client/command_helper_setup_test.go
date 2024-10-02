// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"

	cloud "github.com/ory/client-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/x/pointerx"
)

func TestDetermineIDs(t *testing.T) {
	wsID := uuid.Must(uuid.NewV4())
	wsName := "workspace-name"
	pjID := uuid.Must(uuid.NewV4())
	pjSlug := "slg-" + uuid.Must(uuid.NewV4()).String()
	requests := make(map[string]struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")] = struct{}{}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/workspaces":
			if err := json.NewEncoder(w).Encode(map[string]any{
				"workspaces": []cloud.Workspace{{
					Id:   wsID.String(),
					Name: wsName,
				}},
				"has_next_page":   false,
				"next_page_token": "",
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				t.Logf("unable to encode response: %s", err)
			}
		case "/projects":
			if err := json.NewEncoder(w).Encode([]cloud.ProjectMetadata{{
				Id:   pjID.String(),
				Slug: pjSlug,
			}}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				t.Logf("unable to encode response: %s", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	ctx := context.Background()

	setup := func(t *testing.T) (*CommandHelper, *Config) {
		h := &CommandHelper{
			configLocation:     t.TempDir() + "/config.json",
			cloudConsoleAPIURL: Ptr(ts.URL),
		}
		cfg, err := h.getOrCreateConfig()
		require.NoError(t, err)
		cfg.AccessToken = &oauth2.Token{
			AccessToken: "req_id_" + uuid.Must(uuid.NewV4()).String(),
			Expiry:      time.Now().Add(time.Hour),
		}
		cfg.isAuthenticated = true
		h.config = cfg
		return h, cfg
	}

	t.Run("case=no values", func(t *testing.T) {
		h, cfg := setup(t)

		require.NoError(t, h.determineWorkspaceID(ctx, cfg))
		require.NoError(t, h.determineProjectID(ctx, cfg))

		assert.Zero(t, h.workspaceID)
		assert.Zero(t, h.projectID)

		assert.NotContains(t, requests, cfg.AccessToken.AccessToken)
	})

	t.Run("case=from config", func(t *testing.T) {
		h, cfg := setup(t)

		cfg.SelectedWorkspace = uuid.Must(uuid.NewV4())
		cfg.SelectedProject = uuid.Must(uuid.NewV4())

		require.NoError(t, h.determineWorkspaceID(ctx, cfg))
		require.NoError(t, h.determineProjectID(ctx, cfg))

		assert.Equal(t, cfg.SelectedWorkspace, h.workspaceID)
		assert.Equal(t, cfg.SelectedProject, h.projectID)

		assert.NotContains(t, requests, cfg.AccessToken.AccessToken)
	})

	t.Run("case=with workspace API key", func(t *testing.T) {
		t.Run("errors when workspace override is set", func(t *testing.T) {
			h, cfg := setup(t)

			h.workspaceAPIKey = Ptr("workspace-api-key")
			h.workspaceOverride = Ptr("workspace-override")

			assert.ErrorContains(t, h.determineWorkspaceID(ctx, cfg), "workspace API key is set but workspace flag is also set, please remove one")
		})

		h, cfg := setup(t)

		h.workspaceAPIKey = Ptr("workspace-api-key")

		require.NoError(t, h.determineWorkspaceID(ctx, cfg))
		assert.Equal(t, wsID, h.workspaceID)
		assert.Contains(t, requests, *h.workspaceAPIKey)
	})

	t.Run("case=with project API key", func(t *testing.T) {
		t.Run("errors when project override is set", func(t *testing.T) {
			h, cfg := setup(t)

			h.projectAPIKey = Ptr("project-api-key")
			h.projectOverride = Ptr("project-override")

			assert.ErrorContains(t, h.determineProjectID(ctx, cfg), "project API key is set but project flag is also set, please remove one")
		})
		t.Run("errors when workspace is set", func(t *testing.T) {
			h, cfg := setup(t)

			h.projectAPIKey = Ptr("project-api-key")
			h.workspaceID = uuid.Must(uuid.NewV4())

			assert.ErrorContains(t, h.determineProjectID(ctx, cfg), "project API key is set but workspace is also set, please remove one")
		})

		h, cfg := setup(t)

		h.projectAPIKey = Ptr("project-api-key")

		require.NoError(t, h.determineProjectID(ctx, cfg))
		assert.Equal(t, pjID, h.projectID)
		assert.Contains(t, requests, *h.projectAPIKey)
	})

	t.Run("case=with workspace overrides", func(t *testing.T) {
		for _, tc := range []struct {
			name, value         string
			expectHandlerCalled bool
		}{
			{
				name:                "full ID",
				value:               wsID.String(),
				expectHandlerCalled: false,
			},
			{
				name:                "partial ID",
				value:               wsID.String()[0:5],
				expectHandlerCalled: true,
			},
			{
				name:                "full name",
				value:               wsName,
				expectHandlerCalled: true,
			},
			{
				name:                "partial name",
				value:               wsName[0:5],
				expectHandlerCalled: true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				for _, setValue := range []func(*CommandHelper){
					func(h *CommandHelper) { h.workspaceOverride = Ptr(tc.value) },
					func(*CommandHelper) { t.Setenv(WorkspaceKey, tc.value) },
				} {
					h, cfg := setup(t)

					setValue(h)

					require.NoError(t, h.determineWorkspaceID(ctx, cfg))
					assert.Equal(t, wsID, h.workspaceID)
					if tc.expectHandlerCalled {
						assert.Contains(t, requests, cfg.AccessToken.AccessToken)
					} else {
						assert.NotContains(t, requests, cfg.AccessToken.AccessToken)
					}
				}
			})
		}
	})

	t.Run("case=with project overrides", func(t *testing.T) {
		for _, tc := range []struct {
			name, value         string
			expectHandlerCalled bool
		}{
			{
				name:                "full ID",
				value:               pjID.String(),
				expectHandlerCalled: false,
			},
			{
				name:                "partial ID",
				value:               pjID.String()[0:5],
				expectHandlerCalled: true,
			},
			{
				name:                "full slug",
				value:               pjSlug,
				expectHandlerCalled: true,
			},
			{
				name:                "partial slug",
				value:               pjSlug[0:5],
				expectHandlerCalled: true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				for _, setValue := range []func(*CommandHelper){
					func(h *CommandHelper) { h.projectOverride = Ptr(tc.value) },
					func(*CommandHelper) { t.Setenv(ProjectKey, tc.value) },
				} {
					h, cfg := setup(t)

					setValue(h)

					require.NoError(t, h.determineProjectID(ctx, cfg))
					assert.Equal(t, pjID, h.projectID)
					if tc.expectHandlerCalled {
						assert.Contains(t, requests, cfg.AccessToken.AccessToken)
					} else {
						assert.NotContains(t, requests, cfg.AccessToken.AccessToken)
					}
				}
			})
		}
	})
}
