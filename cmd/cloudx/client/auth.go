package client

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid/v3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/randx"
)

func RegisterAuthHelpers(cmd *cobra.Command) {
	var (
		h  *CommandHelper
		ac *AuthContext
	)
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		h, err = NewCommandHelper(cmd)
		if err != nil {
			return err
		}
		ac, err = h.EnsureContext()
		if err != nil {
			return err
		}
		cmd.SetContext(context.WithValue(h.Ctx, cloud.ContextOAuth2, ac.TokenSource()))
		h.Ctx = cmd.Context()
		return nil
	}
	cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		return h.WriteConfig(ac)
	}
}

func (h *CommandHelper) Authenticate() (*AuthContext, error) {
	if h.IsQuiet {
		return nil, errors.New("can not sign in or sign up when flag --quiet is set")
	}

	ac, err := h.readConfig()
	if err != nil {
		if !errors.Is(err, ErrNoConfig) {
			return nil, err
		}
	}

	if ac.AccessToken != nil {
		fmt.Fprintf(h.VerboseWriter, "You are already logged in.\n")
		return ac, nil
	}

	ac, err = h.loginOAuth2()
	if err != nil {
		return nil, err
	}

	if err := h.WriteConfig(ac); err != nil {
		return nil, err
	}

	return ac, nil
}

var oac = oauth2.Config{
	ClientID: "ory-cli",
	Endpoint: oauth2.Endpoint{
		AuthURL:   makeCloudConsoleURL("project") + "/oauth2/auth",
		TokenURL:  makeCloudConsoleURL("project") + "/oauth2/token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
}

type data struct {
	OK          bool
	Error, Desc string
}

func (h *CommandHelper) loginOAuth2() (*AuthContext, error) {
	state := randx.MustString(32, randx.AlphaNum)
	callbackURL, code, errs, outcome, stop := h.runOAuth2CallbackServer(state)
	defer stop()

	oac.RedirectURL = callbackURL
	pkceVerifier := oauth2.GenerateVerifier()
	url := oac.AuthCodeURL(state,
		oauth2.S256ChallengeOption(pkceVerifier),
		oauth2.SetAuthURLParam("scope", "offline_access"),
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("prompt", "login consent"),
		oauth2.SetAuthURLParam("audience", makeCloudConsoleURL("api")),
	)

	_ = webbrowser.Open(url)
	fmt.Fprintf(h.VerboseErrWriter,
		`A browser should have opened for you to complete your login to Ory Network.
If no browser opened, visit the below page to continue:

		%s 

`, url)

	var authCode string
	select {
	case authCode = <-code:
		// ok
	case err := <-errs:
		fmt.Fprintf(h.VerboseErrWriter, "An error occured logging into Ory Network: %v\n", err)
		return nil, fmt.Errorf("failed OAuth2 authorization: %w", err)
	}

	token, err := oac.Exchange(
		h.Ctx,
		authCode,
		oauth2.VerifierOption(pkceVerifier),
	)
	if err != nil {
		outcome <- data{OK: false, Error: "token exchange", Desc: "An error occured during the OAuth2 token exchange: " + err.Error()}
		fmt.Fprintf(h.VerboseErrWriter, "An error occured logging into Ory Network: %v\n", err)
		return nil, fmt.Errorf("failed OAuth2 token exchange: %w", err)
	}
	outcome <- data{OK: true}

	scope, _ := token.Extra("scope").(string)
	if !slices.Contains(strings.Split(scope, " "), "offline_access") {
		fmt.Fprintf(h.VerboseErrWriter,
			"You have not granted the 'offline_access' permission during login and will have to authenticate again in %v.\n",
			time.Until(token.Expiry).Round(time.Second),
		)
	}

	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, token))
	c := newCloudClient()
	activeProject, _, err := c.ProjectAPI.GetActiveProjectInConsole(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get active project: %w", err)
	}

	fmt.Fprintf(h.VerboseErrWriter, "Successfully logged into Ory Network.\n")

	return &AuthContext{
		AccessToken:     token,
		SelectedProject: uuid.FromStringOrNil(activeProject.GetProjectId()),
	}, nil
}

func (h *CommandHelper) runOAuth2CallbackServer(state string) (callbackURL string, code <-chan string, errs <-chan error, outcome chan<- data, cleanup func()) {
	var (
		l     net.Listener
		err   error
		ports = []int{12345, 34525, 49763, 51238, 59724, 60582, 62125}
	)
	rand.Shuffle(len(ports), func(i, j int) { ports[i], ports[j] = ports[j], ports[i] })
	for _, port := range ports {
		l, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			callbackURL = fmt.Sprintf("http://localhost:%d/callback", port)
			break
		}
	}
	if l == nil {
		fmt.Fprintln(h.VerboseErrWriter, "Failed to allocate port for OAuth2 callback handler")
		os.Exit(1)
	}
	_code, _errs, _outcome := make(chan string), make(chan error), make(chan data)
	srv := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer close(_code)
			if err := r.ParseForm(); err != nil {
				redirectErr(w, r, "parse form", "An error occured during CLI authentication. Please try again")
				_errs <- err
				return
			}
			if s := r.Form.Get("state"); s != state {
				redirectErr(w, r, "state mismatch", "An error occured during CLI authentication. Please try again")
				_errs <- fmt.Errorf("state mismatch: expected %s, got %s", state, s)
				return
			}
			code := r.Form.Get("code")
			if code == "" {
				error, desc := r.Form.Get("error"), r.Form.Get("error_description")
				redirectErr(w, r, error, desc)
				_errs <- fmt.Errorf("%s: %s", error, desc)
				return
			}
			_code <- code
			if outcome := <-_outcome; !outcome.OK {
				redirectErr(w, r, outcome.Error, outcome.Desc)
				return
			}
			redirectOK(w, r)
		}),
	}
	go srv.Serve(l)
	return callbackURL, _code, _errs, _outcome, func() {
		_ = srv.Close()
	}
}

func redirectOK(w http.ResponseWriter, r *http.Request) {
	location := CloudConsoleURL("")
	location.Path = "/projects/current/dashboard"
	location.RawQuery = url.Values{"cli_auth": []string{"success"}}.Encode()
	http.Redirect(w, r, location.String(), http.StatusFound)
}

func redirectErr(w http.ResponseWriter, r *http.Request, err, desc string) {
	location := CloudConsoleURL("")
	location.Path = "/error"
	location.RawQuery = url.Values{"error": []string{err}, "error_description": []string{desc}}.Encode()
	http.Redirect(w, r, location.String(), http.StatusFound)
}

func (h *CommandHelper) SignOut() error {
	ac, err := h.readConfig()
	if err != nil {
		return err
	}
	if ac.AccessToken == nil {
		return h.WriteConfig(new(AuthContext))
	}
	revoke, err := url.Parse(oac.Endpoint.AuthURL)
	if err != nil {
		return err
	}
	revoke.Path = "/oauth2/revoke"
	res, err := http.PostForm(revoke.String(), url.Values{
		"client_id": []string{oac.ClientID},
		"token":     []string{ac.AccessToken.RefreshToken}, // this also revokes the associated access token
	})
	if err != nil {
		fmt.Fprintf(h.VerboseErrWriter, "failed to revoke access token: %v\n", err)
	} else {
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			body, _ := io.ReadAll(res.Body)
			fmt.Fprintf(h.VerboseErrWriter, "failed to revoke access token: %v\n", string(body))
		}
	}
	return h.WriteConfig(new(AuthContext))
}
