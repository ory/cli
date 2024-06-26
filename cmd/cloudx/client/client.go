// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	cloud "github.com/ory/client-go"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
	hydracli "github.com/ory/hydra/v2/cmd/cliclient"
	kratoscli "github.com/ory/kratos/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func Client(cmd *cobra.Command) (*retryablehttp.Client, *Config, *cloud.Project, error) {
	ctx := cmd.Context()
	sc, err := NewCobraCommandHelper(cmd)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize HTTP Client: %s\n", err)
		return nil, nil, nil, cmdx.FailSilently(cmd)
	}

	ac, err := sc.GetAuthenticatedConfig(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	project, err := sc.GetSelectedProject(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	c := retryablehttp.NewClient()
	c.Logger = nil
	return c, ac, project, nil
}

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

		return CloudAPIsURL(project.Slug)
	})

	ctx = context.WithValue(ctx, hydracli.ClientContextKey, func(cmd *cobra.Command) (*hydra.APIClient, *url.URL, error) {
		c, ac, p, err := Client(cmd)
		if err != nil {
			return nil, nil, err
		}

		conf := hydra.NewConfiguration()
		conf.HTTPClient = &http.Client{
			Transport: &bearerTokenTransporter{RoundTripper: c.StandardClient().Transport, bearerToken: ac.SessionToken},
			Timeout:   time.Second * 30,
		}

		consoleProjectURL := cloudConsoleURL(p.Slug + ".projects")
		// We use the cloud console API because it works with ory cloud session tokens.
		conf.Servers = hydra.ServerConfigurations{{URL: consoleProjectURL.String()}}
		return hydra.NewAPIClient(conf), consoleProjectURL, nil
	})

	ctx = context.WithValue(ctx, kratoscli.ClientContextKey, func(cmd *cobra.Command) (*kratoscli.ClientContext, error) {
		c, ac, p, err := Client(cmd)
		if err != nil {
			return nil, err
		}

		// We use the cloud console API because it works with ory cloud session tokens.
		return &kratoscli.ClientContext{
			Endpoint: cloudConsoleURL(p.Slug + ".projects").String(),
			HTTPClient: &http.Client{
				Transport: &bearerTokenTransporter{
					RoundTripper: c.StandardClient().Transport,
					bearerToken:  ac.SessionToken,
				},
				Timeout: time.Second * 30,
			},
		}, nil
	})
	return ctx
}
