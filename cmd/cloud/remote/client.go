package remote

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	kratos "github.com/ory/kratos-client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/stringsx"
)

const (
	FlagProject        = "project"
	FlagEndpoint       = "endpoint"
	projectEnvKey      = "ORY_PROJECT_ID"
	projectAccessToken = "ORY_ACCESS_TOKEN"
)

type tokenTransporter struct {
	http.RoundTripper
	token string
}

func (t *tokenTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.RoundTripper.RoundTrip(req)
}

func NewHTTPClient(cmd *cobra.Command) *http.Client {
	token := os.Getenv(projectAccessToken)
	if len(token) == 0 {
		cmdx.Fatalf(`Ory API Token could not be detected! Did you forget to set the environment variable "%s"?

You can create an API Token in the Ory Console. Once created, set the environment variable as follows.

**Unix (Linux, macOS)**

$ export ORY_ACCESS_TOKEN="<your-api-token-here>"
$ ory ...

**Windows (Powershell)**

> $env:ORY_ACCESS_TOKEN = '<your-api-token-here>'
> ory ...

**Windows (cmd.exe)**

> set ORY_ACCESS_TOKEN = "<your-api-token-here>"
> ory ...
`, projectAccessToken)
		return nil
	}

	return &http.Client{
		Transport: &tokenTransporter{
			RoundTripper: http.DefaultTransport,
			token:        token,
		},
		Timeout: time.Second * 10,
	}
}

func NewAdminClient(cmd *cobra.Command) *kratos.APIClient {
	project := stringsx.Coalesce(flagx.MustGetString(cmd, FlagProject), os.Getenv(projectEnvKey))
	if project == "" {
		cmdx.Fatalf("You have to set the Ory Cloud Project ID, try --help for details.")
	}

	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: flagx.MustGetString(cmd, FlagEndpoint)}}
	return kratos.NewAPIClient(conf)
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.StringP(FlagProject, FlagProject[:1], "", fmt.Sprintf("Must be set to your Ory Cloud Project Slug. Alternatively set using the %s environmental variable.", projectEnvKey))
	flags.String(FlagEndpoint, "https://oryapis.com", "Use a different endpoint.")
}
