package monorepo

import (
	"github.com/spf13/cobra"
)

var rootDirectory string
var changedComponentIds string
var verbose bool
var debug bool

var Main = &cobra.Command{
	Use:   "monorepo",
	Short: "Helpers for CircleCI Monorepo Support",
}

func init() {
	Main.PersistentFlags().StringVarP(&rootDirectory, "root", "r", ".", "Root directory to be used to traverse and search for dependency configurations.")
	Main.PersistentFlags().StringVarP(&changedComponentIds, "changed", "c", "", "Changed Components IDs.")
	Main.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	Main.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug output")
}
