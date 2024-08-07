// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewListProjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List your Ory Network projects.",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			projects, err := h.ListProjects(cmd.Context(), h.WorkspaceID())
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintTable(cmd, &outputProjectCollection{projects})
			return nil
		},
	}
	client.RegisterWorkspaceFlag(cmd.Flags())

	return cmd
}
