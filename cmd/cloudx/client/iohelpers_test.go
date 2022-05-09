package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/assertx"
)

func TestReadConfigFiles(t *testing.T) {
	configs, err := ReadConfigFiles([]string{
		"fixtures/a.yaml",
		"fixtures/b.yml",
		"fixtures/c.json",
	})
	require.NoError(t, err)
	assertx.EqualAsJSON(t, json.RawMessage(`[{"a":true},{"b":true},{"c":true}]`), configs)
}
