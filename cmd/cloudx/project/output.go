// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	cloud "github.com/ory/client-go"
)

type (
	outputConfig            map[string]interface{}
	outputProject           cloud.Project
	outputProjectCollection struct {
		projects []cloud.ProjectMetadata
	}
)

func (i outputConfig) String() string {
	return fmt.Sprintf("%+v", map[string]interface{}(i))
}

func (i *outputProject) ID() string {
	return i.Id
}

func (*outputProject) Header() []string {
	return []string{"ID", "SLUG", "STATE", "NAME"}
}

func (i *outputProject) Columns() []string {
	return []string{
		i.Id,
		i.Slug,
		i.State,
		i.Name,
	}
}

func (i *outputProject) Interface() interface{} {
	return i
}

func (*outputProjectCollection) Header() []string {
	return []string{"ID", "SLUG", "STATE", "NAME"}
}

func (c *outputProjectCollection) Table() [][]string {
	rows := make([][]string, len(c.projects))
	for i, ident := range c.projects {
		rows[i] = []string{
			ident.Id,
			func() string {
				if ident.Slug != nil {
					return *ident.Slug
				}
				return "<none>"
			}(),
			ident.State,
			ident.Name,
		}
	}
	return rows
}

func (c *outputProjectCollection) Interface() interface{} {
	return c.projects
}

func (c *outputProjectCollection) Len() int {
	return len(c.projects)
}
