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
			// not just "valid", which is also a substring of "invalid"
			assert.Contains(t, stdout, "file is valid")
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

	t.Run("reports an empty error list as JSON", func(t *testing.T) {
		t.Parallel()

		config := writeFile(t, validOPL)
		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(nil, "validate", "opl", "--file", config, "--format", "json")
			require.NoError(t, err, stderr)

			// the key is always present, so that `jq '.errors | length'` works
			errs := gjson.Get(stdout, "errors")
			require.True(t, errs.IsArray(), stdout)
			assert.Empty(t, errs.Array(), stdout)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
	})

	t.Run("accepts an empty file", func(t *testing.T) {
		t.Parallel()

		config := writeFile(t, "")
		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(nil, "validate", "opl", "--file", config)
			require.NoError(t, err, stderr)
			assert.Contains(t, stdout, "file is valid")
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
	})

	t.Run("reads the file from stdin", func(t *testing.T) {
		t.Parallel()

		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(strings.NewReader(validOPL), "validate", "opl", "--file", "-")
			require.NoError(t, err, stderr)
			assert.Contains(t, stdout, "file is valid")

			// errors are reported against "stdin", not against the literal "-"
			_, stderr, err = exec(strings.NewReader(invalidOPL), "validate", "opl", "--file", "-")
			assert.ErrorIs(t, err, cmdx.ErrNoPrintButFail, stderr)
			assert.Contains(t, stderr, "stdin:", stderr)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
		runWithProjectAsFlag(ctx, t, extraProject.Id, validate)
	})

	t.Run("stays quiet on success but still reports syntax errors", func(t *testing.T) {
		t.Parallel()

		valid, invalid := writeFile(t, validOPL), writeFile(t, invalidOPL)
		validate := func(t *testing.T, exec execFunc) {
			stdout, stderr, err := exec(nil, "validate", "opl", "--file", valid, "--quiet")
			require.NoError(t, err, stderr)
			assert.Empty(t, stdout)
			assert.Empty(t, stderr)

			// --quiet silences the success message, not the syntax errors
			stdout, stderr, err = exec(nil, "validate", "opl", "--file", invalid, "--quiet")
			assert.ErrorIs(t, err, cmdx.ErrNoPrintButFail, stderr)
			assert.Empty(t, stdout)
			assert.Contains(t, stderr, invalid+":", stderr)
		}

		runWithProjectAsDefault(ctx, t, defaultProject.Id, validate)
		runWithProjectAsFlag(ctx, t, extraProject.Id, validate)
	})
}
