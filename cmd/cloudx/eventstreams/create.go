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

func NewCreateEventStreamCmd() *cobra.Command {
	c := streamConfig{}

	cmd := &cobra.Command{
		Use:   "event-stream [--project=PROJECT_ID] --type=sns --aws-iam-role-arn=arn:aws:iam::123456789012:role/MyRole --aws-sns-topic-arn=arn:aws:sns:us-east-1:123456789012:MyTopic",
		Short: "Create a new event stream",
		Args:  cobra.NoArgs,
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

			if err := c.Validate(); err != nil {
				return err
			}
			stream, err := h.CreateEventStream(ctx, projectID, cloud.CreateEventStreamBody(c))
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Event stream created successfully!")
			cmdx.PrintRow(cmd, output(*stream))
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())

	registerStreamConfigFlags(cmd.Flags(), &c)

	return cmd
}
