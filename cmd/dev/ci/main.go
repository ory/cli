package ci

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/ci/github"
	"github.com/ory/cli/cmd/dev/ci/orbs"
)

var Main = &cobra.Command{
	Use:   "ci",
	Short: "Continuous Integration Helpers",
}

func init() {
	Main.AddCommand(orbs.Main)
	Main.AddCommand(github.Main)
}
