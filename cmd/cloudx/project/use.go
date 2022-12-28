// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewUseProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Set the project as the default. When no id is provided, prints the currently used default project.",
		Example: `$ ory use project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89

ID		ecaaa3cb-0730-4ee8-a6df-9553cdfeef89

$ ory use project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format json

{
  "id": "ecaaa3cb-0730-4ee8-a6df-9553cdfeef89
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			var id string
			if len(args) == 0 {
				if id, err = h.GetDefaultProjectID(); err != nil {
					return cmdx.PrintOpenAPIError(cmd, err)
				}
			} else {
				id = args[0]
				err = h.SetDefaultProject(id)
				if err != nil {
					return cmdx.PrintOpenAPIError(cmd, err)
				}
			}

			cmdx.PrintRow(cmd, &selectedProject{Id: id})
			return nil
		},
	}

	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
