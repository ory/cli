// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openapi

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "openapi",
	Short: "Helpers for OpenAPI 3.0",
}

func init() {
	Main.AddCommand(
		migrateCmd,
	)
}
