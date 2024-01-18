// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewListEventStreamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-streams",
		Args:  cobra.NoArgs,
		Short: "List your event streams",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			id, err := client.ProjectOrDefault(cmd, h)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			streams, err := h.ListEventStreams(id)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintTable(cmd, outputList(*streams))
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
