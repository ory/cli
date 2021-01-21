package deps

import (
	"github.com/spf13/cobra"
)
var pOs string
var pArch string
var pConfig string

var Main = &cobra.Command{
	Use:   "deps",
	Short: "Helpers for binary dependencies in Makefiles.",
}

func init() {
	Main.PersistentFlags().StringVarP(&pOs, "os", "o", "","OS the binary should run on. Currently only 'linux' and 'darwin' are supported.")
	Main.PersistentFlags().StringVarP(&pArch, "architecture", "a", "", "Architecture the binary should run on. Currently only 'amd64' is supported.")
	Main.PersistentFlags().StringVarP(&pConfig, "config", "c", "", "Path to config files.")
}
