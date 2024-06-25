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

func NewUpdateEventStreamCmd() *cobra.Command {
	c := streamConfig{}

	cmd := &cobra.Command{
		Use:   "event-stream id [--project=PROJECT_ID] [--type=sns] [--aws-iam-role-arn=arn:aws:iam::123456789012:role/MyRole] [--aws-sns-topic-arn=arn:aws:sns:us-east-1:123456789012:MyTopic]",
		Args:  cobra.ExactArgs(1),
		Short: "Update the event stream with the given ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := h.ProjectID()
			if err != nil {
				return err
			}
			streamID := args[0]

			if err := c.Validate(); err != nil {
				return err
			}
			stream, err := h.UpdateEventStream(ctx, projectID, streamID, cloud.SetEventStreamBody(c))
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Event stream updated successfully!")
			cmdx.PrintRow(cmd, output(*stream))
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())

	return cmd
}
