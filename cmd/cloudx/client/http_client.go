// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"golang.org/x/oauth2"
	"net/http"
	"time"
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
