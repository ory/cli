// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// revokeRecorder answers OAuth2 token revocation requests with the given
// status and fails the test on any other request, such as a token refresh.
func revokeRecorder(t *testing.T, status int) (*http.Client, *[]url.Values) {
	var revoked []url.Values
	return &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		res := &http.Response{
			StatusCode: status,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    r,
		}
		if r.Method != http.MethodPost || r.URL.Path != "/oauth2/revoke" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			res.StatusCode = http.StatusBadRequest
			return res, nil
		}
		if assert.NoError(t, r.ParseForm()) {
			revoked = append(revoked, r.PostForm)
		}
		return res, nil
	})}, &revoked
}

func TestSignOut(t *testing.T) {
	expired := time.Now().Add(-time.Hour)

	writeConfig := func(t *testing.T, config any) string {
		location := testhelpers.NewConfigFile(t)
		raw, err := json.Marshal(config)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(location, raw, 0600))
		return location
	}

	signOut := func(t *testing.T, location string, status int, opts ...client.CommandHelperOption) (revoked []url.Values, stderr string, err error) {
		httpClient, rec := revokeRecorder(t, status)
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
		out := bytes.Buffer{}
		h, err := client.NewCommandHelper(ctx, append([]client.CommandHelperOption{
			client.WithConfigLocation(location),
			client.WithVerboseErrWriter(&out),
			client.WithOpenBrowserHook(func(string) error { return testhelpers.ErrAuthFlowTriggered }),
		}, opts...)...)
		require.NoError(t, err)
		err = h.SignOut(ctx)
		return *rec, out.String(), err
	}

	t.Run("revokes the refresh token of an expired session without signing in", func(t *testing.T) {
		location := writeConfig(t, client.Config{
			Version:     client.ConfigVersion,
			AccessToken: &oauth2.Token{AccessToken: "access", RefreshToken: "refresh", Expiry: expired},
		})

		revoked, _, err := signOut(t, location, http.StatusOK)
		require.NoError(t, err)

		assert.Equal(t, []url.Values{{
			"client_id":       {"ory-cli"},
			"token":           {"refresh"},
			"token_type_hint": {"refresh_token"},
		}}, revoked)
		assert.Nil(t, testhelpers.ReadConfig(t, location).AccessToken)
	})

	t.Run("revokes the access token if there is no refresh token", func(t *testing.T) {
		location := writeConfig(t, client.Config{
			Version:     client.ConfigVersion,
			AccessToken: &oauth2.Token{AccessToken: "access", Expiry: expired},
		})

		revoked, _, err := signOut(t, location, http.StatusOK)
		require.NoError(t, err)

		assert.Equal(t, []url.Values{{
			"client_id":       {"ory-cli"},
			"token":           {"access"},
			"token_type_hint": {"access_token"},
		}}, revoked)
		assert.Nil(t, testhelpers.ReadConfig(t, location).AccessToken)
	})

	t.Run("signs out locally if revocation fails", func(t *testing.T) {
		location := writeConfig(t, client.Config{
			Version:     client.ConfigVersion,
			AccessToken: &oauth2.Token{AccessToken: "access", RefreshToken: "refresh", Expiry: time.Now().Add(time.Hour)},
		})

		revoked, stderr, err := signOut(t, location, http.StatusServiceUnavailable)
		require.NoError(t, err)

		assert.Len(t, revoked, 1)
		assert.Contains(t, stderr, "failed to revoke")
		assert.Nil(t, testhelpers.ReadConfig(t, location).AccessToken)
	})

	t.Run("does nothing without a config file", func(t *testing.T) {
		location := testhelpers.NewConfigFile(t)

		revoked, _, err := signOut(t, location, http.StatusOK)
		require.NoError(t, err)

		assert.Empty(t, revoked)
		assert.NoFileExists(t, location)
	})

	t.Run("clears a legacy config", func(t *testing.T) {
		location := writeConfig(t, map[string]string{"version": "v0alpha0", "session_token": "legacy"})

		revoked, _, err := signOut(t, location, http.StatusOK, client.WithNoConfirm(true))
		require.NoError(t, err)

		assert.Empty(t, revoked)
		raw, err := os.ReadFile(location)
		require.NoError(t, err)
		assert.NotContains(t, string(raw), "legacy")
		assert.Equal(t, client.ConfigVersion, testhelpers.ReadConfig(t, location).Version)
	})
}
