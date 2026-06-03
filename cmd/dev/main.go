// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/ci"
	"github.com/ory/cli/cmd/dev/headers"
	"github.com/ory/cli/cmd/dev/markdown"
	"github.com/ory/cli/cmd/dev/newsletter"
	"github.com/ory/cli/cmd/dev/openapi"
	"github.com/ory/cli/cmd/dev/release"
	"github.com/ory/cli/cmd/dev/schema"
	"github.com/ory/cli/cmd/dev/swagger"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "dev",
		Short: "Developer tools for writing Ory software",
		Long: `Developer tools and convenience functions for writing Ory software.
Please check the individual commands for more information!`,
		Hidden: true,
	}
	c.AddCommand(
		newsletter.NewCommand(),
		markdown.NewCommand(),
		release.NewCommand(),
		swagger.NewCommand(),
		ci.NewCommand(),
		schema.NewCommand(),
		openapi.NewCommand(),
		headers.NewCommand(),
	)
	return c
}
