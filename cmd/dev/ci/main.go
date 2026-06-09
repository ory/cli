// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package ci

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/ci/deps"
	"github.com/ory/cli/cmd/dev/ci/github"
	"github.com/ory/cli/cmd/dev/ci/orbs"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "ci",
		Short: "Continuous Integration helpers",
	}
	c.AddCommand(
		orbs.NewCommand(),
		github.NewCommand(),
		deps.NewCommand(),
	)
	return c
}
