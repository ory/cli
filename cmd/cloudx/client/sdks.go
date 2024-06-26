// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/stringsx"
)

const (
	RateLimitHeaderKey = "ORY_RATE_LIMIT_HEADER"
	ConsoleURLKey      = "ORY_CONSOLE_URL"
	OryAPIsURLKey      = "ORY_ORYAPIS_URL"
)

var rateLimitHeader = os.Getenv(RateLimitHeaderKey)

func cloudConsoleURL(prefix string) *url.URL {
	// we load the URL from the env here instead of init() because the tests might want to change this
	consoleURL, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv(ConsoleURLKey), "https://console.ory.sh"))
	if err != nil {
		consoleURL = &url.URL{Scheme: "https", Host: "console.ory.sh"}
	}
	consoleURL.Host = prefix + "." + consoleURL.Host
	if consoleURL.Port() == "" {
		consoleURL.Host = consoleURL.Host + ":443"
	}

	return consoleURL
}

func CloudAPIsURL(slug string) *url.URL {
	// we load the URL from the env here instead of init() because the tests might want to change this
	oryAPIsURL, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv(OryAPIsURLKey), "https://projects.oryapis.com"))
	if err != nil {
		oryAPIsURL = &url.URL{Scheme: "https", Host: "projects.oryapis.com"}
	}
	oryAPIsURL.Host = slug + "." + oryAPIsURL.Host
	if oryAPIsURL.Port() == "" {
		oryAPIsURL.Host = oryAPIsURL.Host + ":443"
	}

	return oryAPIsURL
}

func NewOryProjectClient() (*cloud.APIClient, error) {
	conf := cloud.NewConfiguration()
	conf.Servers = cloud.ServerConfigurations{{URL: cloudConsoleURL("project").String()}}
	conf.HTTPClient = &http.Client{Timeout: time.Second * 30}
	if rateLimitHeader != "" {
		conf.AddDefaultHeader("Ory-RateLimit-Action", rateLimitHeader)
	}

	return cloud.NewAPIClient(conf), nil
}

func (h *CommandHelper) newCloudClient(ctx context.Context) (*cloud.APIClient, error) {
	config, err := h.GetAuthenticatedConfig(ctx)
	if err != nil {
		return nil, err
	}

	conf := cloud.NewConfiguration()
	conf.OperationServers = nil
	conf.Servers = cloud.ServerConfigurations{{URL: cloudConsoleURL("api").String()}}
	conf.HTTPClient = newBearerTokenClient(config.SessionToken)
	if rateLimitHeader != "" {
		conf.AddDefaultHeader("Ory-RateLimit-Action", rateLimitHeader)
	}

	return cloud.NewAPIClient(conf), nil
}
