// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "github",
		Short: "Helpers for GitHub",
	}
	c.AddCommand(newEnvCmd())
	return c
}
