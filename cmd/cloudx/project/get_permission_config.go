// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewGetKetoConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "permission-config",
		Aliases: []string{"pc", "keto-config"},
		Args:    cobra.NoArgs,
		Short:   "Get Ory Permissions configuration.",
		Long:    "Get the Ory Permissions configuration for an Ory Network project.",
		Example: `$ ory get permission-config --project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format yaml > permission-config.yaml

$ ory get permission-config --format json   # uses currently selected project

{
  "namespaces": [
    {
      "name": "files",
      "id": 1
	},1
    // ...
  ]
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

			cmdx.PrintJSONAble(cmd, outputConfig(project.Services.Permission.Config))
			return nil
		},
	}

	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	return cmd
}
