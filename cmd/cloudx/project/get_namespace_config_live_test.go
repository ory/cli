// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetOPL covers the round-trip #321 asks for, in both directions:
//
//	ory update opl --file namespace_config.ts   # what `get opl` must return
//	ory get opl > namespace_config.ts           # must be usable as input again
//
// The second direction is the one that matters for checking the permission
// model into version control: if the two commands disagreed on the file format,
// committing `get opl` output and applying it later would corrupt the model
// rather than restore it.
func TestGetOPL(t *testing.T) {
	if testing.Short() {
		// this test needs internet, typically not available when you're on a (german) train
		t.Skip("skipping test that requires internet access")
	}

	t.Parallel()

	content := `class Default implements Namespace {}`
	config := writeFile(t, content)

	// The read-after-write assertions below would be racy against the shared
	// fixtures: other tests in this package rewrite their permission config in
	// parallel.
	project := newProject(t)

	_, stderr, err := defaultCmd.Exec(nil, "update", "opl", "--project", project, "--file", config)
	require.NoError(t, err, stderr)

	opl, stderr, err := defaultCmd.Exec(nil, "get", "opl", "--project", project)
	require.NoError(t, err, stderr)
	require.Equal(t, content, opl, "`get opl` must return what `update opl` uploaded")

	// Feed `get opl` output straight back in, exactly as a user committing the
	// file and re-applying it would, and read it once more. Writing it to disk
	// verbatim also catches any stray wrapping or trailing newline that would
	// accumulate over repeated round-trips.
	_, stderr, err = defaultCmd.Exec(nil, "update", "opl", "--project", project, "--file", writeFile(t, opl))
	require.NoError(t, err, stderr, "`get opl` output must be accepted by `update opl`")

	opl, stderr, err = defaultCmd.Exec(nil, "get", "opl", "--project", project)
	require.NoError(t, err, stderr)
	assert.Equal(t, content, opl, "the OPL file must survive a full get/update round-trip unchanged")
}
