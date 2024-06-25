// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/ory/client-go"
)

type streamConfig client.CreateEventStreamBody

func (c *streamConfig) Validate() error {
	switch "" {
	case c.Type:
		return fmt.Errorf("flag --type must be set")
	case c.RoleArn:
		return fmt.Errorf("flag --aws-iam-role-arn must be set")
	case c.TopicArn:
		return fmt.Errorf("flag --aws-sns-topic-arn must be set")
	}
	return nil
}

func registerStreamConfigFlags(f *pflag.FlagSet, c *streamConfig) {
	f.StringVar(&c.Type, "type", "", `The type of the event stream destination. Only "sns" is supported at the moment.`)
	f.StringVar(&c.RoleArn, "aws-iam-role-arn", "", "The ARN of the AWS IAM role to assume when publishing messages to the SNS topic.")
	f.StringVar(&c.TopicArn, "aws-sns-topic-arn", "", "The ARN of the AWS SNS topic.")
}
