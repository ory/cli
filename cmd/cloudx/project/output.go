// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	"github.com/ory/x/cmdx"

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
	return []string{"ID", "NAME", "ENVIRONMENT", "WORKSPACE", "SLUG", "STATE"}
}

func (i *outputProject) Columns() []string {
	return []string{
		i.Id,
		i.Name,
		i.Environment,
		nullableStringOrNone(i.WorkspaceId.Get()),
		i.Slug,
		i.State,
	}
}

func (i *outputProject) Interface() interface{} {
	return i
}

func (*outputProjectCollection) Header() []string {
	return []string{"ID", "NAME", "ENVIRONMENT", "WORKSPACE", "SLUG", "STATE"}
}

func (c *outputProjectCollection) Table() [][]string {
	rows := make([][]string, len(c.projects))
	for i, ident := range c.projects {
		rows[i] = []string{
			ident.Id,
			ident.Name,
			ident.Environment,
			nullableStringOrNone(ident.WorkspaceId.Get()),
			ident.Slug,
			ident.State,
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

type selectedProject struct {
	ID string `json:"id"`
}

func (i selectedProject) String() string {
	return i.ID
}

func (i *selectedProject) ProjectID() string {
	return i.ID
}

func (*selectedProject) Header() []string {
	return []string{"ID"}
}

func (i *selectedProject) Columns() []string {
	return []string{
		i.ID,
	}
}

func (i *selectedProject) Interface() interface{} {
	return i
}

func nullableStringOrNone(s *string) string {
	if s == nil {
		return cmdx.None
	}
	return *s
}
