// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	stderrs "errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid/v3"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
	"golang.org/x/term"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/jsonx"
	"github.com/ory/x/randx"
	"github.com/ory/x/stringsx"
)

const (
	fileName   = ".ory-cloud.json"
	ConfigFlag = "config"
	osEnvVar   = "ORY_CLOUD_CONFIG_PATH"
	Version    = "v0alpha0"
	yesFlag    = "yes"
)

func RegisterConfigFlag(f *pflag.FlagSet) {
	f.StringP(ConfigFlag, ConfigFlag[:1], "", "Path to the Ory Network configuration file.")
}

func RegisterYesFlag(f *pflag.FlagSet) {
	f.BoolP(yesFlag, yesFlag[:1], false, "Confirm all dialogs with yes.")
}

type AuthContext struct {
	Version         string        `json:"version"`
	SessionToken    string        `json:"session_token"`
	SelectedProject uuid.UUID     `json:"selected_project"`
	IdentityTraits  AuthIdentity  `json:"session_identity_traits"`
	AccessToken     *oauth2.Token `json:"oauth_token"`
}

func (i *AuthContext) ID() string {
	return i.IdentityTraits.ID
}

func (*AuthContext) Header() []string {
	return []string{"ID", "SELECTED_PROJECT"}
}

func (i *AuthContext) Columns() []string {
	return []string{
		i.ID(),
		i.SelectedProject.String(),
	}
}

func (i *AuthContext) Interface() interface{} {
	return i
}

type AuthIdentity struct {
	ID    string
	Email string `json:"email"`
}

type AuthProject struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
}

var ErrNoConfig = stderrs.New("no ory configuration file present")
var ErrNoConfigQuiet = stderrs.New("please run `ory auth` to initialize your configuration or remove the `--quiet` flag")

func getConfigPath(cmd *cobra.Command) (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "unable to guess your home directory")
	}

	return stringsx.Coalesce(
		os.Getenv(osEnvVar),
		flagx.MustGetString(cmd, ConfigFlag),
		filepath.Join(path, fileName),
	), nil
}

type CommandHelper struct {
	Ctx              context.Context
	VerboseWriter    io.Writer
	VerboseErrWriter io.Writer
	ConfigLocation   string
	NoConfirm        bool
	IsQuiet          bool
	Stdin            *bufio.Reader
	PwReader         passwordReader
}

type PasswordReader struct{}

// NewCommandHelper creates a new CommandHelper instance which handles cobra CLI commands.
func NewCommandHelper(cmd *cobra.Command) (*CommandHelper, error) {
	location, err := getConfigPath(cmd)
	if err != nil {
		return nil, err
	}

	var out = cmd.OutOrStdout()
	if flagx.MustGetBool(cmd, cmdx.FlagQuiet) {
		out = io.Discard
	}

	var outErr = cmd.ErrOrStderr()
	if flagx.MustGetBool(cmd, cmdx.FlagQuiet) {
		outErr = io.Discard
	}

	pwReader := func() ([]byte, error) {
		return term.ReadPassword(int(os.Stdin.Fd()))
	}
	if p, ok := cmd.Context().Value(PasswordReader{}).(passwordReader); ok {
		pwReader = p
	}

	return &CommandHelper{
		ConfigLocation:   location,
		NoConfirm:        flagx.MustGetBool(cmd, yesFlag),
		IsQuiet:          flagx.MustGetBool(cmd, cmdx.FlagQuiet),
		VerboseWriter:    out,
		VerboseErrWriter: outErr,
		Stdin:            bufio.NewReader(cmd.InOrStdin()),
		Ctx:              cmd.Context(),
		PwReader:         pwReader,
	}, nil
}

func (h *CommandHelper) GetDefaultProjectID() string {
	conf, err := h.readConfig()
	if err != nil {
		return ""
	}

	if conf.SelectedProject != uuid.Nil {
		return conf.SelectedProject.String()
	}

	return ""
}

func (h *CommandHelper) SetDefaultProject(id string) error {
	conf, err := h.readConfig()
	if err != nil {
		return err
	}

	uid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	conf.SelectedProject = uid
	return h.WriteConfig(conf)
}

