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

func NewCreateEventStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-stream [--project=PROJECT_ID] --type=sns --aws-iam-role-arn=arn:aws:iam::123456789012:role/MyRole --aws-sns-topic-arn=arn:aws:sns:us-east-1:123456789012:MyTopic",
		Short: "Create a new event stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := client.ProjectOrDefault(cmd, h)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			stream, err := h.CreateEventStream(projectID, cloud.CreateEventStreamBody{
				Type:     flagx.MustGetString(cmd, "type"),
				RoleArn:  flagx.MustGetString(cmd, "aws-iam-role-arn"),
				TopicArn: flagx.MustGetString(cmd, "aws-sns-topic-arn"),
			})
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Event stream created successfully!")
			cmdx.PrintRow(cmd, output(*stream))
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	cmd.Flags().String("type", "", `The type of the event stream destination. Only "sns" is supported at the moment.`)
	cmd.Flags().String("aws-iam-role-arn", "", "The ARN of the AWS IAM role to assume when publishing messages to the SNS topic.")
	cmd.Flags().String("aws-sns-topic-arn", "", "The ARN of the AWS SNS topic.")

	return cmd
}
