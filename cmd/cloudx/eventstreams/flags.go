// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"fmt"
	"net/url"

	"github.com/spf13/pflag"

	"github.com/ory/client-go"
)

// Event stream statuses. A paused stream does not forward any events until it
// is set back to active.
const (
	StatusActive = "active"
	StatusPaused = "paused"
)

type streamConfig client.CreateEventStreamBody

func (c *streamConfig) Validate() error {
	// The status flag is optional. An empty value is normalized to nil so the
	// server keeps the current status (on update) or applies its default (on create).
	if c.Status != nil {
		switch *c.Status {
		case "":
			c.Status = nil
		case StatusActive, StatusPaused:
		default:
			return fmt.Errorf(`flag --status must be one of %q or %q`, StatusActive, StatusPaused)
		}
	}

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

// toSetBody maps the shared stream config onto the update (set) request body.
// The two bodies are no longer convertible by type assertion: SetEventStreamBody.Type
// is a pointer (optional on update) whereas CreateEventStreamBody.Type is required.
func (c streamConfig) toSetBody() client.SetEventStreamBody {
	body := client.SetEventStreamBody{
		HttpsEndpoint: c.HttpsEndpoint,
		RoleArn:       c.RoleArn,
		Status:        c.Status,
		TopicArn:      c.TopicArn,
	}
	if c.Type != "" {
		t := c.Type
		body.Type = &t
	}
	return body
}

func registerStreamConfigFlags(f *pflag.FlagSet, c *streamConfig) {
	f.StringVar(&c.Type, "type", "", `The type of the event stream destination. Supported values are "sns" for AWS SNS topics and "https" for generic HTTPS endpoints.`)
	c.RoleArn = f.String("aws-iam-role-arn", "", "The ARN of the AWS IAM role to assume when publishing messages to the SNS topic.")
	c.TopicArn = f.String("aws-sns-topic-arn", "", "The ARN of the AWS SNS topic.")
	c.HttpsEndpoint = f.String("https-endpoint", "", "The URL of the HTTPS endpoint.")
	c.Status = f.String("status", "", fmt.Sprintf("The status of the event stream. Supported values are %q and %q. Defaults to %q.", StatusActive, StatusPaused, StatusActive))
}
