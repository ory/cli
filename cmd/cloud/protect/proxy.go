package protect

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var ProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Secure Endpoint Using the Ory Reverse Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is yet to be implemented.")
		os.Exit(1)
	},
}
