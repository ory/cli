package remote

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/tidwall/gjson"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	kratos "github.com/ory/kratos-client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const (
	FlagEndpoint       = "endpoint"
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
			RoundTripper: retryablehttp.NewClient().StandardClient().Transport,
			token:        token,
		},
		Timeout: time.Second * 10,
	}
}

func GetProjectSlug(cmd *cobra.Command) (string, error) {
	url := flagx.MustGetString(cmd, FlagEndpoint)
	client := NewHTTPClient(cmd)

	rsp, err := client.Get(fmt.Sprintf("%s/token/slug", url))
	if err != nil {
		return "", errors.WithStack(err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return gjson.GetBytes(body, "slug").String(), nil
}

func NewAdminClient(cmd *cobra.Command) *kratos.APIClient {
	_, err := GetProjectSlug(cmd)
	if err != nil {
		cmdx.Fatalf("Could not retrieve project slug: %s", errors.WithStack(err).Error())
	}

	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: flagx.MustGetString(cmd, FlagEndpoint)}}
	return kratos.NewAPIClient(conf)
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.String(FlagEndpoint, "https://oryapis.com", "Use a different endpoint.")
}
