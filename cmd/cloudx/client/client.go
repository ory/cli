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

	"github.com/gofrs/uuid/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	hydra "github.com/ory/hydra-client-go"
	hydracli "github.com/ory/hydra/cmd/cliclient"
	kratoscli "github.com/ory/kratos/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const projectFlag = "project"

func RegisterProjectFlag(f *flag.FlagSet) {
	f.String(projectFlag, "", "The project to use")
}

// ProjectID returns the ID the user set with the `--project` flag, or prints a warning and returns an error
// if none was set.
func ProjectID(cmd *cobra.Command) (uuid.UUID, error) {
	project := uuid.FromStringOrNil(flagx.MustGetString(cmd, projectFlag))
	if project == uuid.Nil {
		_, _ = fmt.Fprintf(os.Stderr, "No project selected! Please use the flag --%s to specify one.\n", projectFlag)
		return uuid.Nil, cmdx.FailSilently(cmd)
	}
	return project, nil
}

func ProjectIDFromSlug(cmd *cobra.Command, pjs []cloud.ProjectMetadata) string {
	sl := flagx.MustGetString(cmd, projectFlag)
	var pid string
	for _, pm := range pjs {
		if pm.GetSlug() == sl {
			pid = pm.GetId()
		}
	}
	return pid
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

	pjs, err := sc.ListProjects()
	if err != nil {
		return nil, nil, nil, err
	}

	pid := ProjectIDFromSlug(cmd, pjs)
	if pid == "" {
		pid2, err := ProjectID(cmd)
		pid = pid2.String()
		if err != nil {
			return nil, nil, nil, cmdx.FailSilently(cmd)
		}
	}

	p, err := sc.GetProject(pid)
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
		c, ac, p, err := Client(cmd)
		if err != nil {
			return nil, nil, err
		}

		conf := hydra.NewConfiguration()
		conf.HTTPClient = &http.Client{
			Transport: &bearerTokenTransporter{RoundTripper: c.StandardClient().Transport, bearerToken: ac.SessionToken},
			Timeout:   time.Second * 10,
		}

		consoleURL, err := url.ParseRequestURI(makeCloudConsoleURL(p.Slug + ".projects"))
		if err != nil {
			return nil, nil, err
		}
		// We use the cloud console API because it works with ory cloud session tokens.
		conf.Servers = hydra.ServerConfigurations{{URL: consoleURL.String()}}
		return hydra.NewAPIClient(conf), consoleURL, nil
	})

	ctx = context.WithValue(ctx, kratoscli.ClientContextKey, func(cmd *cobra.Command) (*kratoscli.ClientContext, error) {
		c, ac, p, err := Client(cmd)
		if err != nil {
			return nil, err
		}

		// We use the cloud console API because it works with ory cloud session tokens.
		return &kratoscli.ClientContext{
			Endpoint: makeCloudConsoleURL(p.Slug + ".projects"),
			HTTPClient: &http.Client{
				Transport: &bearerTokenTransporter{
					RoundTripper: c.StandardClient().Transport,
					bearerToken:  ac.SessionToken,
				},
				Timeout: time.Second * 10,
			},
		}, nil
	})
	return ctx
}
