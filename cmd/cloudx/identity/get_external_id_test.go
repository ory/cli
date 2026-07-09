// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/randx"
)

func TestGetIdentityByExternalID(t *testing.T) {
	t.Parallel()

	// The external_id column is VARCHAR(64), so keep this well under that limit.
	externalID := "customer-" + randx.MustString(16, randx.AlphaLowerNum)
	email := testhelpers.FakeEmail()

	importPath := filepath.Join(t.TempDir(), "import.json")
	require.NoError(t, os.WriteFile(importPath, []byte(`{
  "schema_id": "preset://username",
  "traits": {
    "username": "`+email+`"
  },
  "external_id": "`+externalID+`"
}`), 0o600))

	stdout, stderr, err := defaultCmd.Exec(nil, "import", "identities", "--format", "json", "--project", defaultProject.Id, importPath)
	require.NoError(t, err, stderr)
	imported := gjson.Parse(stdout)
	require.True(t, gjson.Valid(stdout))
	userID := imported.Get("id").String()
	require.NotEmpty(t, userID)
	// the import must persist the external ID
	assert.Equal(t, externalID, imported.Get("external_id").String(), stdout)

	t.Run("is able to get identity by external ID", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "get", "identity", "--external-id", "--format", "json", "--project", defaultProject.Id, externalID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, userID, out.Get("id").String(), stdout)
		assert.Equal(t, externalID, out.Get("external_id").String(), stdout)
	})

	t.Run("fails to get identity by unknown external ID", func(t *testing.T) {
		_, _, err := defaultCmd.Exec(nil, "get", "identity", "--external-id", "--project", defaultProject.Id, "unknown-"+randx.MustString(16, randx.AlphaLowerNum))
		require.Error(t, err)
	})
}
