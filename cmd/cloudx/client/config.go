// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/gofrs/uuid"
	"github.com/spf13/pflag"

	"github.com/ory/x/cmdx"

	"github.com/ory/x/stringsx"
)

var (
	ErrNoConfig         = errors.New("no ory configuration file present")
	ErrNoConfigQuiet    = errors.New("please authenticate the CLI or remove the `--quiet` flag")
	ErrNotAuthenticated = errors.New("you are not authenticated, please run `ory auth` to authenticate")
	ErrReauthenticate   = errors.New("your session or key has expired or has otherwise become invalid, re-authenticate to continue")
)

const (
	ConfigFileName = ".ory-cloud.json"
	FlagConfig     = "config"
	ConfigPathKey  = "ORY_CONFIG_PATH"
	ConfigVersion  = "v1"
)

func RegisterConfigFlag(f *pflag.FlagSet) {
	f.StringP(FlagConfig, FlagConfig[:1], "", "Path to the Ory Network configuration file.")
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to guess your home directory: %w", err)
	}

	return stringsx.Coalesce(
		os.Getenv(ConfigPathKey),
		filepath.Join(homeDir, ConfigFileName),
	), nil
}

func (c *Config) writeUpdate() error {
	c.Version = ConfigVersion

	f, err := os.OpenFile(c.location, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file %q for writing: %w", c.location, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(c); err != nil {
		return fmt.Errorf("unable to write configuration file %q: %w", c.location, err)
	}
	return nil
}

func (h *CommandHelper) UpdateConfig(c *Config) error {
	h.config = c
	return c.writeUpdate()
}

func (h *CommandHelper) getOrCreateConfig() (*Config, error) {
	c, err := h.getConfig()
	if errors.Is(err, ErrNoConfig) {
		return &Config{
			location: h.configLocation,
		}, nil
	}
	return c, err
}

func (h *CommandHelper) getConfig() (*Config, error) {
	if h.config == nil {
		c, err := readConfig(h.configLocation)
		if err != nil {
			return nil, err
		}
		switch c.Version {
		case "v0alpha0":
			if h.isQuiet {
				return nil, fmt.Errorf("you have to authenticate the Ory CLI now differently, plese see ory auth for details")
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Thanks for upgrading! You will now be prompted to log in to the Ory CLI through the Ory Console.")
			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Press enter to continue...")
			_, err := h.Stdin.ReadString('\n')
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("unable to read from stdin: %w", err)
			}
			fallthrough
		default:
			return nil, ErrNoConfig
		case ConfigVersion:
			// pass
		}
		h.config = c
	}
	return h.config, nil
}

func readConfig(location string) (*Config, error) {
	f, err := os.Open(location)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrNoConfig
		}
		return nil, fmt.Errorf("unable to open ory config file %q: %w", location, err)
	}
	defer f.Close()

	c := Config{
		location: location,
	}
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, fmt.Errorf("unable to JSON decode the ory config file %q: %w", location, err)
	}

	return &c, nil
}

func (h *CommandHelper) SelectWorkspace(id string) error {
	conf, err := h.getOrCreateConfig()
	if err != nil {
		return err
	}

	uid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	if conf.SelectedWorkspace == uid {
		// nothing to do
		return nil
	}

	conf.SelectedWorkspace = uid
	h.workspaceOverride = &id
	return h.UpdateConfig(conf)
}

func (h *CommandHelper) SelectProject(id string) error {
	conf, err := h.getOrCreateConfig()
	if err != nil {
		return err
	}

	uid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	if conf.SelectedProject == uid {
		// nothing to do
		return nil
	}

	conf.SelectedProject = uid
	h.projectOverride = &id
	return h.UpdateConfig(conf)
}

type Config struct {
	Version           string        `json:"version"`
	AccessToken       *oauth2.Token `json:"access_token"`
	SelectedProject   uuid.UUID     `json:"selected_project"`
	SelectedWorkspace uuid.UUID     `json:"selected_workspace"`
	IdentityTraits    Identity      `json:"session_identity_traits"`

	// isAuthenticated is a flag that we set once the session was checked and is valid.
	// Because this is not stored to the config file, it means that every command execution does at most one session check.
	isAuthenticated bool
	// location is the path to the configuration file
	location string
}

func (c *Config) ID() string {
	return c.IdentityTraits.ID.String()
}

func (*Config) Header() []string {
	return []string{"ID", "EMAIL", "SELECTED PROJECT", "SELECTED WORKSPACE"}
}

func (c *Config) Columns() []string {
	project, workspace := cmdx.None, cmdx.None
	if c.SelectedProject != uuid.Nil {
		project = c.SelectedProject.String()
	}
	if c.SelectedWorkspace != uuid.Nil {
		workspace = c.SelectedWorkspace.String()
	}
	return []string{
		c.ID(),
		c.IdentityTraits.Email,
		project,
		workspace,
	}
}

func (c *Config) Interface() any {
	return c
}

type Identity struct {
	ID    uuid.UUID
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (c *Config) TokenSource(ctx context.Context) oauth2.TokenSource {
	return &autoStoreRefreshedTokenSource{
		ctx: ctx,
		c:   c,
	}
}

// autoStoreRefreshedTokenSource is a token source that automatically stores the refreshed token in the config file.
// Because it holds the context, it should not be re-used. Always create a new one using Config.TokenSource() with the current context.
type autoStoreRefreshedTokenSource struct {
	ctx context.Context
	c   *Config
}

func (s *autoStoreRefreshedTokenSource) Token() (*oauth2.Token, error) {
	newToken, err := oauth2ClientConfig().TokenSource(s.ctx, s.c.AccessToken).Token()
	if err != nil {
		return nil, err
	}
	if newToken.AccessToken != s.c.AccessToken.AccessToken {
		s.c.AccessToken = newToken
		if err := s.c.writeUpdate(); err != nil {
			return nil, err
		}
	}
	return newToken, nil
}
