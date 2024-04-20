// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/ory/cli/buildinfo"
	cloud "github.com/ory/client-go"
	"github.com/ory/x/stringsx"
)

var RateLimitHeader = os.Getenv("ORY_RATE_LIMIT_HEADER")

func CloudConsoleURL(prefix string) *url.URL {
	u, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv("ORY_CLOUD_CONSOLE_URL"), "https://console.ory.sh"))
	if err != nil {
		u = &url.URL{Scheme: "https", Host: "console.ory.sh"}
	}
	u.Host = prefix + "." + u.Host
	return u
}

func makeCloudConsoleURL(prefix string) string {
	u := CloudConsoleURL(prefix)

	return u.Scheme + "://" + u.Host
}

func CloudAPIsURL(prefix string) *url.URL {
	u, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv("ORY_CLOUD_ORYAPIS_URL"), "https://oryapis.com"))
	if err != nil {
		u = &url.URL{Scheme: "https", Host: "oryapis.com"}
	}
	u.Host = prefix + "." + u.Host
	return u
}

func makeCloudAPIsURL(prefix string) string {
	u := CloudAPIsURL(prefix)

	return u.Scheme + "://" + u.Host
}

func NewKratosClient() (*cloud.APIClient, error) {
	conf := cloud.NewConfiguration()
	conf.Servers = cloud.ServerConfigurations{{URL: makeCloudConsoleURL("project")}}
	conf.HTTPClient = &http.Client{Timeout: time.Second * 10}
	if RateLimitHeader != "" {
		conf.AddDefaultHeader("Ory-RateLimit-Action", RateLimitHeader)
	}

	return cloud.NewAPIClient(conf), nil
}

func newCloudClient(t *oauth2.Token) *cloud.APIClient {
	u := makeCloudConsoleURL("api")

	conf := cloud.NewConfiguration()
	conf.Servers = cloud.ServerConfigurations{{URL: u}}
	// conf.HTTPClient = oac.Client(context.Background(), t)
	conf.UserAgent = "ory-cli/" + buildinfo.Version
	if RateLimitHeader != "" {
		conf.AddDefaultHeader("Ory-RateLimit-Action", RateLimitHeader)
	}
	conf.Debug = true

	return cloud.NewAPIClient(conf)
}
