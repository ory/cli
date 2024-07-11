// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/randx"
	"github.com/ory/x/urlx"
)

func (h *CommandHelper) checkAuthenticated(_ context.Context) error {
	c, err := h.getConfig()
	if err != nil {
		return err
	}
	if c.isAuthenticated {
		return nil
	}
	if c.AccessToken == nil {
		return ErrNotAuthenticated
	}
	// TODO should we do some API call here to check whether the token is still valid?
	// return ErrReauthenticate
	return nil
}

func (h *CommandHelper) GetAuthenticatedConfig(ctx context.Context) (*Config, error) {
	if err := h.checkAuthenticated(ctx); err == nil {
		return h.getConfig()
	} else if stderrors.Is(err, ErrReauthenticate) {
		if h.isQuiet {
			return nil, ErrNoConfigQuiet
		}
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "Your session has expired or has otherwise become invalid. Please re-authenticate to continue.\n")
	} else if stderrors.Is(err, ErrNoConfig) || stderrors.Is(err, ErrNotAuthenticated) {
		if h.isQuiet {
			return nil, ErrNoConfigQuiet
		}
	}
	if err := h.ClearConfig(); err != nil {
		return nil, err
	}

	if err := h.Authenticate(ctx); err != nil {
		return nil, err
	}

	return h.getConfig()
}

func (c *Config) fromUserinfo(info *cloud.OidcUserInfo) error {
	c.IdentityTraits = Identity{}
	if info.Email != nil {
		c.IdentityTraits.Email = *info.Email
	} else {
		return fmt.Errorf("userinfo response did not contain email")
	}
	if info.Name != nil {
		c.IdentityTraits.Name = *info.Name
	} else {
		return fmt.Errorf("userinfo response did not contain name")
	}
	if info.Sub != nil {
		var err error
		c.IdentityTraits.ID, err = uuid.FromString(*info.Sub)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("userinfo response did not contain sub")
	}
	return nil
}

func (h *CommandHelper) Authenticate(ctx context.Context) error {
	if h.isQuiet {
		return stderrors.New("can not sign in or sign up when flag --quiet is set")
	}

	config, err := h.getConfig()
	if stderrors.Is(err, ErrNoConfig) {
		config = &Config{
			location: h.configLocation,
		}
	} else if err != nil {
		return err
	}

	if config.AccessToken != nil {
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "You are already logged in. Use the logout command to log out.\n")
		return nil
	}

	config, err = h.loginOAuth2(ctx)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(h.VerboseErrWriter, "You are now signed in as: %s\n", config.IdentityTraits.Email)
	return nil
}

func (h *CommandHelper) ClearConfig() error {
	return h.UpdateConfig(&Config{
		location: h.configLocation,
	})
}

