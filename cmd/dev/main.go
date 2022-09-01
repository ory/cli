package dev

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/ci"
	"github.com/ory/cli/cmd/dev/headers"
	"github.com/ory/cli/cmd/dev/markdown"
	"github.com/ory/cli/cmd/dev/newsletter"
	"github.com/ory/cli/cmd/dev/openapi"
	"github.com/ory/cli/cmd/dev/pop"
	"github.com/ory/cli/cmd/dev/release"
	"github.com/ory/cli/cmd/dev/schema"
	"github.com/ory/cli/cmd/dev/swagger"
)

var Main = &cobra.Command{
	Use:   "dev",
	Short: "Developer tools for writing Ory software",
	Long: `Developer tools and convenience functions for writing Ory software.
Please check the individual commands for more information!`,
	Hidden: true,
}

func init() {
	Main.AddCommand(
		pop.Main,
		newsletter.Main,
		markdown.Main,
		release.Main,
		swagger.Main,
		ci.Main,
		schema.Main,
		openapi.Main,
		headers.Main,
	)
}
