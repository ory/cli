// Copyright © 2022 Ory Corp

package swagger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitize(t *testing.T) {
	fp := filepath.Join(os.TempDir(), uuid.New().String()+".json")

	require.NoError(t, sanitize("stub/in.json", fp))

	actual, err := os.ReadFile(fp)
	require.NoError(t, err)

	expected, err := os.ReadFile("stub/expected.json")
	require.NoError(t, err)

	require.NotEmpty(t, actual)
	require.NotEmpty(t, expected)

	assert.JSONEq(t, string(expected), string(actual), fp)
}
