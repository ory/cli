package cmd

import (
	"github.com/ory/cli/cmd/cloud/identities"
)

func init() {
	rootCmd.AddCommand(
		identities.Main,
	)
}
