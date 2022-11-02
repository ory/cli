// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package monorepo

import "github.com/spf13/cobra"

var rootDirectory string
var verbose bool
var debug bool
var pr string
var revisionRange string
var branch string

// Main cobra command for monorepo support
var Main = &cobra.Command{
	Use:   "monorepo",
	Short: "Helpers for CircleCI monorepo support",
}

func isPR() bool {
	return len(pr) > 0
}

func init() {
	Main.PersistentFlags().StringVarP(&rootDirectory, "root", "r", ".", "Root directory to be used to traverse and search for dependency configurations.")
	Main.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	Main.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug output")
	Main.PersistentFlags().StringVar(&pr, "pr", "", "Pull Request")
	Main.PersistentFlags().StringVar(&branch, "branch", "", "Branch")
	Main.PersistentFlags().StringVar(&revisionRange, "revisionRange", "", "RevisionRange used to determine changes!")
}
