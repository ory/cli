// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package organizations

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewListOrganizationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organizations",
		Args:  cobra.NoArgs,
		Short: "List your Ory Network organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := h.ProjectID()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			organizations, err := h.ListOrganizations(cmd.Context(), projectID)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintTable(cmd, &outputOrganizations{organizations})
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
