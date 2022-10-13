// Copyright © 2022 Ory Corp

package cmd

import (
	"fmt"

	"github.com/ory/cli/buildinfo"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this binary's version, build time, and git hash of this build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", buildinfo.Version)
		fmt.Printf("Git Hash:   %s\n", buildinfo.GitHash)
		fmt.Printf("Build Time: %s\n", buildinfo.Time)
	},
}
