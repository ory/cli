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
			opts := make([]client.CommandHelperOption, 0, 1)
			if len(args) == 1 {
				opts = append(opts, client.WithProjectOverride(args[0]))
			}
			h, err := client.NewCobraCommandHelper(cmd, opts...)
			if err != nil {
				return err
			}

			id, err := h.ProjectID()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			if err := h.SelectProject(id); err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintRow(cmd, &selectedProject{ID: id})
			return nil
		},
	}

	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
