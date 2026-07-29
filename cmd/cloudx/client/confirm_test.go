// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// refusingReader fails the test if anything tries to read from it. Reading
// stdin is exactly what --yes and --quiet have to avoid: in CI there is nobody
// to answer, so the command would block forever.
type refusingReader struct{ t *testing.T }

func (r refusingReader) Read([]byte) (int, error) {
	r.t.Error("stdin must not be read without an interactive confirmation")
	return 0, io.EOF
}

// TestConfirm covers https://github.com/ory/cli/issues/219: --yes was parsed
// and stored but never read, so it did not skip any prompt.
func TestConfirm(t *testing.T) {
	t.Run("case=--yes answers yes without prompting", func(t *testing.T) {
		var out bytes.Buffer
		h := &CommandHelper{
			noConfirm:        true,
			Stdin:            bufio.NewReader(refusingReader{t}),
			VerboseErrWriter: &out,
		}

		ok, err := h.Confirm("do you want to sign in?")
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Empty(t, out.String(), "no prompt should be shown when --yes is set")
	})

	for _, tc := range []struct {
		name  string
		input string
		want  bool
	}{
		{name: "case=prompts and accepts yes", input: "y\n", want: true},
		{name: "case=prompts and accepts no", input: "n\n", want: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			h := &CommandHelper{
				noConfirm:        false,
				Stdin:            bufio.NewReader(strings.NewReader(tc.input)),
				VerboseErrWriter: &out,
			}

			ok, err := h.Confirm("do you want to sign in?")
			require.NoError(t, err)
			assert.Equal(t, tc.want, ok)
			assert.Contains(t, out.String(), "do you want to sign in?")
		})
	}

	t.Run("case=--quiet does not block on stdin", func(t *testing.T) {
		h := &CommandHelper{
			isQuiet:          true,
			Stdin:            bufio.NewReader(refusingReader{t}),
			VerboseErrWriter: io.Discard,
		}

		ok, err := h.Confirm("do you want to sign in?")
		assert.Error(t, err)
		assert.False(t, ok)
	})
}

// TestTemporaryAPIKeyConfirmation guards the wiring of the confirmation: the
// prompt must actually be reachable. Gating on Authenticate instead of
// checkAuthenticated performs the sign-in instead of checking for it, which
// makes the prompt (and with it --yes) dead code.
func TestTemporaryAPIKeyConfirmation(t *testing.T) {
	errBrowserOpened := errors.New("browser opened")

	newHelper := func(t *testing.T, stdin io.Reader, out io.Writer, opts ...CommandHelperOption) *CommandHelper {
		h, err := NewCommandHelper(context.Background(), append([]CommandHelperOption{
			WithConfigLocation(filepath.Join(t.TempDir(), "config.json")),
			WithStdin(stdin),
			WithVerboseErrWriter(out),
			WithOpenBrowserHook(func(string) error { return errBrowserOpened }),
		}, opts...)...)
		require.NoError(t, err)
		return h
	}

	t.Run("case=asks before signing in", func(t *testing.T) {
		var out bytes.Buffer
		h := newHelper(t, strings.NewReader("n\n"), &out)

		key, cleanup, err := h.TemporaryAPIKey(context.Background(), "test", 0)
		require.NoError(t, err)
		require.NoError(t, cleanup())
		assert.Empty(t, key)
		assert.Contains(t, out.String(), "Do you want to sign in?")
	})

	t.Run("case=--yes skips the prompt and signs in", func(t *testing.T) {
		var out bytes.Buffer
		h := newHelper(t, refusingReader{t}, &out, WithNoConfirm(true))

		_, _, err := h.TemporaryAPIKey(context.Background(), "test", 0)
		assert.ErrorIs(t, err, errBrowserOpened)
	})

	t.Run("case=--quiet continues without an API key", func(t *testing.T) {
		var out bytes.Buffer
		h := newHelper(t, refusingReader{t}, &out, WithQuiet(true))

		key, cleanup, err := h.TemporaryAPIKey(context.Background(), "test", 0)
		require.NoError(t, err)
		require.NoError(t, cleanup())
		assert.Empty(t, key)
		assert.Contains(t, out.String(), "Remove the `--quiet` flag")
	})
}
