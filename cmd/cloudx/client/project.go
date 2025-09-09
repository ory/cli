// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/client-go"
	"github.com/ory/x/jsonx"
)

func (h *CommandHelper) ListProjects(ctx context.Context, workspace *string) ([]client.ProjectMetadata, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	var projects []client.ProjectMetadata
	if workspace != nil {
		list, res, err := c.WorkspaceAPI.ListWorkspaceProjects(ctx, *workspace).Execute()
		if err != nil {
			return nil, handleError("unable to list workspace projects"+*workspace, res, err)
		}
		projects = list.Projects
	} else {
		var res *http.Response
		projects, res, err = c.ProjectAPI.ListProjects(ctx).Execute()
		if err != nil {
			return nil, handleError("unable to list projects", res, err)
		}
	}

	return projects, nil
}

func (h *CommandHelper) GetSelectedProject(ctx context.Context) (*client.ProjectMetadata, error) {
	id, err := h.ProjectID()
	if err != nil {
		return nil, err
	}

	if h.projectAPIKey != nil {
		pjs, err := h.ListProjects(ctx, nil)
		if err != nil {
			return nil, err
		}
		if len(pjs) != 1 {
			return nil, errors.Errorf("got unexpected number of projects when fetching with API key: %d", len(pjs))
		}
		return &pjs[0], nil
	}
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	project, res, err := c.ProjectAPI.GetProject(ctx, id).Execute()
	if err != nil {
		return nil, handleError("unable to get project", res, err)
	}

	return &client.ProjectMetadata{
		CreatedAt:            time.Time{},
		Environment:          project.Environment,
		HomeRegion:           project.HomeRegion,
		Hosts:                nil,
		Id:                   project.Id,
		Name:                 project.Name,
		Slug:                 project.Slug,
		State:                project.State,
		SubscriptionId:       client.NullableString{},
		SubscriptionPlan:     client.NullableString{},
		UpdatedAt:            time.Time{},
		Workspace:            nil,
		WorkspaceId:          project.WorkspaceId,
		AdditionalProperties: nil,
	}, nil
}

func (h *CommandHelper) GetProject(ctx context.Context, idOrSlug string, workspace *string) (*client.Project, error) {
	if idOrSlug == "" {
		return nil, errors.Errorf("No project selected! Please see the help message on how to set one.")
	}

	id, err := uuid.FromString(idOrSlug)
	if err != nil {
		projectMeta, err := h.findProject(ctx, idOrSlug, workspace)
		if err != nil {
			return nil, err
		}
		id = uuid.FromStringOrNil(projectMeta.GetId())
	}

	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	project, res, err := c.ProjectAPI.GetProject(ctx, id.String()).Execute()
	if err != nil {
		return nil, handleError("unable to get project", res, err)
	}

	return project, nil
}

func (h *CommandHelper) findProject(ctx context.Context, semiIdentifier string, workspace *string) (project *client.ProjectMetadata, _ error) {
	pjs, err := h.ListProjects(ctx, workspace)
	if err != nil {
		return nil, err
	}

	candidateSlugs := make([]string, 0, len(pjs))
	candidateIDs := make([]string, 0, len(pjs))
	allSlugs := make([]string, 0, len(pjs))
	allIDs := make([]string, 0, len(pjs))
	for _, pj := range pjs {
		allSlugs = append(allSlugs, pj.Slug)
		allIDs = append(allIDs, pj.Id)
		if strings.HasPrefix(pj.Slug, semiIdentifier) {
			candidateSlugs = append(candidateSlugs, pj.Slug)
			project = &pj
		}
		if strings.HasPrefix(pj.Id, semiIdentifier) {
			candidateIDs = append(candidateIDs, pj.Id)
			project = &pj
		}
	}
	if len(candidateSlugs)+len(candidateIDs) > 1 {
		return nil, errors.Errorf("The slug or ID prefix %q is not unique, please use more characters.\nMatching slugs: %v\nMatching IDs: %v", semiIdentifier, candidateSlugs, candidateIDs)
	}
	if project == nil {
		return nil, errors.Errorf("No project found with slug or ID %s.\nAll known slugs: %v\nAll known IDs: %v", semiIdentifier, allSlugs, allIDs)
	}
	return project, nil
}

func (h *CommandHelper) CreateProject(ctx context.Context, name, environment string, workspace *string, setDefault bool) (*client.Project, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	project, res, err := c.ProjectAPI.CreateProject(ctx).CreateProjectBody(client.CreateProjectBody{
		Name:        strings.TrimSpace(name),
		Environment: environment,
		WorkspaceId: workspace,
	}).Execute()
	if err != nil {
		return nil, handleError("unable to create project", res, err)
	}

	if setDefault || h.projectID == uuid.Nil {
		if err := h.SelectProject(project.Id); err != nil {
			return nil, fmt.Errorf("project created successfully, but could not select it: %w", err)
		}
		if workspace != nil {
			if err := h.SelectWorkspace(*workspace); err != nil {
				return nil, fmt.Errorf("project created successfully, but could not select workspace: %w", err)
			}
		}
	}

	return project, nil
}

