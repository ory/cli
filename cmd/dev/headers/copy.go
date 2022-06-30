// Tool for copying files inside https://github.com/ory/meta/blob/master/scripts/sync.sh
// and adding a link to the original as a header comment.

package headers

import (
	"fmt"

	"github.com/spf13/cobra"
)

// template for the header
const LINK_TEMPLATE = "AUTO-GENERATED, DO NOT EDIT! Please edit the original at %s."

// the token that identifies comments containing the license
const LINK_TOKEN = "Copyright ©"

var copy = &cobra.Command{
	Use:   "cp",
	Short: "Copies the given files and adds a header pointing to their original",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("HELLO %v\n", args)
		return nil
	},
}

func init() {
	Main.AddCommand(copyright)
	copy.Flags().StringVarP(&glob, "files", "f", "", "Commands to be run if current component is affected.")
	copy.Flags().BoolVarP(&recursive, "recursive", "R", false, "Whether to copy files in subdirectories")
}

var glob string
var recursive bool
