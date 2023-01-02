// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

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
