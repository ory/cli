// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/pkg/browser"

	"github.com/ory/x/pointerx"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/jsonx"
)

const (
	WorkspaceKey = "ORY_WORKSPACE"
	ProjectKey   = "ORY_PROJECT"

	WorkspaceAPIKey = "ORY_WORKSPACE_API_KEY"
	ProjectAPIKey   = "ORY_PROJECT_API_KEY"
)

var ErrProjectNotSet = fmt.Errorf("no project was specified")

type (
	CommandHelper struct {
		config            *Config
		projectOverride   *string
		workspaceOverride *string
		configLocation    string
		noConfirm         bool
		isQuiet           bool
		VerboseErrWriter  io.Writer
		Stdin             *bufio.Reader
		openBrowserHook   func(string) error
		projectAPIKey     *string
		workspaceAPIKey   *string
	}
	helperOptionsContextKey struct{}
	CommandHelperOption     func(*CommandHelper)
)

func ContextWithOptions(ctx context.Context, opts ...CommandHelperOption) context.Context {
	baseOpts, _ := ctx.Value(helperOptionsContextKey{}).([]CommandHelperOption)
	newOpts := make([]CommandHelperOption, len(baseOpts)+len(opts))
	copy(newOpts, baseOpts)
	copy(newOpts[len(baseOpts):], opts)
	return context.WithValue(ctx, helperOptionsContextKey{}, newOpts)
}

func WithConfigLocation(location string) CommandHelperOption {
	return func(h *CommandHelper) {
		h.configLocation = location
	}
}

func WithNoConfirm(noConfirm bool) CommandHelperOption {
	return func(h *CommandHelper) {
		h.noConfirm = noConfirm
	}
}

func WithQuiet(isQuiet bool) CommandHelperOption {
	return func(h *CommandHelper) {
		h.isQuiet = isQuiet
	}
}

func WithVerboseErrWriter(w io.Writer) CommandHelperOption {
	return func(h *CommandHelper) {
		h.VerboseErrWriter = w
	}
}

func WithStdin(r io.Reader) CommandHelperOption {
	return func(h *CommandHelper) {
		h.Stdin = bufio.NewReader(r)
	}
}

func WithProjectOverride(project string) CommandHelperOption {
	return func(h *CommandHelper) {
		h.projectOverride = &project
	}
}

func WithWorkspaceOverride(workspace string) CommandHelperOption {
	return func(h *CommandHelper) {
		h.workspaceOverride = pointerx.Ptr(workspace)
	}
}

func WithProjectAPIKey(apiKey string) CommandHelperOption {
	return func(h *CommandHelper) {
		h.projectAPIKey = &apiKey
	}
}

func WithWorkspaceAPIKey(apiKey string) CommandHelperOption {
	return func(h *CommandHelper) {
		h.workspaceAPIKey = &apiKey
	}
}

func WithSessionToken(_ testing.TB, sessionToken string) CommandHelperOption {
	return func(h *CommandHelper) {
		h.workspaceAPIKey = &sessionToken
		h.projectAPIKey = &sessionToken
	}
}

func WithOpenBrowserHook(openBrowser func(string) error) CommandHelperOption {
	return func(h *CommandHelper) {
		h.openBrowserHook = openBrowser
	}
}

// NewCobraCommandHelper creates a new CommandHelper instance which handles cobra CLI commands.
func NewCobraCommandHelper(cmd *cobra.Command, opts ...CommandHelperOption) (*CommandHelper, error) {
	stdErr := cmd.ErrOrStderr()
	quiet := flagx.MustGetBool(cmd, cmdx.FlagQuiet)
	if quiet {
		stdErr = io.Discard
	}
	defaultOpts := []CommandHelperOption{
		WithVerboseErrWriter(stdErr),
		WithStdin(cmd.InOrStdin()),
		WithQuiet(quiet),
		WithNoConfirm(flagx.MustGetBool(cmd, FlagYes)),
	}
	// we explicitly ignore the error here, because the command might not support the project flag (most do)
	if project, _ := cmd.Flags().GetString(FlagProject); project != "" {
		defaultOpts = append(defaultOpts, WithProjectOverride(project))
	}
	// we explicitly ignore the error here, because the command might not support the workspace flag (most do)
	if workspace, _ := cmd.Flags().GetString(FlagWorkspace); workspace != "" {
		defaultOpts = append(defaultOpts, WithWorkspaceOverride(workspace))
	}
	if config := flagx.MustGetString(cmd, FlagConfig); config != "" {
		defaultOpts = append(defaultOpts, WithConfigLocation(config))
	}
	return NewCommandHelper(cmd.Context(), append(defaultOpts, opts...)...)
}

