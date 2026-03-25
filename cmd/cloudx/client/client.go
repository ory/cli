// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"net/url"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
	hydracli "github.com/ory/hydra/v2/cmd/cliclient"
	kratoscli "github.com/ory/kratos/cmd/cliclient"
)

func ContextWithClient(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, hydracli.OAuth2URLOverrideContextKey, func(cmd *cobra.Command) *url.URL {
		h, err := NewCobraCommandHelper(cmd)
		if err != nil {
			return nil
		}
		project, err := h.GetSelectedProject(cmd.Context())
		if err != nil {
			return nil
		}

		return CloudAPIsURL(project.Slug + ".projects")
	})

	ctx = context.WithValue(ctx, hydracli.ClientContextKey, func(cmd *cobra.Command) (*hydra.APIClient, *url.URL, error) {
		ctx := cmd.Context()
		h, err := NewCobraCommandHelper(cmd)
		if err != nil {
			return nil, nil, err
		}

		c, baseURL, err := h.newProjectHTTPClient(ctx)
		if err != nil {
			return nil, nil, err
		}

		p, err := h.GetSelectedProject(ctx)
		if err != nil {
			return nil, nil, err
		}

		conf := hydra.NewConfiguration()
		conf.HTTPClient = c

		consoleProjectURL := baseURL(p.Slug + ".projects")
		// We use the cloud console API because it works with ory cloud session tokens.
		conf.Servers = hydra.ServerConfigurations{{URL: consoleProjectURL.String()}}
		return hydra.NewAPIClient(conf), consoleProjectURL, nil
	})

	ctx = context.WithValue(ctx, kratoscli.ClientContextKey, func(cmd *cobra.Command) (*kratoscli.ClientContext, error) {
		ctx := cmd.Context()
		h, err := NewCobraCommandHelper(cmd)
		if err != nil {
			return nil, err
		}

		c, baseURL, err := h.newProjectHTTPClient(ctx)
		if err != nil {
			return nil, err
		}

		p, err := h.GetSelectedProject(ctx)
		if err != nil {
			return nil, err
		}

		clientContext := &kratoscli.ClientContext{
			HTTPClient: c,
			Endpoint:   baseURL(p.Slug + ".projects").String(),
		}

		return clientContext, nil
	})
	return ctx
}
