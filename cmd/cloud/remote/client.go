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

	c := retryablehttp.NewClient()
	c.Logger = nil

	return &http.Client{
		Transport: &tokenTransporter{
			RoundTripper: c.StandardClient().Transport,
			token:        token,
		},
		Timeout: time.Second * 10,
	}
}

func IsUrl(str string) (*url.URL, error) {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return nil, err
	}
	if u.Host == "" {
		return nil, errors.New(fmt.Sprintf("Could not parse requested url: %s", str))
	}
	return u, nil
}

func GetProjectSlug(consoleURL string) (string, error) {
	client := NewHTTPClient()
	u, err := IsUrl(consoleURL)
	if err != nil || u == nil {
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

func CreateKratosURL(apiURL, consoleURL string) (*url.URL, error) {
	slug, err := GetProjectSlug(consoleURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if slug == "" {
		return nil, errors.New("Could not retrieve slug from requested url")
	}
	api, err := IsUrl(apiURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	uu := ""
	if api.Scheme != "" {
		uu = fmt.Sprintf("%s://%s.projects.%s/%s", api.Scheme, slug, api.Host, kratosAdminPath)
	} else {
		uu = fmt.Sprintf("https://%s.projects.%s/%s", slug, api.Host, kratosAdminPath)
	}
	upstream, err := IsUrl(uu)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return upstream, nil
}

func NewAdminClient(apiURL, consoleURL string) *kratos.APIClient {
	u, err := CreateKratosURL(apiURL, consoleURL)
	if err != nil {
		cmdx.Fatalf("Could not retrieve project slug: %s", errors.WithStack(err).Error())
	}
	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: u.String()}}
	conf.HTTPClient = NewHTTPClient()

	return kratos.NewAPIClient(conf)
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.String(FlagAPIEndpoint, "https://oryapis.com", "Use a different endpoint.")
	flags.String(FlagConsoleAPI, "https://api.console.ory.sh", "Use a different URL.")
}
