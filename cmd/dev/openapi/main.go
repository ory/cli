// Copyright © 2022 Ory Corp

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
