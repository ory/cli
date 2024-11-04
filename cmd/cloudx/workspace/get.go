// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workspace <id>",
		Aliases: []string{"workspaces", "ws"},
		Short:   "Get an Ory Network workspaces",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			workspace, err := h.GetWorkspace(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			cmdx.PrintRow(cmd, (*outputWorkspace)(workspace))
			return nil
		},
	}
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
