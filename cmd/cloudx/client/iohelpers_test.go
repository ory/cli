// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/assertx"
)

func TestReadConfigFiles(t *testing.T) {
	configs, err := ReadConfigFiles([]string{
		"fixtures/iohelpers/a.yaml",
		"fixtures/iohelpers/b.yml",
		"fixtures/iohelpers/c.json",
	})
	require.NoError(t, err)
	assertx.EqualAsJSON(t, json.RawMessage(`[{"a":true},{"b":true},{"c":true}]`), configs)
}
