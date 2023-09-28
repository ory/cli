// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/http"
	"net/url"
	"os"
	"time"

	cloud "github.com/ory/client-go"
	oldCloud "github.com/ory/client-go/114"
	"github.com/ory/x/stringsx"
)

var RateLimitHeader = os.Getenv("ORY_RATE_LIMIT_HEADER")

func CloudConsoleURL(prefix string) *url.URL {
	u, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv("ORY_CLOUD_CONSOLE_URL"), "https://console.ory.sh"))
	if err != nil {
		u = &url.URL{Scheme: "https", Host: "console.ory.sh"}
	}
	u.Host = prefix + "." + u.Host
	if u.Port() == "" {
		u.Host = u.Host + ":443"
	}

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
	if u.Port() == "" {
		u.Host = u.Host + ":443"
	}

	return u
}

func makeCloudAPIsURL(prefix string) string {
	u := CloudAPIsURL(prefix)

	return u.Scheme + "://" + u.Host
}

func NewKratosClient() (*oldCloud.APIClient, error) {
	conf := oldCloud.NewConfiguration()
	conf.Servers = oldCloud.ServerConfigurations{{URL: makeCloudConsoleURL("project")}}
	conf.HTTPClient = &http.Client{Timeout: time.Second * 10}
	if RateLimitHeader != "" {
		conf.AddDefaultHeader("Ory-RateLimit-Action", RateLimitHeader)
	}

	return oldCloud.NewAPIClient(conf), nil
}

func newCloudClient(token string) (*cloud.APIClient, error) {
	u := makeCloudConsoleURL("api")

	conf := cloud.NewConfiguration()
	conf.Servers = cloud.ServerConfigurations{{URL: u}}
	conf.HTTPClient = newBearerTokenClient(token)
	if RateLimitHeader != "" {
		conf.AddDefaultHeader("Ory-RateLimit-Action", RateLimitHeader)
	}

	return cloud.NewAPIClient(conf), nil
}
