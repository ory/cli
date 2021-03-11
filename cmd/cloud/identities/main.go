package identities

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/kratos/cmd/identities"
	"github.com/ory/x/cmdx"
)

var Main = &cobra.Command{
	Use:   "identities",
	Short: "Manage your identities",
}

func init() {
	cmdx.RegisterFormatFlags(identities.ListCmd.PersistentFlags())
	cmdx.RegisterFormatFlags(identities.GetCmd.PersistentFlags())
	remote.RegisterClientFlags(Main.PersistentFlags())

	Main.AddCommand(
		identities.DeleteCmd,
		identities.GetCmd,
		identities.ListCmd,
	)
}
