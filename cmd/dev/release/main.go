// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/release/notify"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "release",
		Short: "Release infrastructure for ORY and related components",
	}
	c.AddCommand(
		notify.NewCommand(),
		newCompileCmd(),
		newPublishCmd(),
	)
	return c
}
