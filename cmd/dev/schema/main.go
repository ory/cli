// Copyright © 2022 Ory Corp

package schema

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "schema",
	Short: "JSON Schema related helpers",
}

func init() {
	Main.AddCommand(
		RenderVersion,
	)
}
