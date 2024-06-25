// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"workspace", "ws"},
		Short:   "List Ory Network workspaces",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			workspaces, err := h.ListWorkspaces(cmd.Context())
			if err != nil {
				return err
			}

			cmdx.PrintTable(cmd, (outputWorkspaces)(workspaces))
			return nil
		},
	}
}
