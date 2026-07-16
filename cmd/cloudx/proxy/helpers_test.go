// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/proxy"
)

func TestReqMiddleware(t *testing.T) {
	oryURL := &url.URL{Scheme: "https", Host: "example.projects.oryapis.com"}
	publicURL := &url.URL{Scheme: "http", Host: "localhost:4000"}
	const apiKey = "ory_apikey_sentinel"

	const (
		headerToken    = "Ory-Base-URL-Rewrite-Token"
		headerRewrite  = "Ory-Base-URL-Rewrite"
		headerNoCustom = "Ory-No-Custom-Domain-Redirect"
	)

	newRequest := func(t *testing.T, outHost string) *httputil.ProxyRequest {
		t.Helper()
		in, err := http.NewRequest(http.MethodGet, "http://localhost:4000/foo", nil)
		require.NoError(t, err)
		out, err := http.NewRequest(http.MethodGet, "http://"+outHost+"/foo", nil)
		require.NoError(t, err)
		return &httputil.ProxyRequest{In: in, Out: out}
	}

	t.Run("case=ory-bound request receives Ory headers", func(t *testing.T) {
		conf := &config{publicURL: publicURL}
		r := newRequest(t, oryURL.Host)

		_, err := reqMiddleware(conf, oryURL, apiKey)(r, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Equal(t, apiKey, r.Out.Header.Get(headerToken))
		assert.Equal(t, publicURL.String(), r.Out.Header.Get(headerRewrite))
		assert.Equal(t, "true", r.Out.Header.Get(headerNoCustom))
	})

	t.Run("case=upstream-bound request does not receive Ory headers", func(t *testing.T) {
		conf := &config{publicURL: publicURL}
		r := newRequest(t, "localhost:3000")

		_, err := reqMiddleware(conf, oryURL, apiKey)(r, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Empty(t, r.Out.Header.Get(headerToken), "API key must not leak to the upstream app")
		assert.Empty(t, r.Out.Header.Get(headerRewrite))
		assert.Empty(t, r.Out.Header.Get(headerNoCustom))
	})

	t.Run("case=rewriteHost upstream sets X-Forwarded-Host but not the API key", func(t *testing.T) {
		conf := &config{publicURL: publicURL, rewriteHost: true}
		r := newRequest(t, "localhost:3000")

		_, err := reqMiddleware(conf, oryURL, apiKey)(r, &proxy.HostConfig{UpstreamHost: "upstream.internal"}, nil)
		require.NoError(t, err)

		assert.Equal(t, r.In.Host, r.Out.Header.Get("X-Forwarded-Host"))
		assert.Equal(t, "upstream.internal", r.Out.Host)
		assert.Empty(t, r.Out.Header.Get(headerToken), "API key must not leak to the upstream app")
	})
}
