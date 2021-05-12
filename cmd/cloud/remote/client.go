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

	"github.com/spf13/pflag"

	kratos "github.com/ory/kratos-client-go"
	"github.com/ory/x/cmdx"
)

const (
	FlagAPIEndpoint    = "api-endpoint"
	FlagConsoleAPI     = "console-url"
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

func IsUrl(str string) (*url.URL, bool, error) {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return nil, false, err
	}
	if u.Host == "" {
		return nil, false, errors.New(fmt.Sprintf("Could not parse requested url: %s", str))
	}
	return u, true, nil
}

func GetProjectSlug(consoleURL string) (string, error) {
	client := NewHTTPClient()
	u, ok, err := IsUrl(consoleURL)
	if err != nil || !ok || u == nil {
		return "", errors.WithStack(err)
	}
	uu := ""
	if u.Scheme != "" {
		uu = fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, tokenPath)
	} else {
		uu = fmt.Sprintf("%s/%s", u.Host, tokenPath)
	}
	rsp, err := client.Get(uu)
	if err != nil {
		return "", errors.WithStack(err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return gjson.GetBytes(body, "slug").String(), nil
}

func NewAdminClient(apiURL, consoleURL string) *kratos.APIClient {
	slug, err := GetProjectSlug(consoleURL)
	if err != nil {
		cmdx.Fatalf("Could not retrieve project slug: %s", errors.WithStack(err).Error())
	}
	if slug == "" {
		cmdx.Fatalf("No slug associated with given token")
	}
	api, err := url.ParseRequestURI(apiURL)
	if err != nil {
		cmdx.Must(err, "Unable to parse upstream URL because: %s", err)
	}
	uu := ""
	if api.Scheme != "" {
		uu = fmt.Sprintf("%s://%s.projects.%s/%s",api.Scheme, slug, api.Host, kratosAdminPath)
	} else {
		uu = fmt.Sprintf("https://%s.projects.%s/%s", slug, api.Host, kratosAdminPath)
	}
	upstream, err := url.ParseRequestURI(uu)
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
	flags.String(FlagConsoleAPI, "https://console.ory.sh", "Use a different URL.")
}
