// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package swagger

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "swagger",
		Short: "Helpers for Swagger 2.0 / OpenAPI spec",
	}
	c.AddCommand(newSanitizeCmd())
	return c
}
