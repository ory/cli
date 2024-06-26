// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewGetKratosConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "identity-config",
		Aliases: []string{"ic", "kratos-config"},
		Args:    cobra.NoArgs,
		Short:   "Get Ory Identities configuration.",
		Long:    "Get the Ory Identities configuration for an Ory Network project.",
		Example: `$ ory get identity-config --project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format yaml > identity-config.yaml

$ ory get identity-config --format json   # uses currently selected project

{
  "selfservice": {
	"methods": {
	  "password": { "enabled": false }
	}
	// ...
  }
}`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			project, err := h.GetSelectedProject(cmd.Context())
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintJSONAble(cmd, outputConfig(project.Services.Identity.Config))
			return nil
		},
	}

	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	return cmd
}