func (h *CommandHelper) WriteConfig(c *AuthContext) error {
	c.Version = Version
	file, err := os.OpenFile(h.ConfigLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "unable to open file for writing at location: %s", file.Name())
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(c); err != nil {
		return errors.Wrapf(err, "unable to write configuration to file: %s", h.ConfigLocation)
	}

	return nil
}

func (h *CommandHelper) readConfig() (*AuthContext, error) {
	contents, err := os.ReadFile(h.ConfigLocation)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return new(AuthContext), ErrNoConfig
		}
		return nil, errors.Wrapf(err, "unable to open ory config file location: %s", h.ConfigLocation)
	}

	var c AuthContext
	if err := json.Unmarshal(contents, &c); err != nil {
		return nil, errors.Wrapf(err, "unable to JSON decode the ory config file: %s", h.ConfigLocation)
	}

	return &c, nil
}

func (h *CommandHelper) HasValidContext() (*AuthContext, bool, error) {
	c, err := h.readConfig()
	if err != nil {
		if errors.Is(err, ErrNoConfig) {
			if h.IsQuiet {
				return nil, false, errors.WithStack(ErrNoConfigQuiet)
			}
			// No context
			return nil, false, nil
		}

		return nil, false, err
	}

	if c.AccessToken != nil {
		return c, true, nil
	}

	if len(c.SessionToken) > 0 {
		client, err := NewKratosClient()
		if err != nil {
			return nil, false, err
		}

		sess, _, err := client.FrontendApi.ToSession(h.Ctx).XSessionToken(c.SessionToken).Execute()
		if err != nil {
			return nil, false, nil
		} else if sess == nil {
			return nil, false, nil
		}
		return c, true, nil
	}

	return nil, false, nil
}

