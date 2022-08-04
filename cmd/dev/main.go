package dev

import (
	"github.com/ory/cli/cmd/dev/ci"
	"github.com/ory/cli/cmd/dev/headers"
	"github.com/ory/cli/cmd/dev/markdown"
	"github.com/ory/cli/cmd/dev/newsletter"
	"github.com/ory/cli/cmd/dev/openapi"
	"github.com/ory/cli/cmd/dev/pop"
	"github.com/ory/cli/cmd/dev/release"
	"github.com/ory/cli/cmd/dev/schema"
	"github.com/ory/cli/cmd/dev/swagger"
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "dev",
	Short: "Developer tools for writing Ory software",
	Long: `This section of the Ory CLI is full of convenience functions needed when contributing software to Ory.
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
