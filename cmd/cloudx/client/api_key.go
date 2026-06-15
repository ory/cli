// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
)

// CreateProjectAPIKey creates a project API key. If expiresIn is greater than
// zero, the key is set to expire that duration from now so it is cleaned up
// automatically on the server side even if local cleanup fails. An expiresIn of
// zero creates a key without expiry; a negative value is rejected.
func (h *CommandHelper) CreateProjectAPIKey(ctx context.Context, projectID, name string, expiresIn time.Duration) (*cloud.ProjectApiKey, error) {
	if expiresIn < 0 {
		return nil, errors.New("API key expiry must not be negative")
	}

	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	req := cloud.CreateProjectApiKeyRequest{Name: name}
	if expiresIn > 0 {
		expiresAt := time.Now().Add(expiresIn)
		req.ExpiresAt = &expiresAt
	}

	token, res, err := c.ProjectAPI.CreateProjectApiKey(ctx, projectID).CreateProjectApiKeyRequest(req).Execute()
	if err != nil {
		return nil, handleError("unable to create project API key", res, err)
	}

	return token, nil
}

func (h *CommandHelper) DeleteProjectAPIKey(ctx context.Context, projectID, keyID string) error {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return err
	}

	if res, err := c.ProjectAPI.DeleteProjectApiKey(ctx, projectID, keyID).Execute(); err != nil {
		return handleError("unable to delete project API key", res, err)
	}

	return nil
}

func (h *CommandHelper) CreateWorkspaceAPIKey(ctx context.Context, workspaceID, name string) (*cloud.WorkspaceApiKey, error) {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	key, res, err := c.WorkspaceAPI.CreateWorkspaceApiKey(ctx, workspaceID).CreateWorkspaceApiKeyBody(cloud.CreateWorkspaceApiKeyBody{Name: name}).Execute()
	if err != nil {
		return nil, handleError("unable to create workspace API key", res, err)
	}
	return key, nil
}

func (h *CommandHelper) DeleteWorkspaceAPIKey(ctx context.Context, workspaceID, keyID string) error {
	c, err := h.newConsoleAPIClient(ctx)
	if err != nil {
		return err
	}

	if res, err := c.WorkspaceAPI.DeleteWorkspaceApiKey(ctx, workspaceID, keyID).Execute(); err != nil {
		return handleError("unable to delete workspace API key", res, err)
	}
	return nil
}

// TemporaryAPIKey creates a short-lived project API key that is deleted via the
// returned cleanup function. The key is additionally set to expire after
// expiresIn so that it is removed automatically should the cleanup fail. An
// expiresIn of zero creates a key without expiry.
func (h *CommandHelper) TemporaryAPIKey(ctx context.Context, name string, expiresIn time.Duration) (apiKey string, cleanup func() error, err error) {
	if h.projectAPIKey != nil {
		return *h.projectAPIKey, noop, nil
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
	ak, err := h.CreateProjectAPIKey(ctx, projectID, name, expiresIn)
	if err != nil {
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "Unable to create API key. Do you have the required permissions to use the Ory CLI with project %q? Continuing without API key.", projectID)
		return "", noop, nil
	}

	if !ak.HasValue() {
		return "", noop, nil
	}

	return *ak.Value, func() error {
		return h.DeleteProjectAPIKey(ctx, projectID, ak.Id)
	}, nil
}

func noop() error { return nil }
