package cloudx

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

type sessionTokenTransporter struct {
	http.RoundTripper
	sessionToken string
}

func (t *sessionTokenTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.sessionToken != "" {
		req.Header.Set("X-Session-Token", "session "+t.sessionToken)
	}
	return t.RoundTripper.RoundTrip(req)
}

func newSessionTokenClient(token string) *http.Client {
	return &http.Client{
		Transport: &sessionTokenTransporter{
			RoundTripper: http.DefaultTransport,
			sessionToken: token,
		},
		Timeout: time.Second * 30,
	}
}
