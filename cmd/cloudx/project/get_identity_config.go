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
		Use:     "identity-config [project-id]",
		Aliases: []string{"ic", "kratos-config"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Ory Identities configuration.",
		Long:    "Get the Ory Identities configuration for the specified Ory Network project.",
		Example: `$ ory get identity-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format yaml > identity-config.yaml

$ ory get identity-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format json

{
  "selfservice": {
	"methods": {
	  "password": { "enabled": false }
	}
	// ...
  }
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			id, err := selectedProjectID(h, args)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			project, err := h.GetProject(id)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintJSONAble(cmd, outputConfig(project.Services.Identity.Config))
			return nil
		},
	}

	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	return cmd
}
