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
		Use:     "permission-config [project-id]",
		Aliases: []string{"pc", "keto-config"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Ory Permissions configuration.",
		Long:    "Get the Ory Permissions configuration for the specified Ory Network project.",
		Example: `$ ory get permission-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format yaml > permission-config.yaml

$ ory get permission-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format json

{
  "namespaces": [
    {
      "name": "files",
      "id": 1
	},1
    // ...
  ]
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			id, err := getSelectedProjectId(h, args)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			project, err := h.GetProject(id)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintJSONAble(cmd, outputConfig(project.Services.Permission.Config))
			return nil
		},
	}

	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	return cmd
}
