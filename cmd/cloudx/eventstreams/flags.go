// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"fmt"
	"net/url"

	"github.com/spf13/pflag"

	"github.com/ory/client-go"
)

type streamConfig client.CreateEventStreamBody

func (c *streamConfig) Validate() error {
	switch c.Type {
	case "":
		return fmt.Errorf("flag --type must be set")
	case "sns":
		if c.RoleArn == nil || *c.RoleArn == "" {
			return fmt.Errorf("flag --aws-iam-role-arn must be set")
		}
		if c.TopicArn == nil || *c.TopicArn == "" {
			return fmt.Errorf("flag --aws-sns-topic-arn must be set")
		}
		if c.HttpsEndpoint != nil && *c.HttpsEndpoint != "" {
			return fmt.Errorf("flag --https-endpoint cannot be set when type is sns")
		}
	case "https":
		if c.HttpsEndpoint == nil || *c.HttpsEndpoint == "" {
			return fmt.Errorf("flag --https-endpoint must be set")
		}
		e, err := url.Parse(*c.HttpsEndpoint)
		if err != nil {
			return fmt.Errorf("invalid URL for flag --https-endpoint: %w", err)
		}
		if e.Scheme != "https" {
			return fmt.Errorf("flag --https-endpoint must have https scheme")
		}
		if c.RoleArn != nil && *c.RoleArn != "" {
			return fmt.Errorf("flag --aws-iam-role-arn cannot be set when type is https")
		}
		if c.TopicArn != nil && *c.TopicArn != "" {
			return fmt.Errorf("flag --aws-sns-topic-arn cannot be set when type is https")
		}
	default:
		return fmt.Errorf("unsupported event stream type: %s", c.Type)
	}
	return nil
}

func registerStreamConfigFlags(f *pflag.FlagSet, c *streamConfig) {
	f.StringVar(&c.Type, "type", "", `The type of the event stream destination. Supported values are "sns" for AWS SNS topics and "https" for generic HTTPS endpoints.`)
	c.RoleArn = f.String("aws-iam-role-arn", "", "The ARN of the AWS IAM role to assume when publishing messages to the SNS topic.")
	c.TopicArn = f.String("aws-sns-topic-arn", "", "The ARN of the AWS SNS topic.")
	c.HttpsEndpoint = f.String("https-endpoint", "", "The URL of the HTTPS endpoint.")
}
