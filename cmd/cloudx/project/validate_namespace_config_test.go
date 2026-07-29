// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/cmdx"
)

const (
	validOPL   = `class Default implements Namespace {}`
	invalidOPL = `class Default implements Namespace {`
)

func TestValidateNamespaceConfig(t *testing.T) {
	if testing.Short() {
		// this test needs internet, typically not available when you're on a (german) train
		return
	}

	t.Parallel()

	t.Run("accepts a valid file", func(t *testing.T) {
		t.Parallel()

		config := writeFile(t, validOPL)
		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(nil, "validate", "opl", "--file", config)
			require.NoError(t, err, stderr)
			assert.Contains(t, stdout, "valid")
			assert.Empty(t, stderr)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
		runWithProjectAsFlag(ctx, t, extraProject.Id, validate)
	})

	t.Run("reports syntax errors and exits non-zero", func(t *testing.T) {
		t.Parallel()

		config := writeFile(t, invalidOPL)
		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(nil, "validate", "opl", "--file", config)
			assert.ErrorIs(t, err, cmdx.ErrNoPrintButFail, stderr)
			assert.Empty(t, stdout)
			// the location is rendered as "<source>:<line>:<column>: <message>"
			assert.Contains(t, stderr, config+":", stderr)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
		runWithProjectAsFlag(ctx, t, extraProject.Id, validate)
	})

	t.Run("reports syntax errors as JSON", func(t *testing.T) {
		t.Parallel()

		config := writeFile(t, invalidOPL)
		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(nil, "validate", "opl", "--file", config, "--format", "json")
			assert.ErrorIs(t, err, cmdx.ErrNoPrintButFail, stderr)

			errs := gjson.Get(stdout, "errors")
			require.True(t, errs.IsArray(), stdout)
			require.NotEmpty(t, errs.Array(), stdout)
			assert.NotEmpty(t, errs.Array()[0].Get("message").String(), stdout)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
		runWithProjectAsFlag(ctx, t, extraProject.Id, validate)
	})

	t.Run("reads the file from stdin", func(t *testing.T) {
		t.Parallel()

		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(strings.NewReader(validOPL), "validate", "opl", "--file", "-")
			require.NoError(t, err, stderr)
			assert.Contains(t, stdout, "valid")
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
		runWithProjectAsFlag(ctx, t, extraProject.Id, validate)
	})

	t.Run("prints nothing on success when quiet", func(t *testing.T) {
		t.Parallel()

		config := writeFile(t, validOPL)
		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(nil, "validate", "opl", "--file", config, "--quiet")
			require.NoError(t, err, stderr)
			assert.Empty(t, stdout)
			assert.Empty(t, stderr)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
	})
}
