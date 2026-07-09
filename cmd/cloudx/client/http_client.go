// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func newOAuth2TokenClient(token oauth2.TokenSource) *http.Client {
	return &http.Client{
		Transport: &oauth2.Transport{
			Base:   http.DefaultTransport,
			Source: token,
		},
		Timeout: time.Second * 30,
	}
}

// setHeaderTransport sets a fixed header on every outgoing request. It is used
// to attach the Ory-RateLimit-Action header (see ORY_RATE_LIMIT_HEADER) to the
// project HTTP client so that admin API calls issued by the wrapped Ory Kratos
// and Ory Hydra CLI commands are not rate limited during E2E tests.
type setHeaderTransport struct {
	base       http.RoundTripper
	key, value string
}

func (t *setHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set(t.key, t.value)
	return t.base.RoundTrip(req)
}
