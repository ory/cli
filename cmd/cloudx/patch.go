// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/project"
)

func NewPatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patch",
		Short: "Patch resources",
	}
	client.RegisterConfigFlag(cmd.PersistentFlags())
	cmd.AddCommand(
		project.NewProjectsPatchCmd(),
		project.NewPatchKratosConfigCmd(),
		project.NewPatchKetoConfigCmd(),
		project.NewPatchOAuth2ConfigCmd(),
		project.NewUpdateNamespaceConfigCmd(),
	)
	return cmd
}
