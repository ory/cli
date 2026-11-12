// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

const linkPreviewUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/601.2.4 (KHTML, like Gecko) Version/9.0.1 Safari/601.2.4 facebookexternalhit/1.1 Facebot Twitterbot/1.0"

type fakeAuthorizationHooks struct {
	// Each hook runs while the browser waits for the answer to its request, so
	// whatever the hook sends reaches the CLI before the browser continues.
	beforeCallback, duringRedemption func(callbackURL string)
}

// fakeAuthorizationServer signs in every browser at once: it answers an
// authorization request with a redirect to the client that carries a code, or
// an error if access is denied. Link previews cannot sign in and get the login
// page instead. Like Ory Network, it revokes the token issued for a code that
// is redeemed twice. It also serves the Console pages the CLI redirects the
// browser to.
type fakeAuthorizationServer struct {
	t     *testing.T
	URL   string
	deny  bool
	hooks fakeAuthorizationHooks

	mu sync.Mutex
	// flows are keyed by their code.
	flows   map[string]*fakeFlow
	revoked bool
}

type fakeFlow struct {
	authorizeRequest url.Values
	callbackURL      string
	redemptions      int
}

func newFakeAuthorizationServer(t *testing.T, deny bool, hooks fakeAuthorizationHooks) *fakeAuthorizationServer {
	f := &fakeAuthorizationServer{t: t, deny: deny, hooks: hooks, flows: map[string]*fakeFlow{}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /oauth2/auth", f.authorize)
	mux.HandleFunc("POST /oauth2/token", f.token)
	mux.Handle("GET /cli-auth-success", fakePage("Signed in"))
	mux.Handle("GET /error", fakePage("Error"))

	srv := httptest.NewUnstartedServer(mux)
	// Browsers treat localhost and 127.0.0.1 as different sites, as they do
	// Ory Network and the CLI's callback server.
	f.URL = fmt.Sprintf("http://localhost:%d", srv.Listener.Addr().(*net.TCPAddr).Port)
	srv.Start()
	t.Cleanup(srv.Close)
	return f
}

func fakePage(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprintf(w, "<!doctype html><html><body>%s</body></html>", body)
	}
}

func (f *fakeAuthorizationServer) authorize(w http.ResponseWriter, r *http.Request) {
	if r.UserAgent() == linkPreviewUserAgent {
		fakePage("Sign in")(w, r)
		return
	}
	q := r.URL.Query()
	assert.Equal(f.t, "ory-cli", q.Get("client_id"))
	assert.Equal(f.t, "code", q.Get("response_type"))
	assert.Equal(f.t, "S256", q.Get("code_challenge_method"))
	assert.NotEmpty(f.t, q.Get("code_challenge"))
	assert.NotEmpty(f.t, q.Get("state"))

	params := url.Values{"state": {q.Get("state")}}
	if f.deny {
		params.Set("error", "access_denied")
		params.Set("error_description", "The resource owner denied the request")
	} else {
		params.Set("code", rand.Text())
	}
	callback := q.Get("redirect_uri") + "?" + params.Encode()
	if !f.deny {
		f.mu.Lock()
		f.flows[params.Get("code")] = &fakeFlow{authorizeRequest: q, callbackURL: callback}
		f.mu.Unlock()
	}

	if f.hooks.beforeCallback != nil {
		f.hooks.beforeCallback(callback)
	}
	http.Redirect(w, r, callback, http.StatusSeeOther)
}

