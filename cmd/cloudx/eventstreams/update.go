// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

func NewUpdateEventStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-stream id [--project=PROJECT_ID] [--type=sns] [--aws-iam-role-arn=arn:aws:iam::123456789012:role/MyRole] [--aws-sns-topic-arn=arn:aws:sns:us-east-1:123456789012:MyTopic]",
		Args:  cobra.ExactArgs(1),
		Short: "Update the event stream with the given ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := client.ProjectOrDefault(cmd, h)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			streamID := args[0]

			stream, err := h.UpdateEventStream(projectID, streamID, cloud.SetEventStreamBody{
				Type:     flagx.MustGetString(cmd, "type"),
				RoleArn:  flagx.MustGetString(cmd, "aws-iam-role-arn"),
				TopicArn: flagx.MustGetString(cmd, "aws-sns-topic-arn"),
			})
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Event stream updated successfully!")
			cmdx.PrintRow(cmd, output(*stream))
			return nil
		},
	}

	cmd.Flags().String("type", "", `The type of the event stream destination. Only "sns" is supported at the moment.`)
	cmd.Flags().String("aws-iam-role-arn", "", "The ARN of the AWS IAM role to assume when publishing messages to the SNS topic.")
	cmd.Flags().String("aws-sns-topic-arn", "", "The ARN of the AWS SNS topic.")
	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())

	return cmd
}
