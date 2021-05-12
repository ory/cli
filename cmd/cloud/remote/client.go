package remote

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/tidwall/gjson"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	kratos "github.com/ory/kratos-client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const (
	FlagAPIEndpoint    = "api-endpoint"
	FlagConsoleURL     = "console-url"
	projectAccessToken = "ORY_ACCESS_TOKEN"
	tokenPath          = "backoffice/token/slug"
	kratosAdminPath    = "api/kratos/admin"
)

type tokenTransporter struct {
	http.RoundTripper
	token string
}

func (t *tokenTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.RoundTripper.RoundTrip(req)
}

func NewHTTPClient() *http.Client {
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
	//if s, ok := cmd.Context().Value(TestKeyConstSlug).(string); ok {
	//	return s, nil
	//}
	client := NewHTTPClient()
	url, err := url.ParseRequestURI(flagx.MustGetString(cmd, FlagConsoleURL))
	if err != nil {
		return "", errors.WithStack(err)
	}
	rsp, err := client.Get(fmt.Sprintf("https://api.%s/%s", url.Host, tokenPath))
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
	slug, err := GetProjectSlug(cmd)
	if err != nil {
		cmdx.Fatalf("Could not retrieve project slug: %s", errors.WithStack(err).Error())
	}
	if slug == "" {
		cmdx.Fatalf("Could not retrieve valid project slug from %s", flagx.MustGetString(cmd, FlagConsoleURL))
	}
	upstream, err := url.ParseRequestURI(fmt.Sprintf("https://%s.projects.%s/%s", slug, flagx.MustGetString(cmd, FlagAPIEndpoint), kratosAdminPath))
	if err != nil {
		cmdx.Must(err, "Unable to parse upstream URL because: %s", err)
	}

	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: upstream.String()}}
	conf.HTTPClient = NewHTTPClient()

	return kratos.NewAPIClient(conf)
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.String(FlagAPIEndpoint, "https://oryapis.com", "Use a different endpoint.")
	flags.String(FlagConsoleURL, "https://console.ory.sh", "Use a different URL.")
}
