// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"encoding/json"

	client "github.com/ory/client-go"
)

type (
	outputList client.ListEventStreams
	output     client.EventStream
)

func (output) Header() []string {
	return []string{"ID", "TYPE", "IAM_ROLE_ARN", "SNS_TOPIC_ARN", "HTTPS_ENDPOINT"}
}

func (o output) Columns() []string {
	return []string{
		coalesce(o.Id),
		coalesce(o.Type),
		coalesce(o.RoleArn),
		coalesce(o.TopicArn),
		coalesce(o.HttpsEndpoint.Get()),
	}
}

func coalesce(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (o output) Interface() interface{} {
	return o
}

func (outputList) Header() []string {
	return new(output).Header()
}

func (o outputList) Table() [][]string {
	rows := make([][]string, len(o.EventStreams))
	for i, stream := range o.EventStreams {
		rows[i] = (output)(stream).Columns()
	}
	return rows
}

func (o outputList) Interface() interface{} {
	return o
}

func (o outputList) Len() int {
	return len(o.EventStreams)
}

func (o outputList) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.EventStreams)
}
