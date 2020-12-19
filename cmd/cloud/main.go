package cloud

import (
	"github.com/ory/cli/cmd/cloud/identities"
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "cloud",
	Short: "Manage your ORY Cloud projects.",
}

func init() {
	Main.AddCommand(
		identities.Main,
	)
}
