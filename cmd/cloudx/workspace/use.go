// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewUseWorkspaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workspace [id]",
		Aliases: []string{"ws"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Set the workspace as the default. When no id is provided, prints the currently used default workspace.",
		Example: `$ ory use workspace ecaaa3cb-0730-4ee8-a6df-9553cdfeef89

ID		ecaaa3cb-0730-4ee8-a6df-9553cdfeef89

$ ory use workspace ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format json

{
  "id": "ecaaa3cb-0730-4ee8-a6df-9553cdfeef89"
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := make([]client.CommandHelperOption, 0, 1)
			if len(args) == 1 {
				opts = append(opts, client.WithWorkspaceOverride(args[0]))
			}
			h, err := client.NewCobraCommandHelper(cmd, opts...)
			if err != nil {
				return err
			}

			// The helper resolves a (partial) name or ID to the workspace ID,
			// and falls back to the one already stored in the config when no
			// argument is given.
			id := h.WorkspaceID()
			if id == nil {
				return client.ErrWorkspaceNotSet
			}

			if err := h.SelectWorkspace(*id); err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintRow(cmd, &selectedWorkspace{ID: *id})
			return nil
		},
	}

	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
