// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"errors"
	"fmt"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
)

func (h *CommandHelper) CreateAPIKey(ctx context.Context, projectID, name string) (*cloud.ProjectApiKey, error) {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return nil, err
	}

	token, _, err := c.ProjectAPI.CreateProjectApiKey(ctx, projectID).CreateProjectApiKeyRequest(cloud.CreateProjectApiKeyRequest{Name: name}).Execute()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (h *CommandHelper) DeleteAPIKey(ctx context.Context, projectIdOrSlug, id string) error {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return err
	}

	if _, err := c.ProjectAPI.DeleteProjectApiKey(ctx, projectIdOrSlug, id).Execute(); err != nil {
		return err
	}

	return nil
}

func (h *CommandHelper) TemporaryAPIKey(ctx context.Context, name string) (apiKey string, cleanup func() error, err error) {
	if ak := GetProjectAPIKeyFromEnvironment(); len(ak) > 0 {
		return ak, noop, nil
	}

	// For all other projects, except the playground, we need to authenticate.
	if err := h.Authenticate(ctx); errors.Is(err, ErrNoConfigQuiet) {
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "Because you are not authenticated, the Ory CLI can not configure your project automatically. You can still use the Ory Proxy / Ory Tunnel, but complex flows such as Social Sign In will not work. Remove the `--quiet` flag or run `ory auth login` to authenticate.")
		return "", noop, nil
	} else if errors.Is(err, ErrNotAuthenticated) {
		ok, err := cmdx.AskScannerForConfirmation("To support complex flows such as Social Sign In, the Ory CLI can configure your project automatically. To do so, you need to be signed in. Do you want to sign in?", h.Stdin, h.VerboseErrWriter)
		if err != nil {
			return "", noop, err
		}

		if !ok {
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "Because you are not authenticated, the Ory CLI can not configure your project automatically. You can still use the Ory Proxy / Ory Tunnel, but complex flows such as Social Sign In will not work.")
			return "", noop, nil
		}

		if err := h.Authenticate(ctx); err != nil {
			return "", noop, err
		}
	} else if err != nil {
		return "", noop, err
	}

	projectID, err := h.ProjectID()
	if err != nil {
		return "", noop, err
	}
	ak, err := h.CreateAPIKey(ctx, projectID, name)
	if err != nil {
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "Unable to create API key. Do you have the required permissions to use the Ory CLI with project %q? Continuing without API key.", projectID)
		return "", noop, nil
	}

	if !ak.HasValue() {
		return "", noop, nil
	}

	return *ak.Value, func() error {
		return h.DeleteAPIKey(ctx, projectID, ak.Id)
	}, nil
}

func noop() error { return nil }
