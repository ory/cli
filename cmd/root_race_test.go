// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"sync"
	"testing"
)

// TestNewRootCmdConcurrent guards against data races in NewRootCmd. The command
// executer (github.com/ory/x/cmdx.CommandExecuter) calls New() for every command
// invocation, and the cloudx test suites run those invocations in parallel, so
// NewRootCmd must be safe to call concurrently and must return fully independent
// command trees. Run with -race to detect regressions.
func TestNewRootCmdConcurrent(t *testing.T) {
	const goroutines = 16

	// Mix of args so we exercise both command execution and full-tree help
	// rendering, which walks every (shared) subcommand.
	argSets := [][]string{{"version"}, {"--help"}, {"dev", "--help"}}

	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Go(func() {
			args := argSets[i%len(argSets)]
			c := NewRootCmd()
			c.SetArgs(args)
			c.SetOut(&bytes.Buffer{})
			c.SetErr(&bytes.Buffer{})
			_ = c.Execute()
		})
	}
	wg.Wait()
}
