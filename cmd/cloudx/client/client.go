package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	kratos "github.com/ory/kratos-client-go"
	"github.com/ory/kratos/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const projectFlag = "project"

func RegisterProjectFlag(f *flag.FlagSet) {
	f.String(projectFlag, "", "The project to use")
}

func ContextWithClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, cliclient.ClientContextKey, func(cmd *cobra.Command) (*kratos.APIClient, error) {
		sc, err := NewCommandHelper(cmd)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize HTTP Client: %s\n", err)
			return nil, cmdx.FailSilently(cmd)
		}

		ac, err := sc.EnsureContext()
		if err != nil {
			return nil, err
		}

		project := uuid.FromStringOrNil(flagx.MustGetString(cmd, projectFlag))
		if project == uuid.Nil {
			_, _ = fmt.Fprintf(os.Stderr, "No project selected! Please use the flag --%s to specify one.\n", projectFlag)
			return nil, cmdx.FailSilently(cmd)
		}

		p, err := sc.GetProject(project.String())
		if err != nil {
			return nil, err
		}

		c := retryablehttp.NewClient()
		c.Logger = nil

		conf := kratos.NewConfiguration()
		conf.HTTPClient = &http.Client{
			Transport: &bearerTokenTransporter{RoundTripper: c.StandardClient().Transport, bearerToken: ac.SessionToken},
			Timeout:   time.Second * 10}

		conf.Servers = kratos.ServerConfigurations{{URL: makeCloudConsoleURL(p.Slug + ".projects")}}
		return kratos.NewAPIClient(conf), nil
	})
}
