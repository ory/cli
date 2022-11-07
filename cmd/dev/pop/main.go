// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pop

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/pop/migration"
)

var Main = &cobra.Command{
	Use:   "pop",
	Short: "Helpers for pop",
}

func init() {
	Main.AddCommand(migration.Main)
}
