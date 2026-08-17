// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRedactInZip pins the guarantee the Playwright traces depend on: this
// repository is public and CI uploads the traces as an artifact, so the
// rate-limit header value must not survive anywhere inside one.
func TestRedactInZip(t *testing.T) {
	const secret = `s3cret"value\with-escapes`

	writeArchive := func(t *testing.T, entries map[string]string) string {
		path := filepath.Join(t.TempDir(), "trace.zip")
		f, err := os.Create(path)
		require.NoError(t, err)
		defer f.Close()

		w := zip.NewWriter(f)
		for name, content := range entries {
			e, err := w.Create(name)
			require.NoError(t, err)
			_, err = e.Write([]byte(content))
			require.NoError(t, err)
		}
		require.NoError(t, w.Close())
		return path
	}

	readArchive := func(t *testing.T, path string) map[string]string {
		r, err := zip.OpenReader(path)
		require.NoError(t, err)
		defer r.Close()

		out := make(map[string]string, len(r.File))
		for _, f := range r.File {
			src, err := f.Open()
			require.NoError(t, err)
			content, err := io.ReadAll(src)
			require.NoError(t, err)
			require.NoError(t, src.Close())
			out[f.Name] = string(content)
		}
		return out
	}

	// The trace stores headers as JSON string values, so the secret appears in
	// its escaped spelling rather than verbatim.
	escaped, err := json.Marshal(secret)
	require.NoError(t, err)
	jsonEncoded := string(escaped)

	t.Run("case=removes the secret in both spellings", func(t *testing.T) {
		path := writeArchive(t, map[string]string{
			"trace.network": `{"headers":[{"name":"Ory-RateLimit-Action","value":` + jsonEncoded + `}]}`,
			"trace.trace":   "prefix " + secret + " suffix",
			"resources/1":   "a response body mentioning nothing",
		})

		require.NoError(t, redactInZip(path, secret))

		for name, content := range readArchive(t, path) {
			assert.NotContains(t, content, secret, "%s still holds the raw secret", name)
			assert.NotContains(t, content, jsonEncoded[1:len(jsonEncoded)-1], "%s still holds the escaped secret", name)
		}
	})

	t.Run("case=leaves the rest of the trace intact", func(t *testing.T) {
		path := writeArchive(t, map[string]string{
			"trace.trace": "keep me " + secret + " keep me too",
			"resources/1": "untouched",
		})

		require.NoError(t, redactInZip(path, secret))

		got := readArchive(t, path)
		assert.Equal(t, "keep me "+redactedPlaceholder+" keep me too", got["trace.trace"])
		assert.Equal(t, "untouched", got["resources/1"])
	})

	t.Run("case=no configured secret leaves the archive alone", func(t *testing.T) {
		path := writeArchive(t, map[string]string{"trace.trace": "verbatim"})

		require.NoError(t, redactInZip(path, ""))

		assert.Equal(t, "verbatim", readArchive(t, path)["trace.trace"])
	})

	t.Run("case=an unreadable archive is an error, so the caller can delete it", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "not-a-zip.zip")
		require.NoError(t, os.WriteFile(path, []byte("definitely not a zip"), 0o600))

		assert.Error(t, redactInZip(path, secret))
	})
}
