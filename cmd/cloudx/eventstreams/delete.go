// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewDeleteEventStream() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-stream <id> [--project=PROJECT_ID]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete the event stream with the given ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := h.ProjectID()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			streamID := args[0]

			err = h.DeleteEventStream(cmd.Context(), projectID, streamID)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Event stream deleted successfully!")
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	return cmd
}
