package cmd

import (
	"fmt"

	"github.com/ory/cli/x"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this binary's version, build time, and git hash of this build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", x.BuildVersion)
		fmt.Printf("Git Hash:   %s\n", x.BuildGitHash)
		fmt.Printf("Build Time: %s\n", x.BuildTime)
	},
}
