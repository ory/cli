// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package deps

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "deps",
		Short: "Helpers for binary dependencies in Makefiles.",
	}
	c.PersistentFlags().StringP("os", "o", "", "OS the binary should run on. Currently only 'linux' and 'darwin' are supported.")
	c.PersistentFlags().StringP("architecture", "a", "", "Architecture the binary should run on. Currently only 'amd64' is supported.")
	c.PersistentFlags().StringP("config", "c", "", "Path to config files.")
	c.AddCommand(newURLCmd())
	return c
}
