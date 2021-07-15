package proxy

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "proxy",
	Short: "Easily protect applications with the Ory Proxy",
}

func init() {
	Main.AddCommand(
		NewProxyLocalCmd(),
		//NewProxyProductionCmd(),
	)
}
