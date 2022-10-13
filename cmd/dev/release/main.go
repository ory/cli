// Copyright © 2022 Ory Corp

package release

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/release/notify"
)

var Main = &cobra.Command{
	Use:   "release",
	Short: "Release infrastructure for ORY and related components",
}

func init() {
	Main.AddCommand(
		notify.Main,
	)
}
