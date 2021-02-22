package protect

import (
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "secure",
	Short: "Secure HTTP(s) Endpoints",
}

func init() {
	remote.RegisterClientFlags(Main.PersistentFlags())

	Main.AddCommand(
		ProxyCmd,
	)
}