func (f *fakeAuthorizationServer) token(w http.ResponseWriter, r *http.Request) {
	tokenError := func(err, description string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err, "error_description": description})
	}
	if !assert.NoError(f.t, r.ParseForm()) {
		tokenError("invalid_request", "unparsable form")
		return
	}

	f.mu.Lock()
	flow, ok := f.flows[r.PostForm.Get("code")]
	valid := ok &&
		r.PostForm.Get("grant_type") == "authorization_code" &&
		r.PostForm.Get("client_id") == flow.authorizeRequest.Get("client_id") &&
		r.PostForm.Get("redirect_uri") == flow.authorizeRequest.Get("redirect_uri") &&
		oauth2.S256ChallengeFromVerifier(r.PostForm.Get("code_verifier")) == flow.authorizeRequest.Get("code_challenge")
	var reused bool
	if valid {
		flow.redemptions++
		reused = flow.redemptions > 1
		f.revoked = f.revoked || reused
	}
	f.mu.Unlock()

	switch {
	case !valid:
		tokenError("invalid_grant", "The provided authorization grant is invalid.")
		return
	case reused:
		tokenError("invalid_grant", "The authorization code has already been used.")
		return
	}
	if f.hooks.duringRedemption != nil {
		f.hooks.duringRedemption(flow.callbackURL)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"access_token":  "access-token",
		"refresh_token": "refresh-token",
		"token_type":    "bearer",
		"expires_in":    3600,
		"scope":         "offline_access email profile",
	})
}

func (f *fakeAuthorizationServer) redeemed() (redemptions int, revoked bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, flow := range f.flows {
		redemptions += flow.redemptions
	}
	return redemptions, f.revoked
}

// errorCallback returns the callback that Ory Network sends for the same state
// when a link preview repeats one of the browser's requests to it.
func errorCallback(t *testing.T, callbackURL string) string {
	u, err := url.Parse(callbackURL)
	if !assert.NoError(t, err) {
		return ""
	}
	u.RawQuery = url.Values{
		"error":             {"request_forbidden"},
		"error_description": {"The request is not allowed. No CSRF value available in the session cookie."},
		"state":             {u.Query().Get("state")},
	}.Encode()
	return u.String()
}

// linkPreview fetches a URL the way the link previews of macOS do: like the
// browser, but without its cookies and following redirects. A jar lets it keep
// the cookies it receives itself.
func linkPreview(t *testing.T, jar http.CookieJar) func(string) {
	c := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	return func(u string) {
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			t.Errorf("link preview request: %v", err)
			return
		}
		req.Header.Set("User-Agent", linkPreviewUserAgent)
		res, err := c.Do(req)
		if err != nil {
			t.Logf("link preview of %s: %v", u, err)
			return
		}
		t.Logf("link preview of %s ended at %s with %d", u, res.Request.URL, res.StatusCode)
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}
}

// signInWithBrowser opens the URL the CLI gives the browser and waits for the
// page the CLI sends the browser to at the end.
func signInWithBrowser(page playwright.Page, startURL string, deny bool) error {
	if _, err := page.Goto(startURL); err != nil {
		return err
	}
	landing, landingText := "/cli-auth-success", "Signed in"
	if deny {
		landing, landingText = "/error", "Error"
	}
	// Page.WaitForURL can miss a load event that fires while it starts
	// waiting, so this waits for the content of the page instead.
	if err := page.GetByText(landingText, playwright.PageGetByTextOptions{Exact: new(true)}).WaitFor(); err != nil {
		return fmt.Errorf("the browser did not reach %s but stayed at %s: %w", landing, page.URL(), err)
	}
	if u, err := url.Parse(page.URL()); err != nil || u.Path != landing {
		return fmt.Errorf("the browser is at %s instead of %s", page.URL(), landing)
	}
	return nil
}

// browserLoginCookies returns the cookies the CLI set to recognize the browser.
func browserLoginCookies(browserContext playwright.BrowserContext) ([]playwright.Cookie, error) {
	cookies, err := browserContext.Cookies()
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(cookies, func(c playwright.Cookie) bool {
		return !strings.HasPrefix(c.Name, "ory_cli_login_")
	}), nil
}

