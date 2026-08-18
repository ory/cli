// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/fetcher"
)

// TestOPLLocation covers https://github.com/ory/cli/issues/321: the permission
// config reports where the Ory Permission Language file lives rather than its
// contents, and `ory get opl` has to resolve that pointer.
func TestOPLLocation(t *testing.T) {
	for _, tc := range []struct {
		name    string
		config  map[string]interface{}
		want    string
		wantErr error
	}{
		{
			name:   "case=reads the location of an OPL file",
			config: map[string]interface{}{"namespaces": map[string]interface{}{"location": "https://example.com/opl.bin"}},
			want:   "https://example.com/opl.bin",
		},
		{
			name:    "case=legacy namespace definitions are reported as such",
			config:  map[string]interface{}{"namespaces": []interface{}{map[string]interface{}{"name": "files", "id": 1}}},
			wantErr: errLegacyNamespaces,
		},
		{
			name:    "case=an empty legacy list is not a legacy configuration",
			config:  map[string]interface{}{"namespaces": []interface{}{}},
			wantErr: errNoOPLConfigured,
		},
		{
			name:    "case=missing namespaces key",
			config:  map[string]interface{}{"limit": map[string]interface{}{}},
			wantErr: errNoOPLConfigured,
		},
		{
			name:    "case=null namespaces",
			config:  map[string]interface{}{"namespaces": nil},
			wantErr: errNoOPLConfigured,
		},
		{
			name:    "case=nil config",
			config:  nil,
			wantErr: errNoOPLConfigured,
		},
		{
			name:    "case=namespaces without a location",
			config:  map[string]interface{}{"namespaces": map[string]interface{}{}},
			wantErr: errNoOPLConfigured,
		},
		{
			name:    "case=empty location",
			config:  map[string]interface{}{"namespaces": map[string]interface{}{"location": ""}},
			wantErr: errNoOPLConfigured,
		},
		{
			// Anything that is neither an OPL pointer nor legacy definitions
			// carries no file to print, so it is reported as nothing configured.
			name:    "case=namespaces of an unexpected shape",
			config:  map[string]interface{}{"namespaces": "https://example.com/opl.bin"},
			wantErr: errNoOPLConfigured,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := oplLocation(tc.config)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestOPLLocationIsReadable pins the loader contract the command relies on:
// `ory update opl` writes the file as a base64:// payload, and Ory Network
// hands back an https:// location, so both have to resolve — while file://
// must not, since the location comes from the API. An error response must not
// be mistaken for the file either, or `ory get opl > opl.ts` writes the
// storage provider's error page to disk.
func TestOPLLocationIsReadable(t *testing.T) {
	const opl = "class Example implements Namespace {}"

	f := newOPLFetcher()

	t.Run("case=base64 payloads are read", func(t *testing.T) {
		location := "base64://" + base64.StdEncoding.EncodeToString([]byte(opl))

		got, err := f.FetchBytes(t.Context(), location)
		require.NoError(t, err)
		assert.Equal(t, opl, string(got))
	})

	t.Run("case=remote files are read", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(opl))
		}))
		t.Cleanup(s.Close)

		got, err := f.FetchBytes(t.Context(), s.URL+"/opl.bin")
		require.NoError(t, err)
		assert.Equal(t, opl, string(got))
	})

	t.Run("case=local files are refused", func(t *testing.T) {
		// The file has to exist, otherwise the test would pass even if the
		// file loader were enabled.
		path := filepath.Join(t.TempDir(), "opl.ts")
		require.NoError(t, os.WriteFile(path, []byte(opl), 0o600))

		got, err := f.FetchBytes(t.Context(), "file://"+path)
		require.ErrorIs(t, err, fetcher.ErrUnknownScheme, "a location from the API must never read from the local disk")
		assert.NotContains(t, string(got), opl)
	})

	t.Run("case=error responses are not mistaken for the file", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`<Error><Code>AccessDenied</Code></Error>`))
		}))
		t.Cleanup(s.Close)

		got, err := f.FetchBytes(t.Context(), s.URL+"/opl.bin")
		require.ErrorContains(t, err, "status code 200 but got 403")
		assert.Empty(t, got)
	})
}
