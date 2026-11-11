// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/randx"
	"github.com/ory/x/urlx"
)

func (h *CommandHelper) checkAuthenticated(ctx context.Context) error {
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
	if c.AccessToken.Expiry.Before(time.Now()) {
		// Do a simple call that refreshes the token. If it fails, the user should re-authenticate.
		_, _, err := NewPublicOryProjectClient().OidcAPI.GetOidcUserInfo(context.WithValue(ctx, cloud.ContextOAuth2, c.TokenSource(ctx))).Execute()
		if err != nil {
			return ErrReauthenticate
		}
	}
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
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "Using API key instead of a user session.\n")
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
	userInfo, _, err := NewPublicOryProjectClient().OidcAPI.GetOidcUserInfo(context.WithValue(ctx, cloud.ContextOAuth2, config.TokenSource(ctx))).Execute()
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

func (h *CommandHelper) oAuth2DanceWithServer(ctx context.Context, client *oauth2.Config) (*oauth2.Token, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to allocate port for OAuth2 callback handler, try again later: %w", err)
	}
	client.RedirectURL = fmt.Sprintf("http://%s/callback", l.Addr())

	login := &loopbackLogin{
		// Cookies are not separated by port, so concurrent logins need distinct
		// names.
		cookieName: fmt.Sprintf("ory_cli_login_%d", l.Addr().(*net.TCPAddr).Port),
		authURL: func(state, pkceVerifier string) string {
			return client.AuthCodeURL(state,
				oauth2.S256ChallengeOption(pkceVerifier),
				oauth2.SetAuthURLParam("scope", "offline_access email profile"),
				oauth2.SetAuthURLParam("response_type", "code"),
				oauth2.SetAuthURLParam("prompt", "login consent"),
				oauth2.SetAuthURLParam("audience", CloudConsoleURL("api").String()),
			)
		},
		exchange: func(code, pkceVerifier string) (*oauth2.Token, error) {
			return client.Exchange(ctx, code, oauth2.VerifierOption(pkceVerifier))
		},
		pkceVerifiers: make(map[string]string),
		done:          make(chan struct{}),
	}

	handler := login.handler()
	if h.wrapLoopbackHandler != nil {
		handler = h.wrapLoopbackHandler(handler)
	}
	srv := http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	go func() { _ = srv.Serve(l) }()
	defer srv.Close()
	defer func() {
		// Shutdown lets the browser receive the redirect it waits for.
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	startURL := fmt.Sprintf("http://%s/start", l.Addr())
	if h.noBrowser {
		_, _ = fmt.Fprintln(h.VerboseErrWriter, "Open the below page in a browser on this computer to complete your login to Ory Network:")
	} else {
		if err := h.openBrowserHook(startURL); err != nil {
			return nil, err
		}
		_, _ = fmt.Fprintln(h.VerboseErrWriter, `A browser should have opened for you to complete your login to Ory Network.
If no browser opened, visit the below page to continue:`)
	}
	_, _ = fmt.Fprintf(h.VerboseErrWriter, "\n\t\t%s\n\n", startURL)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-login.done:
		return login.result.token, login.result.err
	}
}

// loopbackLogin is the CLI's side of the browser login: it sends the browser
// to the authorization server and receives the callback with the result.
//
// Every request to the start URL begins an attempt with its own state and PKCE
// verifier, and sets a cookie holding the state. A callback completes the
// login only if its state matches the cookie and belongs to an attempt. Other
// programs on the machine, such as link previews that fetch the pages the
// browser visits, may request the start and callback URLs too. They get
// attempts of their own or are rejected, and leave the login running.
type loopbackLogin struct {
	// All attempts share the cookie name, so in each browser only the latest
	// attempt can complete, as with the authorization server's CSRF cookie.
	cookieName string
	authURL    func(state, pkceVerifier string) string
	exchange   func(code, pkceVerifier string) (*oauth2.Token, error)

	mu sync.Mutex
	// pkceVerifiers are keyed by state.
	pkceVerifiers map[string]string

	once sync.Once
	// done is closed after the first callback from a browser has set result.
	done   chan struct{}
	result loopbackLoginResult
}

const maxLoginAttempts = 32

type loopbackLoginResult struct {
	token *oauth2.Token
	err   error
	// browserError and browserErrorDescription are shown in the browser if
	// err is set.
	browserError, browserErrorDescription string
}

func (l *loopbackLogin) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /start", l.start)
	mux.HandleFunc("GET /callback", l.callback)
	return mux
}