func TestOAuth2LoginInBrowser(t *testing.T) {
	pw, err := playwright.Run()
	require.NoError(t, err)
	t.Cleanup(func() { _ = pw.Stop() })
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: new(true)})
	require.NoError(t, err)
	t.Cleanup(func() { _ = browser.Close() })

	type setup struct {
		hooks fakeAuthorizationHooks
		// beforeBrowser runs before the browser opens the start URL.
		beforeBrowser       func(startURL string)
		wrapLoopbackHandler func(http.Handler) http.Handler
		check               func(t *testing.T)
	}

	for _, tc := range []struct {
		name  string
		deny  bool
		setup func(t *testing.T, browserContext playwright.BrowserContext) setup
	}{
		{
			name: "signs in with a session cookie",
			setup: func(t *testing.T, browserContext playwright.BrowserContext) setup {
				duringLogin := make(chan []playwright.Cookie, 1)
				return setup{
					hooks: fakeAuthorizationHooks{beforeCallback: func(string) {
						cookies, err := browserLoginCookies(browserContext)
						assert.NoError(t, err)
						duringLogin <- cookies
					}},
					check: func(t *testing.T) {
						var cookies []playwright.Cookie
						select {
						case cookies = <-duringLogin:
						default:
							require.Fail(t, "the cookies were not read during the login")
						}
						require.Len(t, cookies, 1)
						assert.EqualValues(t, -1, cookies[0].Expires, "the cookie must not outlive the browser session")
						assert.Equal(t, "/callback", cookies[0].Path)
						assert.True(t, cookies[0].HttpOnly)
						assert.Equal(t, playwright.SameSiteAttributeLax, cookies[0].SameSite)
					},
				}
			},
		},
		{
			name: "fails when the user denies access",
			deny: true,
		},
		{
			name: "ignores an error callback caused by a link preview",
			setup: func(t *testing.T, _ playwright.BrowserContext) setup {
				preview := linkPreview(t, nil)
				return setup{hooks: fakeAuthorizationHooks{beforeCallback: func(callback string) {
					preview(errorCallback(t, callback))
				}}}
			},
		},
		{
			name: "ignores a link preview of the callback",
			setup: func(t *testing.T, _ playwright.BrowserContext) setup {
				return setup{hooks: fakeAuthorizationHooks{duringRedemption: linkPreview(t, nil)}}
			},
		},
		{
			name: "signs in after a link preview opened the login link first",
			setup: func(t *testing.T, _ playwright.BrowserContext) setup {
				return setup{beforeBrowser: linkPreview(t, nil)}
			},
		},
		{
			name: "refuses to begin a login without the nonce",
			setup: func(t *testing.T, _ playwright.BrowserContext) setup {
				return setup{beforeBrowser: func(startURL string) {
					u, err := url.Parse(startURL)
					require.NoError(t, err)
					u.RawQuery = ""
					res, err := http.Get(u.String())
					require.NoError(t, err)
					_ = res.Body.Close()
					assert.Equal(t, http.StatusNotFound, res.StatusCode)
				}}
			},
		},
		{
			name: "ignores a link preview with the cookie of its own login attempt",
			setup: func(t *testing.T, _ playwright.BrowserContext) setup {
				jar, err := cookiejar.New(nil)
				require.NoError(t, err)
				preview := linkPreview(t, jar)
				return setup{
					beforeBrowser: preview,
					hooks: fakeAuthorizationHooks{beforeCallback: func(callback string) {
						preview(errorCallback(t, callback))
					}},
				}
			},
		},
		{
			name: "ignores the callback of an older tab in the same browser",
			setup: func(t *testing.T, browserContext playwright.BrowserContext) setup {
				return setup{beforeBrowser: func(startURL string) {
					// The request API of a browser context shares its cookies,
					// as the tabs of a browser do.
					tabs := browserContext.Request()
					older, err := tabs.Get(startURL, playwright.APIRequestContextGetOptions{MaxRedirects: new(0)})
					require.NoError(t, err)
					_, err = tabs.Get(startURL, playwright.APIRequestContextGetOptions{MaxRedirects: new(0)})
					require.NoError(t, err)
					res, err := tabs.Get(older.Headers()["location"])
					require.NoError(t, err)
					assert.Equal(t, http.StatusForbidden, res.Status(), "the older tab must not end the login")
				}}
			},
		},
		{
			name: "redeems the code once when the browser repeats the callback",
			setup: func(t *testing.T, browserContext playwright.BrowserContext) setup {
				type navigation struct {
					landed string
					err    error
				}
				repeated := make(chan navigation, 1)
				callbacks := make(chan struct{}, 2)
				return setup{
					wrapLoopbackHandler: func(next http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if r.URL.Path == "/callback" {
								callbacks <- struct{}{}
							}
							next.ServeHTTP(w, r)
						})
					},
					hooks: fakeAuthorizationHooks{duringRedemption: func(callback string) {
						// The browser's first callback started this redemption.
						<-callbacks
						page, err := browserContext.NewPage()
						if !assert.NoError(t, err) {
							return
						}
						// Chromium's HTTP cache holds a request until an earlier
						// one for the same URL is answered; routing it bypasses
						// the cache, so both reach the CLI at the same time.
						err = page.Route(func(u string) bool { return u == callback }, func(r playwright.Route) {
							_ = r.Continue()
						})
						if !assert.NoError(t, err) {
							return
						}
						go func() {
							_, err := page.Goto(callback)
							repeated <- navigation{landed: page.URL(), err: err}
						}()
						select {
						case <-callbacks:
						case <-time.After(10 * time.Second):
							t.Error("the repeated callback did not reach the CLI")
						}
					}},
					check: func(t *testing.T) {
						select {
						case n := <-repeated:
							assert.NoError(t, n.err)
							assert.Regexp(t, `/cli-auth-success$`, n.landed)
						case <-time.After(10 * time.Second):
							t.Error("the repeated callback was not answered")
						}
					},
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			browserContext, err := browser.NewContext()
			require.NoError(t, err)
			t.Cleanup(func() { _ = browserContext.Close() })
			page, err := browserContext.NewPage()
			require.NoError(t, err)
			page.SetDefaultTimeout(10_000)

			var s setup
			if tc.setup != nil {
				s = tc.setup(t, browserContext)
			}
			authorizationServer := newFakeAuthorizationServer(t, tc.deny, s.hooks)
			t.Setenv(ConsoleURLKey, authorizationServer.URL)

			signedIn := make(chan error, 1)
			h := &CommandHelper{
				VerboseErrWriter: io.Discard,
				openBrowserHook: func(startURL string) error {
					if s.beforeBrowser != nil {
						s.beforeBrowser(startURL)
					}
					go func() { signedIn <- signInWithBrowser(page, startURL, tc.deny) }()
					return nil
				},
				wrapLoopbackHandler: s.wrapLoopbackHandler,
			}
			conf := oauth2ClientConfig()
			conf.Endpoint.AuthURL = authorizationServer.URL + "/oauth2/auth"
			conf.Endpoint.TokenURL = authorizationServer.URL + "/oauth2/token"

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			token, err := h.oAuth2DanceWithServer(ctx, conf)

			redemptions, revoked := authorizationServer.redeemed()
			if tc.deny {
				assert.ErrorContains(t, err, "access_denied")
				assert.Nil(t, token)
				assert.Zero(t, redemptions)
			} else if assert.NoError(t, err) {
				assert.Equal(t, "access-token", token.AccessToken)
				assert.Equal(t, 1, redemptions, "the code must be redeemed exactly once")
				assert.False(t, revoked, "the issued token must not be revoked")
			}
			select {
			case err := <-signedIn:
				assert.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Error("the browser did not finish the login")
			}
			if s.check != nil {
				s.check(t)
			}

			remaining, err := browserLoginCookies(browserContext)
			require.NoError(t, err)
			assert.Empty(t, remaining, "the callback must clear the cookie")
		})
	}
}