func (h *CommandHelper) PatchProject(ctx context.Context, id string, raw []json.RawMessage, add, replace, del []string) (*client.SuccessfulProjectUpdate, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	var patches []client.JsonPatch
	for _, r := range raw {
		config, err := jsonx.EmbedSources(r, jsonx.WithIgnoreKeys("$id", "$schema"), jsonx.WithOnlySchemes("file"))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var p []client.JsonPatch
		if err := json.NewDecoder(bytes.NewReader(config)).Decode(&p); err != nil {
			return nil, errors.WithStack(err)
		}
		patches = append(patches, p...)
	}

	if v, err := toPatch("add", add); err != nil {
		return nil, err
	} else {
		//revive:disable indent-error-flow
		patches = append(patches, v...)
	}

	if v, err := toPatch("replace", replace); err != nil {
		return nil, err
	} else {
		//revive:disable indent-error-flow
		patches = append(patches, v...)
	}

	for _, del := range del {
		patches = append(patches, client.JsonPatch{Op: "remove", Path: del})
	}

	res, _, err := c.ProjectAPI.PatchProject(ctx, id).JsonPatch(patches).Execute()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *CommandHelper) UpdateProject(ctx context.Context, id string, name string, configs []json.RawMessage) (*client.SuccessfulProjectUpdate, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	for k := range configs {
		config, err := jsonx.EmbedSources(
			configs[k],
			jsonx.WithIgnoreKeys(
				"$id",
				"$schema",
			),
			jsonx.WithOnlySchemes(
				"file",
			),
		)
		if err != nil {
			return nil, err
		}
		configs[k] = config
	}

	interim := make(map[string]interface{})
	for _, config := range configs {
		var decoded map[string]interface{}
		if err := json.Unmarshal(config, &decoded); err != nil {
			return nil, errors.WithStack(err)
		}

		if err := mergo.Merge(&interim, decoded, mergo.WithAppendSlice, mergo.WithOverride); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	_, corsAdminFound := interim["cors_admin"]
	if !corsAdminFound {
		interim["cors_admin"] = map[string]interface{}{}
	}
	_, corsPublicFound := interim["cors_public"]
	if !corsPublicFound {
		interim["cors_public"] = map[string]interface{}{}
	}
	if _, found := interim["name"]; !found {
		interim["name"] = ""
	}
	if _, found := interim["organizations"]; !found {
		interim["organizations"] = []client.BasicOrganization{}
	}

	var payload client.SetProject
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(interim); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := json.NewDecoder(&b).Decode(&payload); err != nil {
		return nil, errors.WithStack(err)
	}

	if payload.Services.Identity == nil && payload.Services.Permission == nil && payload.Services.Oauth2 == nil {
		return nil, errors.Errorf("at least one of the keys `services.identity.config` and `services.permission.config` and `services.oauth2.config` is required and can not be empty")
	}
	if name != "" {
		payload.Name = name
	}

	// If either of the CORS keys is not set after the merge, we need to fetch it from the server
	// If the name is not set, and it was not provided, we need to fetch it from the server
	needsBackfill := !corsAdminFound || !corsPublicFound || payload.Name == ""

	if needsBackfill {
		res, _, err := c.ProjectAPI.GetProject(ctx, id).Execute()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if payload.Name == "" {
			payload.Name = res.Name
		}
		if !corsAdminFound {
			payload.CorsAdmin = *res.CorsAdmin
		}
		if !corsPublicFound {
			payload.CorsPublic = *res.CorsPublic
		}
	}

	res, _, err := c.ProjectAPI.SetProject(ctx, id).SetProject(payload).Execute()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (h *CommandHelper) PrintUpdateProjectWarnings(p *client.SuccessfulProjectUpdate) error {
	if len(p.Warnings) > 0 {
		_, _ = fmt.Fprintln(h.VerboseErrWriter)
		_, _ = fmt.Fprintln(h.VerboseErrWriter, "Warnings were found.")
		for _, warning := range p.Warnings {
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "- %s\n", *warning.Message)
		}
		_, _ = fmt.Fprintln(h.VerboseErrWriter, "It is safe to ignore these warnings unless your intention was to set these keys.")
	}

	_, _ = fmt.Fprintf(h.VerboseErrWriter, "\nProject updated successfully!\n")
	return nil
}
