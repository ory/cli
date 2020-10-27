package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	BuildVersion = "master"
	BuildTime    = "undefined"
	BuildGitHash = "undefined"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this binary's version, build time, and git hash of this build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", BuildVersion)
		fmt.Printf("Git Hash:   %s\n", BuildGitHash)
		fmt.Printf("Build Time: %s\n", BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
