// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// refusingReader fails the test if anything tries to read from it. Reading
// stdin is exactly what --yes has to avoid: in CI there is nobody to answer,
// so the command would block forever.
type refusingReader struct{ t *testing.T }

func (r refusingReader) Read([]byte) (int, error) {
	r.t.Error("stdin must not be read when --yes is set")
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
}
