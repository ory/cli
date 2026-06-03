// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openapi

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "openapi",
		Short: "Helpers for OpenAPI 3.0",
	}
	c.AddCommand(newMigrateCmd())
	return c
}
