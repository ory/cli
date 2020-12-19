package cmd

import (
	"github.com/ory/cli/cmd/cloud"
)

func init() {
	rootCmd.AddCommand(cloud.Main)
}
