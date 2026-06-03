// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package orbs

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "orbs",
		Short: "Helpers for CircleCI",
	}
	c.AddCommand(newBumpCmd())
	return c
}
