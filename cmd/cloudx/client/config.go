// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (h *CommandHelper) UpdateConfig(c *Config) error {
	c.Version = ConfigVersion
	h.config = c

	f, err := os.OpenFile(h.configLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file %q for writing: %w", h.configLocation, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(c); err != nil {
		return fmt.Errorf("unable to write configuration file %q: %w", h.configLocation, err)
	}

	return nil
}

func (h *CommandHelper) getOrCreateConfig() (*Config, error) {
	c, err := h.getConfig()
	if errors.Is(err, ErrNoConfig) {
		return &Config{}, nil
	}
	return c, err
}

func (h *CommandHelper) getConfig() (*Config, error) {
	if h.config == nil {
		c, err := readConfig(h.configLocation)
		if err != nil {
			return nil, err
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

	var c Config
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
	return oauth2ClientConfig().TokenSource(ctx, c.AccessToken)
}
