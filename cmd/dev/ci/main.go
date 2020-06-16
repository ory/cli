package ci

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/ci/orbs"
)

var Main = &cobra.Command{
	Use:   "ci",
	Short: "Helpers for CircleCI",
}

func init() {
	Main.AddCommand(orbs.Main)
}
