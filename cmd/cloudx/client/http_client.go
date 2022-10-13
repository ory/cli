// Copyright © 2022 Ory Corp

package client

import (
	"net/http"
	"time"
)

type bearerTokenTransporter struct {
	http.RoundTripper
	bearerToken string
}

func (t *bearerTokenTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+t.bearerToken)
	}
	return t.RoundTripper.RoundTrip(req)
}

func newBearerTokenClient(token string) *http.Client {
	return &http.Client{
		Transport: &bearerTokenTransporter{
			RoundTripper: http.DefaultTransport,
			bearerToken:  token,
		},
		Timeout: time.Second * 30,
	}
}
