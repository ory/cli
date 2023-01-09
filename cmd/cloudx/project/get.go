// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewGetProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Get the complete configuration of an Ory Network project.",
		Example: `$ ory get project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89

ID		ecaaa3cb-0730-4ee8-a6df-9553cdfeef89
SLUG	good-wright-t7kzy3vugf
STATE	running
NAME	Example Project

$ ory get project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format json

{
  "name": "Example Project",
  "identity": {
	"services": {
	  "config": {
		"courier": {
		  "smtp": {
			"from_name": "..."
		  }
		  // ...
		}
	  }
	}
  }
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			var id string
			if id, err = getSelectedProjectId(h, args); err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			project, err := h.GetProject(id)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintRow(cmd, (*outputProject)(project))
			return nil
		},
	}

	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