func (h *CommandHelper) EnsureContext() (*AuthContext, error) {
	c, valid, err := h.HasValidContext()
	if err != nil {
		return nil, err
	} else if valid {
		return c, nil
	}

	// No valid session, but also quiet mode -> failure!
	if h.IsQuiet {
		return nil, errors.WithStack(ErrNoConfigQuiet)
	}

	// Not valid, but we have a session -> tell the user we need to re-authenticate
	_, _ = fmt.Fprintf(h.VerboseErrWriter, "Your session has expired or has otherwise become invalid. Please re-authenticate to continue.\n")

	if err := h.SignOut(); err != nil {
		return nil, err
	}

	c, err = h.Authenticate()
	if err != nil {
		return nil, err
	}

	return c, nil
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

var (
	oac = oauth2.Config{
		// ClientID: "ory-cli",
		ClientID: "7b29dd0e-3e98-4bf9-a14f-c6efbb35d508",
		Endpoint: oauth2.Endpoint{
			// AuthURL:  makeCloudConsoleURL("project") + "/oauth2/auth",
			// TokenURL: makeCloudConsoleURL("project") + "/oauth2/token",
			AuthURL:   "https://epic-swanson-8q30djkp63.projects.oryapis.com/oauth2/auth",
			TokenURL:  "https://epic-swanson-8q30djkp63.projects.oryapis.com/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	tmpl = template.Must(template.New("").Parse(`<!DOCTYPE html> 
<html>
	<head><title>Ory Network CLI login</title></head>
	<body>
		<h1>{{ .Header }}</h1>
		<p>{{ .Message }}</p>
	</body>
</html>`))
)

type data struct {
	Header, Message string
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
		oauth2.SetAuthURLParam("prompt", "consent"),
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
		outcome <- data{"Login failed", err.Error()}
		fmt.Fprintf(h.VerboseErrWriter, "An error occured logging into Ory Network: %v\n", err)
		return nil, fmt.Errorf("failed OAuth2 token exchange: %w", err)
	}
	outcome <- data{"Successfully logged into Ory Network.", "You may now close this browser tab and continue on with the Ory CLI."}

	fmt.Fprintf(h.VerboseErrWriter, "Successfully logged into Ory Network.\n")

	enc := json.NewEncoder(h.VerboseWriter)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(false)
	enc.Encode(token)

	scope, _ := token.Extra("scope").(string)
	if !slices.Contains(strings.Split(scope, " "), "offline_access") {
		fmt.Fprintf(h.VerboseErrWriter,
			"You have not granted the 'offline_access' permission during login and will have to authenticate again in %v.\n",
			time.Until(token.Expiry).Round(time.Second),
		)
	}
	for range code {
		// drain/wait
	}
	// ok, response written to browser

	return &AuthContext{AccessToken: token}, nil
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
			r.ParseForm()
			if s := r.Form.Get("state"); s != state {
				tmpl.Execute(w, &data{"Login failed", ""})
				_errs <- fmt.Errorf("state mismatch: expected %s, got %s", state, s)
				return
			}
			code := r.Form.Get("code")
			error, desc := r.Form.Get("error"), r.Form.Get("error_description")
			if code == "" {
				tmpl.Execute(w, &data{"Login failed", desc + ": " + error})
				_errs <- fmt.Errorf("%s: %s", error, desc)
				return
			}
			_code <- code
			tmpl.Execute(w, <-_outcome)
		}),
	}
	go srv.Serve(l)
	return callbackURL, _code, _errs, _outcome, func() {
		ctx, cancel := context.WithTimeout(h.Ctx, 3*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}
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

func (h *CommandHelper) ListProjects() ([]cloud.ProjectMetadata, error) {
	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, Token))
	c := newCloudClient(nil)
	projects, res, err := c.ProjectApi.ListProjects(ctx).Execute()
	if err != nil {
		return nil, handleError("unable to list projects", res, err)
	}

	return projects, nil
}

func (h *CommandHelper) GetProject(projectOrSlug string) (*cloud.Project, error) {
	if projectOrSlug == "" {
		return nil, errors.Errorf("No project selected! Please see the help message on how to set one.")
	}

	id := uuid.FromStringOrNil(projectOrSlug)
	if id == uuid.Nil {
		pjs, err := h.ListProjects()
		if err != nil {
			return nil, err
		}

		availableSlugs := make([]string, len(pjs))
		for i, pm := range pjs {
			availableSlugs[i] = pm.GetSlug()
			if strings.HasPrefix(pm.GetSlug(), projectOrSlug) {
				if id != uuid.Nil {
					return nil, errors.Errorf("The slug prefix %q is not unique, please use more characters. Found slugs:\n%s", projectOrSlug, strings.Join(availableSlugs, "\n"))
				}
				id = uuid.FromStringOrNil(pm.GetId())
			}
		}
		if id == uuid.Nil {
			return nil, errors.Errorf("no project found with slug %s, only slugs known are: %v", projectOrSlug, availableSlugs)
		}
	}

	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, Token))
	c := newCloudClient(nil)
	project, res, err := c.ProjectApi.GetProject(ctx, id.String()).Execute()
	if err != nil {
		return nil, handleError("unable to get project", res, err)
	}

	return project, nil
}

func (h *CommandHelper) CreateProject(name string, setDefault bool) (*cloud.Project, error) {
	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, Token))
	c := newCloudClient(nil)
	project, res, err := c.ProjectApi.CreateProject(ctx).CreateProjectBody(*cloud.NewCreateProjectBody(strings.TrimSpace(name))).Execute()
	if err != nil {
		return nil, handleError("unable to list projects", res, err)
	}

	if def := h.GetDefaultProjectID(); setDefault || def == "" {
		_ = h.SetDefaultProject(project.Id)
	}

	return project, nil
}

func handleError(message string, res *http.Response, err error) error {
	if e, ok := err.(*cloud.GenericOpenAPIError); ok {
		return errors.Wrapf(err, "%s: %s", message, e.Body())
	}

	if res == nil {
		return errors.Wrapf(err, "%s", message)
	}

	body, _ := io.ReadAll(res.Body)
	return errors.Wrapf(err, "%s: %s", message, body)
}

func toPatch(op string, values []string) (patches []cloud.JsonPatch, err error) {
	for _, v := range values {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			return nil, errors.Errorf("patches must be in format of `/some/config/key=some-value` but got: %s", v)
		} else if !gjson.Valid(parts[1]) {
			return nil, errors.Errorf("value for %s must be valid JSON but got: %s", parts[0], parts[1])
		}

		config, err := jsonx.EmbedSources(json.RawMessage(parts[1]), jsonx.WithIgnoreKeys("$id", "$schema"), jsonx.WithOnlySchemes("file"))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var value interface{}
		if err := json.Unmarshal(config, &value); err != nil {
			return nil, errors.WithStack(err)
		}

		patches = append(patches, cloud.JsonPatch{Op: op, Path: parts[0], Value: value})
	}
	return patches, nil
}

func (h *CommandHelper) PatchProject(id string, raw []json.RawMessage, add, replace, del []string) (*cloud.SuccessfulProjectUpdate, error) {
	var patches []cloud.JsonPatch
	for _, r := range raw {
		config, err := jsonx.EmbedSources(r, jsonx.WithIgnoreKeys("$id", "$schema"), jsonx.WithOnlySchemes("file"))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var p []cloud.JsonPatch
		if err := json.NewDecoder(bytes.NewReader(config)).Decode(&p); err != nil {
			return nil, errors.WithStack(err)
		}
		patches = append(patches, p...)
	}

	if v, err := toPatch("add", add); err != nil {
		return nil, err
	} else {
		//revive:disable indent-error-flow
		patches = append(patches, v...)
	}

	if v, err := toPatch("replace", replace); err != nil {
		return nil, err
	} else {
		//revive:disable indent-error-flow
		patches = append(patches, v...)
	}

	for _, del := range del {
		patches = append(patches, cloud.JsonPatch{Op: "remove", Path: del})
	}

	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, Token))
	c := newCloudClient(nil)
	res, _, err := c.ProjectApi.PatchProject(ctx, id).JsonPatch(patches).Execute()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *CommandHelper) UpdateProject(id string, name string, configs []json.RawMessage) (*cloud.SuccessfulProjectUpdate, error) {
	for k := range configs {
		config, err := jsonx.EmbedSources(
			configs[k],
			jsonx.WithIgnoreKeys(
				"$id",
				"$schema",
			),
			jsonx.WithOnlySchemes(
				"file",
			),
		)
		if err != nil {
			return nil, err
		}
		configs[k] = config
	}

	interim := make(map[string]interface{})
	for _, config := range configs {
		var decoded map[string]interface{}
		if err := json.Unmarshal(config, &decoded); err != nil {
			return nil, errors.WithStack(err)
		}

		if err := mergo.Merge(&interim, decoded, mergo.WithAppendSlice, mergo.WithOverride); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var payload cloud.SetProject
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(interim); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := json.NewDecoder(&b).Decode(&payload); err != nil {
		return nil, errors.WithStack(err)
	}

	if payload.Services.Identity == nil && payload.Services.Permission == nil && payload.Services.Oauth2 == nil {
		return nil, errors.Errorf("at least one of the keys `services.identity.config` and `services.permission.config` and `services.oauth2.config` is required and can not be empty")
	}

	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, Token))
	c := newCloudClient(nil)
	if name != "" {
		payload.Name = name
	} else if payload.Name == "" {
		res, _, err := c.ProjectApi.GetProject(ctx, id).Execute()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		payload.Name = res.Name
	}

	res, _, err := c.ProjectApi.SetProject(h.Ctx, id).SetProject(payload).Execute()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *CommandHelper) CreateAPIKey(projectIdOrSlug, name string) (*cloud.ProjectApiKey, error) {
	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, Token))
	c := newCloudClient(nil)
	token, _, err := c.ProjectApi.CreateProjectApiKey(ctx, projectIdOrSlug).CreateProjectApiKeyRequest(cloud.CreateProjectApiKeyRequest{Name: name}).Execute()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (h *CommandHelper) DeleteAPIKey(projectIdOrSlug, id string) error {
	ctx := context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, Token))
	c := newCloudClient(nil)
	if _, err := c.ProjectApi.DeleteProjectApiKey(ctx, projectIdOrSlug, id).Execute(); err != nil {
		return err
	}

	return nil
}
