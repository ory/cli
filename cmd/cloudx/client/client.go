// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"golang.org/x/oauth2"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ory/cli/buildinfo"
	cloud "github.com/ory/client-go"
	hydra "github.com/ory/hydra-client-go/v2"
	hydracli "github.com/ory/hydra/v2/cmd/cliclient"
	kratoscli "github.com/ory/kratos/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const projectFlag = "project"

func RegisterProjectFlag(f *flag.FlagSet) {
	f.String(projectFlag, "", "The project to use, either project ID or a (partial) slug.")
}

// ProjectOrDefault returns the slug or ID the user set with the `--project` flag, or the default project, or prints a warning and returns an error
// if none was set.
func ProjectOrDefault(cmd *cobra.Command, h *CommandHelper) (string, error) {
	if flag := flagx.MustGetString(cmd, projectFlag); flag == "" {
		if id := h.GetDefaultProjectID(); id != "" {
			return id, nil
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "No project selected! Please use the flag --%s to specify one.\n", projectFlag)
			return "", cmdx.FailSilently(cmd)
		}
	} else {
		return flag, nil
	}
}

func Client(cmd *cobra.Command) (*retryablehttp.Client, *AuthContext, *cloud.Project, error) {
	sc, err := NewCommandHelper(cmd)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize HTTP Client: %s\n", err)
		return nil, nil, nil, cmdx.FailSilently(cmd)
	}

	ac, err := sc.EnsureContext()
	if err != nil {
		return nil, nil, nil, err
	}

	projectOrSlug, err := ProjectOrDefault(cmd, sc)
	if err != nil {
		return nil, nil, nil, cmdx.FailSilently(cmd)
	}

	p, err := sc.GetProject(projectOrSlug)
	if err != nil {
		return nil, nil, nil, err
	}

	c := retryablehttp.NewClient()
	c.Logger = nil
	return c, ac, p, nil
}

func ContextWithClient(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, hydracli.OAuth2URLOverrideContextKey, func(cmd *cobra.Command) *url.URL {
		_, _, p, err := Client(cmd)
		if err != nil {
			return nil
		}

		apiURL, err := url.ParseRequestURI(makeCloudAPIsURL(p.Slug + ".projects"))
		if err != nil {
			return nil
		}

		// We use the cloud console API because it works with ory cloud session tokens.
		return apiURL
	})

	ctx = context.WithValue(ctx, hydracli.ClientContextKey, func(cmd *cobra.Command) (*hydra.APIClient, *url.URL, error) {
		_, ac, p, err := Client(cmd)
		if err != nil {
			return nil, nil, err
		}

		// We use the cloud console API because it works with ory cloud session tokens.
		consoleURL, err := url.ParseRequestURI(makeCloudConsoleURL(p.Slug + ".projects"))
		if err != nil {
			return nil, nil, err
		}
		conf := hydra.NewConfiguration()
		conf.Servers = hydra.ServerConfigurations{{URL: consoleURL.String(), Variables: make(map[string]hydra.ServerVariable)}}
		conf.Debug = true
		conf.UserAgent = "ory-cli/" + buildinfo.Version

		cmd.SetContext(context.WithValue(cmd.Context(), hydra.ContextOAuth2, oac.TokenSource(cmd.Context(), ac.AccessToken)))

		return hydra.NewAPIClient(conf), consoleURL, nil
	})

	ctx = context.WithValue(ctx, kratoscli.ClientContextKey, func(cmd *cobra.Command) (*kratoscli.ClientContext, error) {
		c, ac, p, err := Client(cmd)
		if err != nil {
			return nil, err
		}

		// We use the cloud console API because it works with ory cloud session tokens.
		return &kratoscli.ClientContext{
			Endpoint:   makeCloudConsoleURL(p.Slug + ".projects"),
			HTTPClient: oac.Client(context.WithValue(context.Background(), oauth2.HTTPClient, c.StandardClient()), ac.AccessToken),
		}, nil
	})
	return ctx
}
