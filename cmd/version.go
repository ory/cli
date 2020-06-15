package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "master"
	Date    = "undefined"
	Commit  = "undefined"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this binary's version, build time, and git hash of this build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git Hash:   %s\n", Commit)
		fmt.Printf("Build Time: %s\n", Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
