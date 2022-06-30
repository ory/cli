package headers

import (
	"fmt"

	"github.com/spf13/cobra"
)

var copy = &cobra.Command{
	Use:   "add",
	Short: "Adds the given header to the given files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("HELLO %v\n", args)
		return nil
	},
}

func init() {
	Main.AddCommand(copyright)
	copy.Flags().StringVarP(&glob, "files", "f", "", "Commands to be run if current component is affected.")
}

var glob string