func oauth2ClientConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: "ory-cli",
		Endpoint: oauth2.Endpoint{
			AuthURL:   urlx.AppendPaths(CloudConsoleURL("project"), "/oauth2/auth").String(),
			TokenURL:  urlx.AppendPaths(CloudConsoleURL("project"), "/oauth2/token").String(),
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

func (h *CommandHelper) loginOAuth2(ctx context.Context) (*Config, error) {
	client := oauth2ClientConfig()
	token, err := h.oAuth2DanceWithServer(ctx, client)
	if err != nil {
		return nil, err
	}

	scope, _ := token.Extra("scope").(string)
	if !slices.Contains(strings.Split(scope, " "), "offline_access") {
		_, _ = fmt.Fprintf(h.VerboseErrWriter,
			"You have not granted the 'offline_access' permission during login and will have to authenticate again in %s.\n",
			time.Until(token.Expiry).Round(time.Second),
		)
	}

	config := &Config{
		AccessToken: token,
		location:    h.configLocation,
	}
	cl := NewPublicOryProjectClient()
	userInfo, _, err := cl.OidcAPI.GetOidcUserInfo(context.WithValue(ctx, cloud.ContextOAuth2, config.TokenSource(ctx))).Execute()
	if err != nil {
		return nil, err
	}

	if err := config.fromUserinfo(userInfo); err != nil {
		return nil, err
	}

	if err := h.UpdateConfig(config); err != nil {
		return nil, err
	}

	_, _ = fmt.Fprintln(h.VerboseErrWriter, "Successfully logged into Ory Network.")
	return config, nil
}

func (h *CommandHelper) oAuth2DanceWithServer(ctx context.Context, client *oauth2.Config) (token *oauth2.Token, err error) {
	var (
		l            net.Listener
		state        = randx.MustString(32, randx.AlphaNum)
		pkceVerifier = oauth2.GenerateVerifier()
		ports        = []int{12345, 34525, 49763, 51238, 59724, 60582, 62125}
	)
	rand.Shuffle(len(ports), func(i, j int) { ports[i], ports[j] = ports[j], ports[i] })
	for _, port := range ports {
		l, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			client.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)
			break
		}
	}
	if l == nil {
		return nil, fmt.Errorf("failed to allocate port for OAuth2 callback handler, try again later: last error: %w", err)
	}

	var (
		serverErr   = make(chan error)
		serverToken = make(chan *oauth2.Token)
	)
	srv := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// for retries the user has to start from the beginning
			defer close(serverErr)
			defer close(serverToken)

			ctx := r.Context()
			if err := r.ParseForm(); err != nil {
				redirectErr(w, r, "parse form", "An error occurred during CLI authentication. Please try again")
				serverErr <- fmt.Errorf("failed to parse form: %w", err)
				return
			}
			if s := r.Form.Get("state"); s != state {
				redirectErr(w, r, "state mismatch", "An error occurred during CLI authentication. Please try again")
				serverErr <- fmt.Errorf("state mismatch: expected %q, got %q", state, s)
				return
			}
			if r.Form.Has("error") {
				e, d := r.Form.Get("error"), r.Form.Get("error_description")
				redirectErr(w, r, e, d)
				serverErr <- fmt.Errorf("upsteam error: %s: %s", e, d)
				return
			}
			code := r.Form.Get("code")
			if code == "" {
				redirectErr(w, r, "missing code", "An error occurred during CLI authentication. Please try again")
				serverErr <- fmt.Errorf("missing code")
				return
			}
			t, err := client.Exchange(
				ctx,
				code,
				oauth2.VerifierOption(pkceVerifier),
			)
			if err != nil {
				redirectErr(w, r, "token exchange", "An error occurred during the OAuth2 token exchange")
				serverErr <- fmt.Errorf("failed OAuth2 token exchange: %w", err)
				return
			}
			serverToken <- t
			redirectOK(w, r)
		}),
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve OAuth2 callback handler: %w", err)
		}
		return nil
	})
	eg.Go(func() (err error) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case token = <-serverToken:
		case err = <-serverErr:
		}
		ctx, cancel := context.WithDeadline(context.WithoutCancel(ctx), time.Now().Add(2*time.Second))
		defer cancel()
		return stderrors.Join(err, srv.Shutdown(ctx))
	})

	u := client.AuthCodeURL(state,
		oauth2.S256ChallengeOption(pkceVerifier),
		oauth2.SetAuthURLParam("scope", "offline_access email profile"),
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("prompt", "login consent"),
		oauth2.SetAuthURLParam("audience", CloudConsoleURL("api").String()),
	)
	if err := h.openBrowserHook(u); err != nil {
		return nil, err
	}
	_, _ = fmt.Fprintf(h.VerboseErrWriter,
		`A browser should have opened for you to complete your login to Ory Network.
If no browser opened, visit the below page to continue:

		%s

`, u)

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to authenticate, please try again: %w", err)
	}
	return token, nil
}

func redirectOK(w http.ResponseWriter, r *http.Request) {
	location := CloudConsoleURL("")
	location.Path = "/cli-auth-success"
	http.Redirect(w, r, location.String(), http.StatusFound)
}

func redirectErr(w http.ResponseWriter, r *http.Request, err, desc string) {
	location := CloudConsoleURL("")
	location.Path = "/error"
	location.RawQuery = url.Values{"error": []string{err}, "error_description": []string{desc}}.Encode()
	http.Redirect(w, r, location.String(), http.StatusFound)
}

func (h *CommandHelper) SignOut(ctx context.Context) error {
	config, err := h.GetAuthenticatedConfig(ctx)
	if err != nil {
		return err
	}
	if config.AccessToken == nil {
		return h.ClearConfig()
	}
	client := oauth2ClientConfig()
	revoke, err := url.Parse(client.Endpoint.AuthURL)
	if err != nil {
		return err
	}
	revoke.Path = "/oauth2/revoke"
	res, err := http.PostForm(revoke.String(), url.Values{
		"client_id": []string{client.ClientID},
		"token":     []string{config.AccessToken.RefreshToken}, // this also revokes the associated access token
	})
	if err != nil {
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "failed to revoke access token: %v\n", err)
	} else {
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			body, _ := io.ReadAll(res.Body)
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "failed to revoke access token: %v\n", string(body))
		}
	}
	return h.ClearConfig()
}