func (l *loopbackLogin) start(w http.ResponseWriter, r *http.Request) {
	state := randx.MustString(32, randx.AlphaNum)
	pkceVerifier := oauth2.GenerateVerifier()
	l.mu.Lock()
	full := len(l.pkceVerifiers) >= maxLoginAttempts
	if !full {
		l.pkceVerifiers[state] = pkceVerifier
	}
	l.mu.Unlock()
	if full {
		http.Error(w, "This login was opened too many times. Run the command again.", http.StatusTooManyRequests)
		return
	}

	http.SetCookie(w, l.cookie(state, 0))
	http.Redirect(w, r, l.authURL(state, pkceVerifier), http.StatusFound)
}

func (l *loopbackLogin) callback(w http.ResponseWriter, r *http.Request) {
	pkceVerifier, ok := l.pkceVerifierOf(r)
	if !ok {
		http.Error(w, "This login was started in another browser window. Continue in the window where you last opened the login link, or run the command again.", http.StatusForbidden)
		return
	}
	http.SetCookie(w, l.cookie("", -1))

	// A repeated callback waits for the first one and gets the same answer.
	l.once.Do(func() {
		l.result = l.complete(r, pkceVerifier)
		close(l.done)
	})
	if l.result.err != nil {
		redirectErr(w, r, l.result.browserError, l.result.browserErrorDescription)
		return
	}
	redirectOK(w, r)
}

func (l *loopbackLogin) cookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     l.cookieName,
		Value:    value,
		Path:     "/callback",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (l *loopbackLogin) pkceVerifierOf(r *http.Request) (string, bool) {
	state := r.URL.Query().Get("state")
	if c, err := r.Cookie(l.cookieName); err != nil || c.Value != state {
		return "", false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	pkceVerifier, ok := l.pkceVerifiers[state]
	return pkceVerifier, ok
}

func (l *loopbackLogin) complete(r *http.Request, pkceVerifier string) loopbackLoginResult {
	const tryAgain = "An error occurred during CLI authentication. Please try again"
	if err := r.ParseForm(); err != nil {
		return loopbackLoginResult{err: fmt.Errorf("failed to parse form: %w", err), browserError: "parse form", browserErrorDescription: tryAgain}
	}
	if r.Form.Has("error") {
		e, d := r.Form.Get("error"), r.Form.Get("error_description")
		return loopbackLoginResult{err: fmt.Errorf("upstream error: %s: %s", e, d), browserError: e, browserErrorDescription: d}
	}
	code := r.Form.Get("code")
	if code == "" {
		return loopbackLoginResult{err: stderrors.New("missing code"), browserError: "missing code", browserErrorDescription: tryAgain}
	}
	token, err := l.exchange(code, pkceVerifier)
	if err != nil {
		return loopbackLoginResult{err: fmt.Errorf("failed OAuth2 token exchange: %w", err), browserError: "token exchange", browserErrorDescription: "An error occurred during the OAuth2 token exchange"}
	}
	return loopbackLoginResult{token: token}
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

// SignOut revokes the stored OAuth2 token and removes it from the local
// configuration. It reads the configuration file directly instead of going
// through GetAuthenticatedConfig, because signing out must not require a valid
// session: an expired or invalid one would otherwise start a new login, and
// RFC 7009 servers accept revocation requests for tokens that are already
// invalid.
func (h *CommandHelper) SignOut(ctx context.Context) error {
	config, err := readConfig(h.configLocation)
	if stderrors.Is(err, ErrNoConfig) {
		return nil
	} else if err != nil {
		return err
	}
	if config.AccessToken != nil {
		if err := revokeToken(ctx, config.AccessToken); err != nil {
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "failed to revoke access token, signing out locally only: %v\n", err)
		}
	}
	return h.ClearConfig()
}

// revokeToken revokes the refresh token, which also revokes the access tokens
// issued with it. Tokens obtained without the offline_access scope have no
// refresh token, so the access token is revoked instead.
func revokeToken(ctx context.Context, token *oauth2.Token) error {
	value, hint := token.RefreshToken, "refresh_token"
	if value == "" {
		value, hint = token.AccessToken, "access_token"
	}
	if value == "" {
		return nil
	}

	client := oauth2ClientConfig()
	revoke, err := url.Parse(client.Endpoint.AuthURL)
	if err != nil {
		return err
	}
	revoke.Path = "/oauth2/revoke"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, revoke.String(), strings.NewReader(url.Values{
		"client_id":       {client.ClientID},
		"token":           {value},
		"token_type_hint": {hint},
	}.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Honor the same context HTTP client that the oauth2 package uses for the
	// token exchange and refresh requests against this server.
	httpClient := http.DefaultClient
	if c, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok && c != nil {
		httpClient = c
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, body)
	}
	return nil
}
