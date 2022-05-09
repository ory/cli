package client

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ory/x/stringsx"

	cloud "github.com/ory/client-go"
)

func makeCloudConsoleURL(prefix string) string {
	u, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv("ORY_CLOUD_CONSOLE_URL"), "https://console.ory.sh"))
	if err != nil {
		u = &url.URL{Scheme: "https", Host: "console.ory.sh"}
	}
	u.Host = prefix + "." + u.Host
	return u.Scheme + "://" + u.Host
}

func NewKratosClient() (*cloud.APIClient, error) {
	conf := cloud.NewConfiguration()
	conf.Servers = cloud.ServerConfigurations{{URL: makeCloudConsoleURL("project")}}
	conf.HTTPClient = &http.Client{Timeout: time.Second * 10}

	return cloud.NewAPIClient(conf), nil
}

func newCloudClient(token string) (*cloud.APIClient, error) {
	u := makeCloudConsoleURL("api")

	conf := cloud.NewConfiguration()
	conf.Servers = cloud.ServerConfigurations{{URL: u}}
	conf.HTTPClient = newBearerTokenClient(token)

	return cloud.NewAPIClient(conf), nil
}
