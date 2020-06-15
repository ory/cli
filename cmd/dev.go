package cmd

import (
	"github.com/ory/cli/cmd/dev"
)

func init() {
	rootCmd.AddCommand(dev.Main)
}
