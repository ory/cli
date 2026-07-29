// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/proxy"
)

func TestRespMiddleware(t *testing.T) {
	const (
		headerAllowOrigin      = "Access-Control-Allow-Origin"
		headerAllowCredentials = "Access-Control-Allow-Credentials"
		headerExposeHeaders    = "Access-Control-Expose-Headers"
	)

	newResponse := func(header http.Header) *http.Response {
		return &http.Response{StatusCode: http.StatusOK, Header: header}
	}

	t.Run("case=strips CORS headers the proxy sets itself", func(t *testing.T) {
		resp := newResponse(http.Header{
			headerAllowOrigin:      []string{"http://localhost:3000"},
			headerAllowCredentials: []string{"true"},
		})

		_, err := respMiddleware(&config{})(resp, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Empty(t, resp.Header.Values(headerAllowOrigin))
		assert.Empty(t, resp.Header.Values(headerAllowCredentials))
	})

	t.Run("case=keeps Access-Control-Expose-Headers from the upstream", func(t *testing.T) {
		resp := newResponse(http.Header{headerExposeHeaders: []string{"X-Custom"}})

		_, err := respMiddleware(&config{})(resp, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Equal(t, []string{"X-Custom"}, resp.Header.Values(headerExposeHeaders),
			"browsers merge duplicates of this header, so stripping it would only lose information")
	})

	t.Run("case=rewrites the welcome page redirect", func(t *testing.T) {
		conf := &config{pathPrefix: "/.ory", defaultRedirectTo: cmdx.URL{URL: url.URL{Scheme: "http", Host: "localhost:3000"}}}
		resp := newResponse(http.Header{"Location": []string{"http://localhost:4000/.ory/ui/welcome"}})
		resp.StatusCode = http.StatusFound

		_, err := respMiddleware(conf)(resp, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Equal(t, "http://localhost:3000", resp.Header.Get("Location"))
	})

	t.Run("case=leaves other redirects alone", func(t *testing.T) {
		conf := &config{pathPrefix: "/.ory", defaultRedirectTo: cmdx.URL{URL: url.URL{Scheme: "http", Host: "localhost:3000"}}}
		resp := newResponse(http.Header{"Location": []string{"http://localhost:4000/.ory/ui/login"}})
		resp.StatusCode = http.StatusFound

		_, err := respMiddleware(conf)(resp, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Equal(t, "http://localhost:4000/.ory/ui/login", resp.Header.Get("Location"))
	})
}

// TestCORSHeadersAreNotDuplicated is the regression test for
// https://github.com/ory/cli/issues/344. An upstream that handles CORS itself
// used to produce two Access-Control-Allow-Origin headers, which browsers
// reject outright even though both values are correct.
func TestCORSHeadersAreNotDuplicated(t *testing.T) {
	const origin = "http://localhost:3000"

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		_, _ = io.WriteString(w, "ok")
	}))
	t.Cleanup(upstream.Close)

	upstreamURL, err := url.Parse(upstream.URL)
	require.NoError(t, err)

	conf := &config{publicURL: &url.URL{Scheme: "http", Host: "localhost:4000"}}

	rp := httputil.NewSingleHostReverseProxy(upstreamURL)
	rp.ModifyResponse = func(resp *http.Response) error {
		_, err := respMiddleware(conf)(resp, &proxy.HostConfig{}, nil)
		return err
	}

	// Built from the same helper runReverseProxy uses, so the test cannot drift
	// away from the CORS configuration the proxy actually runs.
	corsOpts, err := corsOptions(conf)
	require.NoError(t, err)
	ch := cors.New(corsOpts)

	srv := httptest.NewServer(ch.Handler(rp))
	t.Cleanup(srv.Close)

	// HEAD is CORS-safelisted, so it reaches the upstream without a preflight and
	// must still come back with the proxy's own CORS headers.
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		t.Run("method="+method, func(t *testing.T) {
			req, err := http.NewRequest(method, srv.URL+"/api", nil)
			require.NoError(t, err)
			req.Header.Set("Origin", origin)

			res, err := srv.Client().Do(req)
			require.NoError(t, err)
			t.Cleanup(func() { _ = res.Body.Close() })

			assert.Equal(t, []string{origin}, res.Header.Values("Access-Control-Allow-Origin"),
				"a duplicated Access-Control-Allow-Origin makes the browser reject the response, and none at all blocks it too")
			assert.Equal(t, []string{"true"}, res.Header.Values("Access-Control-Allow-Credentials"))
		})
	}
}

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

	t.Run("case=client-spoofed Ory headers are stripped from upstream-bound requests", func(t *testing.T) {
		conf := &config{publicURL: publicURL}
		r := newRequest(t, "localhost:3000")
		r.Out.Header.Set(headerToken, "ory_apikey_spoofed")
		r.Out.Header.Set(headerRewrite, "http://evil.example")
		r.Out.Header.Set(headerNoCustom, "true")

		_, err := reqMiddleware(conf, oryURL, apiKey)(r, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Empty(t, r.Out.Header.Get(headerToken), "spoofed token must not be forwarded to the upstream app")
		assert.Empty(t, r.Out.Header.Get(headerRewrite), "spoofed header must not be forwarded to the upstream app")
		assert.Empty(t, r.Out.Header.Get(headerNoCustom), "spoofed header must not be forwarded to the upstream app")
	})

	t.Run("case=client-spoofed token is stripped from Ory-bound requests when apiKey is empty", func(t *testing.T) {
		conf := &config{publicURL: publicURL}
		r := newRequest(t, oryURL.Host)
		r.Out.Header.Set(headerToken, "ory_apikey_spoofed")

		_, err := reqMiddleware(conf, oryURL, "")(r, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Empty(t, r.Out.Header.Get(headerToken), "spoofed token must not be passed through to Ory when apiKey is empty")
	})

	t.Run("case=real apiKey overwrites a client-spoofed token on Ory-bound requests", func(t *testing.T) {
		conf := &config{publicURL: publicURL}
		r := newRequest(t, oryURL.Host)
		r.Out.Header.Set(headerToken, "ory_apikey_spoofed")

		_, err := reqMiddleware(conf, oryURL, apiKey)(r, &proxy.HostConfig{}, nil)
		require.NoError(t, err)

		assert.Equal(t, apiKey, r.Out.Header.Get(headerToken), "genuine key must overwrite any spoofed token")
	})
}
