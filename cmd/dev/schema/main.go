// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "schema",
		Short: "JSON Schema related helpers",
	}
	c.AddCommand(newRenderVersionCmd())
	return c
}
