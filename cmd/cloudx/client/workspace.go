// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	cloud "github.com/ory/client-go"
)

func (h *CommandHelper) ListWorkspaces(ctx context.Context) ([]cloud.Workspace, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}
	list, _, err := c.WorkspaceAPI.ListWorkspaces(ctx).Execute()
	if err != nil {
		return nil, err
	}
	return list.Workspaces, nil
}

func (h *CommandHelper) findWorkspace(ctx context.Context, semiIdentifier string) (workspace *cloud.Workspace, _ error) {
	wss, err := h.ListWorkspaces(ctx)
	if err != nil {
		return nil, err
	}

	candidateNames := make([]string, 0, len(wss))
	candidateIDs := make([]string, 0, len(wss))
	allNames := make([]string, 0, len(wss))
	allIDs := make([]string, 0, len(wss))
	for _, ws := range wss {
		allNames = append(allNames, ws.Name)
		allIDs = append(allIDs, ws.Id)
		if strings.HasPrefix(ws.Name, semiIdentifier) {
			candidateNames = append(candidateNames, ws.Name)
			workspace = &ws
		}
		if strings.HasPrefix(ws.Id, semiIdentifier) {
			candidateIDs = append(candidateIDs, ws.Id)
			workspace = &ws
		}
	}
	if len(candidateNames)+len(candidateIDs) > 1 {
		return nil, errors.Errorf("Found more than one workspace matching the identifier %q.\nMatching names: %v\nMatching IDs: %v", semiIdentifier, candidateNames, candidateIDs)
	}
	if workspace == nil {
		return nil, errors.Errorf("No workspace found with the identifier %q.\nAll known names: %v\nAll known IDs: %v", semiIdentifier, allNames, allIDs)
	}
	return workspace, nil
}

func (h *CommandHelper) CreateWorkspace(ctx context.Context, name string) (*cloud.Workspace, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	workspace, res, err := c.WorkspaceAPI.CreateWorkspace(ctx).CreateWorkspaceBody(cloud.CreateWorkspaceBody{Name: name}).Execute()
	if err != nil {
		fmt.Printf("raw response: %+v", res)
		return nil, err
	}
	return workspace, nil
}

func (h *CommandHelper) GetWorkspace(ctx context.Context, id string) (*cloud.Workspace, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	workspace, res, err := c.WorkspaceAPI.GetWorkspace(ctx, id).Execute()
	if err != nil {
		return nil, handleError("unable to get workspace", res, err)
	}
	return workspace, nil
}
