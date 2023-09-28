// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package organizations

import (
	"encoding/json"
	"strings"

	client "github.com/ory/client-go"
)

type (
	outputOrganizations struct {
		organizations *client.ListOrganizationsResponse
	}
	output client.Organization
)

func (o output) Header() []string {
	return []string{"ID", "LABEL", "DOMAINS"}
}

func (o output) Columns() []string {
	return []string{
		o.Id,
		o.Label,
		strings.Join(o.Domains, ", "),
	}
}

func (o output) Interface() interface{} {
	return o
}

func (o *outputOrganizations) Header() []string {
	return new(output).Header()
}

func (o *outputOrganizations) Table() [][]string {
	rows := make([][]string, o.Len())
	for i, organization := range o.organizations.Organizations {
		rows[i] = (output)(organization).Columns()
	}
	return rows
}

func (o *outputOrganizations) Interface() interface{} {
	return o
}

func (o *outputOrganizations) Len() int {
	return len(o.organizations.Organizations)
}

func (o *outputOrganizations) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.organizations.Organizations)
}
