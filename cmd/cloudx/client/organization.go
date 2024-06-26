// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	cloud "github.com/ory/client-go"
)

func (h *CommandHelper) ListOrganizations(ctx context.Context, projectID string) (*cloud.ListOrganizationsResponse, error) {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return nil, err
	}

	organizations, res, err := c.ProjectAPI.ListOrganizations(ctx, projectID).Execute()
	if err != nil {
		return nil, handleError("unable to list organizations", res, err)
	}

	return organizations, nil
}

func (h *CommandHelper) CreateOrganization(ctx context.Context, projectID string, body cloud.OrganizationBody) (*cloud.Organization, error) {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return nil, err
	}

	organization, res, err := c.ProjectAPI.
		CreateOrganization(ctx, projectID).
		OrganizationBody(body).
		Execute()
	if err != nil {
		return nil, handleError("unable to create organization", res, err)
	}

	return organization, nil
}

func (h *CommandHelper) UpdateOrganization(ctx context.Context, projectID, orgID string, body cloud.OrganizationBody) (*cloud.Organization, error) {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return nil, err
	}

	organization, res, err := c.ProjectAPI.
		UpdateOrganization(ctx, projectID, orgID).
		OrganizationBody(body).
		Execute()
	if err != nil {
		return nil, handleError("unable to update organization", res, err)
	}

	return organization, nil
}

func (h *CommandHelper) DeleteOrganization(ctx context.Context, projectID, orgID string) error {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return err
	}

	res, err := c.ProjectAPI.
		DeleteOrganization(ctx, projectID, orgID).
		Execute()
	if err != nil {
		return handleError("unable to delete organization", res, err)
	}

	return nil
}