func NewCommandHelper(ctx context.Context, opts ...CommandHelperOption) (*CommandHelper, error) {
	location, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	h := &CommandHelper{
		configLocation:   location,
		noConfirm:        false,
		VerboseErrWriter: io.Discard,
		Stdin:            bufio.NewReader(os.Stdin),
		openBrowserHook: func(uri string) error {
			// we ignore the error in this case, as we also log the URL and we cannot recover in any way
			_ = browser.OpenURL(uri)
			return nil
		},
	}
	if ctxOpts, ok := ctx.Value(helperOptionsContextKey{}).([]CommandHelperOption); ok {
		for _, o := range ctxOpts {
			o(h)
		}
	}
	for _, o := range opts {
		o(h)
	}
	config, err := h.getOrCreateConfig()
	if err != nil {
		return nil, err
	}

	getAPIKey := func(envKey string, override *string) *string {
		if override != nil {
			return override
		}
		if key, ok := os.LookupEnv(envKey); ok {
			return &key
		}
		return nil
	}
	h.workspaceAPIKey = getAPIKey(WorkspaceAPIKey, h.workspaceAPIKey)
	h.projectAPIKey = getAPIKey(ProjectAPIKey, h.projectAPIKey)

	{
		// determine current workspace from all possible sources
		workspace := ""
		if h.workspaceOverride != nil {
			workspace = *h.workspaceOverride
		} else if ws, ok := os.LookupEnv(WorkspaceKey); ok {
			workspace = ws
		} else {
			if config.SelectedWorkspace != uuid.Nil {
				workspace = config.SelectedWorkspace.String()
			}
		}
		workspace = strings.TrimSpace(workspace)

		if id, err := uuid.FromString(workspace); err == nil {
			h.workspaceOverride = pointerx.Ptr(id.String())
		} else if workspace != "" {
			ws, err := h.findWorkspace(ctx, workspace)
			if err != nil {
				return nil, err
			}
			if ws != nil {
				h.workspaceOverride = pointerx.Ptr(ws.Id)
			}
		}
	}
	{
		// determine current project from all possible sources
		project := ""
		if h.projectOverride != nil {
			project = *h.projectOverride
		} else if pj, ok := os.LookupEnv(ProjectKey); ok {
			project = pj
		} else if config.SelectedProject != uuid.Nil {
			project = config.SelectedProject.String()
		}
		project = strings.TrimSpace(project)

		if id, err := uuid.FromString(project); err == nil {
			h.projectOverride = pointerx.Ptr(id.String())
		} else if project != "" {
			pj, err := h.findProject(ctx, project, h.workspaceOverride)
			if err != nil {
				return nil, err
			}
			if pj != nil {
				h.projectOverride = pointerx.Ptr(pj.Id)
			}
		}
	}

	return h, nil
}

func (h *CommandHelper) ProjectID() (string, error) {
	if h.projectOverride == nil {
		return "", ErrProjectNotSet
	}
	return *h.projectOverride, nil
}

func (h *CommandHelper) WorkspaceID() *string {
	return h.workspaceOverride
}

func (h *CommandHelper) UserName(ctx context.Context) string {
	config, err := h.GetAuthenticatedConfig(ctx)
	if err == nil && config.IdentityTraits.Name != "" {
		return config.IdentityTraits.Name
	}
	u, err := user.Current()
	if err != nil {
		return "unknown"
	}
	if u.Name != "" {
		return u.Name
	}
	return u.Username
}

func (h *CommandHelper) OpenURL(uri string) error {
	return h.openBrowserHook(uri)
}

func handleError(message string, res *http.Response, err error) error {
	if e := new(cloud.GenericOpenAPIError); errors.As(err, &e) {
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
		path, value, found := strings.Cut(v, "=")
		if !found {
			return nil, errors.Errorf("patches must be in format of `/some/config/key=some-value` but got: %s", v)
		} else if !gjson.Valid(value) {
			return nil, errors.Errorf("value for %s must be valid JSON but got: %s", path, value)
		}

		config, err := jsonx.EmbedSources(json.RawMessage(value), jsonx.WithIgnoreKeys("$id", "$schema"), jsonx.WithOnlySchemes("file"))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		patches = append(patches, cloud.JsonPatch{Op: op, Path: path, Value: config})
	}
	return patches, nil
}
