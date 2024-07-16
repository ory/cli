// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/oauth2"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/stringsx"
)

const (
	RateLimitHeaderKey = "ORY_RATE_LIMIT_HEADER"
	ConsoleURLKey      = "ORY_CONSOLE_URL"
	OryAPIsURLKey      = "ORY_ORYAPIS_URL"
)

var rateLimitHeader = os.Getenv(RateLimitHeaderKey)

func CloudConsoleURL(prefix string) *url.URL {
	// we load the URL from the env here instead of init() because the tests might want to change this
	consoleURL, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv(ConsoleURLKey), "https://console.ory.sh"))
	if err != nil {
		fmt.Printf("error parsing console url: %s\n", err)
		consoleURL = &url.URL{Scheme: "https", Host: "console.ory.sh"}
	}
	if prefix != "" {
		consoleURL.Host = prefix + "." + consoleURL.Host
	}

	return consoleURL
}

func CloudAPIsURL(slug string) *url.URL {
	// we load the URL from the env here instead of init() because the tests might want to change this
	oryAPIsURL, err := url.ParseRequestURI(stringsx.Coalesce(os.Getenv(OryAPIsURLKey), "https://oryapis.com"))
	if err != nil {
		fmt.Printf("error parsing oryapis url: %s\n", err)
		oryAPIsURL = &url.URL{Scheme: "https", Host: "oryapis.com"}
	}
	oryAPIsURL.Host = slug + "." + oryAPIsURL.Host

	return oryAPIsURL
}

func newSDKConfiguration(uri string) *cloud.Configuration {
	conf := cloud.NewConfiguration()
	conf.Servers = cloud.ServerConfigurations{{URL: uri}}
	conf.OperationServers = nil
	conf.HTTPClient = &http.Client{Timeout: time.Second * 30}
	if rateLimitHeader != "" {
		conf.AddDefaultHeader("Ory-RateLimit-Action", rateLimitHeader)
	}
	return conf
}

func NewPublicOryProjectClient() *cloud.APIClient {
	conf := newSDKConfiguration(CloudConsoleURL("project").String())
	return cloud.NewAPIClient(conf)
}

func (h *CommandHelper) newConsoleAPIClient(ctx context.Context) (_ *cloud.APIClient, err error) {
	conf := newSDKConfiguration(CloudConsoleURL("api").String())
	conf.HTTPClient, err = h.newConsoleHTTPClient(ctx)
	if err != nil {
		return nil, err
	}
	return cloud.NewAPIClient(conf), nil
}

func (h *CommandHelper) newConsoleHTTPClient(ctx context.Context) (*http.Client, error) {
	// use the workspace API key if set
	if h.workspaceAPIKey != nil {
		return newOAuth2TokenClient(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *h.workspaceAPIKey})), nil
	}

	// fall back to interactive OAuth2 flow
	config, err := h.GetAuthenticatedConfig(ctx)
	if err != nil {
		return nil, err
	}

	return newOAuth2TokenClient(config.TokenSource(ctx)), nil
}

func (h *CommandHelper) ProjectAuthToken(ctx context.Context) (oauth2.TokenSource, func(string) *url.URL, error) {
	if h.projectAPIKey != nil {
		return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *h.projectAPIKey}), CloudAPIsURL, nil
	}

	config, err := h.GetAuthenticatedConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	return config.TokenSource(ctx), CloudAPIsURL, nil
}

func (h *CommandHelper) newProjectHTTPClient(ctx context.Context) (*http.Client, func(string) *url.URL, error) {
	tokenSource, baseURL, err := h.ProjectAuthToken(ctx)
	if err != nil {
		return nil, nil, err
	}

	retryable := retryablehttp.NewClient()
	retryable.Logger = nil
	c := retryable.StandardClient()
	c.Transport = &oauth2.Transport{
		Base:   c.Transport,
		Source: tokenSource,
	}

	return c, baseURL, nil
}
