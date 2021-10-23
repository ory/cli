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
	listCmd := identities.NewListCmd()
	getCmd := identities.NewGetCmd()
	deleteCmd := identities.NewDeleteCmd()

	cmdx.RegisterFormatFlags(listCmd.PersistentFlags())
	cmdx.RegisterFormatFlags(getCmd.PersistentFlags())
	remote.RegisterClientFlags(Main.PersistentFlags())

	Main.AddCommand(
		listCmd,
		getCmd,
		deleteCmd,
	)
}
