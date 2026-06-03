// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package monorepo

import "github.com/spf13/cobra"

// Package-level scratch variables read by the helpers in this package. They are
// populated from per-command local flag values in PersistentPreRun (see
// NewCommand) and the leaves' RunE so that constructing the command tree
// concurrently does not race on flag default writes.
var (
	rootDirectory string
	verbose       bool
	debug         bool
	pr            string
	revisionRange string
	branch        string
)

func isPR() bool {
	return len(pr) > 0
}

// NewCommand returns a fresh `monorepo` command tree.
func NewCommand() *cobra.Command {
	var (
		localRoot, localPR, localBranch, localRev string
		localVerbose, localDebug                  bool
	)
	c := &cobra.Command{
		Use:   "monorepo",
		Short: "Helpers for CircleCI monorepo support",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			rootDirectory = localRoot
			verbose = localVerbose
			debug = localDebug
			pr = localPR
			branch = localBranch
			revisionRange = localRev
		},
	}
	c.PersistentFlags().StringVarP(&localRoot, "root", "r", ".", "Root directory to be used to traverse and search for dependency configurations.")
	c.PersistentFlags().BoolVarP(&localVerbose, "verbose", "v", false, "Verbose output")
	c.PersistentFlags().BoolVarP(&localDebug, "debug", "d", false, "Debug output")
	c.PersistentFlags().StringVar(&localPR, "pr", "", "Pull Request")
	c.PersistentFlags().StringVar(&localBranch, "branch", "", "Branch")
	c.PersistentFlags().StringVar(&localRev, "revisionRange", "", "RevisionRange used to determine changes!")
	c.AddCommand(newRunCmd(), newChangesCmd(), newComponentsCmd())
	return c
}
