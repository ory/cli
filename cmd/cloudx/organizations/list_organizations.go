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
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			id, err := client.ProjectOrDefault(cmd, h)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			organizations, err := h.ListOrganizations(id)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintTable(cmd, &outputOrganizations{organizations})
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
