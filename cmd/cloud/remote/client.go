package remote

import (
	"bytes"
	"fmt"
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
			RoundTripper: retryablehttp.NewClient().StandardClient().Transport,
			token:        token,
		},
		Timeout: time.Second * 10,
	}
}

func GetProjectSlug(cmd *cobra.Command) (string, error) {
	url := flagx.MustGetString(cmd, FlagEndpoint)
	client := NewHTTPClient(cmd)

	b := &bytes.Buffer{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/slug", url, os.Getenv(projectEnvKey)), b)
	if err != nil {
		return "", err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return gjson.GetBytes(body, "slug").Str, nil
}

func NewAdminClient(cmd *cobra.Command) *kratos.APIClient {
	project, err := GetProjectSlug(cmd)
	if project == "" || err != nil {
		cmdx.Fatalf("Could not retrieve project slug: %+v", err)
	}

	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: flagx.MustGetString(cmd, FlagEndpoint)}}
	return kratos.NewAPIClient(conf)
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.StringP(FlagProject, FlagProject[:1], "", fmt.Sprintf("Must be set to your Ory Cloud Project Slug. Alternatively set using the %s environmental variable.", projectEnvKey))
	flags.String(FlagEndpoint, "https://oryapis.com", "Use a different endpoint.")
}
