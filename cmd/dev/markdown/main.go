// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package markdown

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "markdown",
		Short: "Utilities for working with markdown",
	}
	c.AddCommand(newRenderCmd())
	return c
}
