package dev

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/ci"
	"github.com/ory/cli/cmd/dev/markdown"
	"github.com/ory/cli/cmd/dev/newsletter"
	"github.com/ory/cli/cmd/dev/pop"
	"github.com/ory/cli/cmd/dev/release"
	"github.com/ory/cli/cmd/dev/swagger"
)

var Main = &cobra.Command{
	Use:   "dev",
	Short: "Tools for developing ORY technology",
}

func init() {
	Main.AddCommand(
		pop.Main,
		newsletter.Main,
		markdown.Main,
		release.Main,
		swagger.Main,
		ci.Main,
	)
}
