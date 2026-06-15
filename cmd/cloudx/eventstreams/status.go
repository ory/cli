// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
)

func NewPauseEventStreamCmd() *cobra.Command {
	return newSetStatusCmd("pause", StatusPaused, "Pause the event stream with the given ID", "A paused event stream does not forward any events until it is resumed.")
}

func NewResumeEventStreamCmd() *cobra.Command {
	return newSetStatusCmd("resume", StatusActive, "Resume the event stream with the given ID", "Resuming a paused event stream makes it forward events again.")
}

func newSetStatusCmd(verb, status, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-stream <id> [--project=PROJECT_ID]",
		Args:  cobra.ExactArgs(1),
		Short: short,
		Long:  short + "\n\n" + long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := h.ProjectID()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			streamID := args[0]

			stream, err := h.UpdateEventStream(ctx, projectID, streamID, cloud.SetEventStreamBody{Status: &status})
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintf(h.VerboseErrWriter, "Event stream %sd successfully!\n", verb)
			cmdx.PrintRow(cmd, output(*stream))
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
