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
	kratos "github.com/ory/kratos-client-go"

	hydracli "github.com/ory/hydra/cmd/cliclient"
	kratoscli "github.com/ory/kratos/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const projectFlag = "project"

func RegisterProjectFlag(f *flag.FlagSet) {
	f.String(projectFlag, "", "The project to use")
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

	project := uuid.FromStringOrNil(flagx.MustGetString(cmd, projectFlag))
	if project == uuid.Nil {
		_, _ = fmt.Fprintf(os.Stderr, "No project selected! Please use the flag --%s to specify one.\n", projectFlag)
		return nil, nil, nil, cmdx.FailSilently(cmd)
	}

	p, err := sc.GetProject(project.String())
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
			Timeout:   time.Second * 10}

		consoleURL, err := url.ParseRequestURI(makeCloudConsoleURL(p.Slug + ".projects"))
		if err != nil {
			return nil, nil, err
		}
		// We use the cloud console API because it works with ory cloud session tokens.
		conf.Servers = hydra.ServerConfigurations{{URL: consoleURL.String()}}
		return hydra.NewAPIClient(conf), consoleURL, nil
	})

	ctx = context.WithValue(ctx, kratoscli.ClientContextKey, func(cmd *cobra.Command) (*kratos.APIClient, error) {
		c, ac, p, err := Client(cmd)
		if err != nil {
			return nil, err
		}
		conf := kratos.NewConfiguration()
		conf.HTTPClient = &http.Client{
			Transport: &bearerTokenTransporter{RoundTripper: c.StandardClient().Transport, bearerToken: ac.SessionToken},
			Timeout:   time.Second * 10}

		// We use the cloud console API because it works with ory cloud session tokens.
		conf.Servers = kratos.ServerConfigurations{{URL: makeCloudConsoleURL(p.Slug + ".projects")}}
		return kratos.NewAPIClient(conf), nil
	})
	return ctx
}
