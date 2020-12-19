package identities

import (
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/kratos/cmd/identities"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
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
