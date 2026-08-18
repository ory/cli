// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/fetcher"
)

// errNoOPLConfigured is returned when the project has no Ory Permission
// Language file, so there is nothing to print.
var errNoOPLConfigured = errors.New("no Ory Permission Language file is configured for this project")

// errLegacyNamespaces is returned for projects still using the legacy list of
// namespace definitions instead of an Ory Permission Language file.
var errLegacyNamespaces = errors.New("this project uses legacy namespace definitions instead of an Ory Permission Language file, use `ory get permission-config` to read them")

// newOPLFetcher returns the loader that resolves an Ory Permission Language
// file location. Local files are deliberately not allowed: the location comes
// from the API and must never make the CLI read from the developer's disk.
//
// The command and its tests share this constructor so that the tests pin the
// loader the command actually uses.
func newOPLFetcher() *fetcher.Fetcher {
	return fetcher.NewFetcher(fetcher.WithAllowedSchemes("http", "https", "base64"))
}

// oplLocation extracts the Ory Permission Language file location from an Ory
// Permissions config.
//
// Ory Network does not inline the OPL source in the project config; it stores a
// pointer to it, which is why `ory get permission-config` shows
// `{"namespaces":{"location":"..."}}` rather than the file itself.
func oplLocation(config map[string]interface{}) (string, error) {
	// The config is a decoded JSON document, so serializing it back cannot fail.
	raw, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("unable to read the Ory Permissions configuration: %w", err)
	}
	namespaces := gjson.GetBytes(raw, "namespaces")

	// Legacy projects list their namespace definitions here instead of pointing
	// at an OPL file. An empty list is not a legacy configuration though, there
	// simply is nothing configured.
	if namespaces.IsArray() && len(namespaces.Array()) > 0 {
		return "", errLegacyNamespaces
	}

	// Anything that is neither an OPL pointer nor legacy definitions carries no
	// file to print, and reads as an empty location here.
	location := namespaces.Get("location").String()
	if location == "" {
		return "", errNoOPLConfigured
	}
	return location, nil
}

func NewGetOPLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "opl",
		Aliases: []string{
			"namespaces-config",
		},
		Args:  cobra.NoArgs,
		Short: "Get the Ory Permission Language file from Ory Network",
		Long: `Get the Ory Permission Language file of an Ory Network project.

The file is written to stdout as-is, so it can be redirected to disk:

	ory get opl > namespace_config.ts

This is the counterpart of ` + "`ory update opl`" + `. Use it to check the
configured Ory Permission Language file into version control.

Note that ` + "`ory get permission-config`" + ` returns the location of this file
rather than its contents, because that is how Ory Network stores it.`,
		Example: `$ {{ .CommandPath }} --project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89

class Example implements Namespace {}
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			pID, err := h.ProjectID()
			if err != nil {
				return err
			}

			project, err := h.GetProject(cmd.Context(), pID, nil)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			location, err := oplLocation(project.Services.GetPermission().Config)
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
				return cmdx.FailSilently(cmd)
			}

			// The location is either a remote URL or an inlined base64 payload.
			// The fetcher verifies the HTTP status code, so an expired storage
			// link cannot be mistaken for the file, honors cancellation of the
			// command context, and redacts base64 payloads from its errors.
			opl, err := newOPLFetcher().FetchBytes(cmd.Context(), location)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "unable to read the Ory Permission Language file: %s\n", err)
				return cmdx.FailSilently(cmd)
			}

			// Writing can fail on a closed pipe (`ory get opl | head`), which is
			// not worth a usage dump.
			if _, err := cmd.OutOrStdout().Write(opl); err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "unable to write the Ory Permission Language file: %s\n", err)
				return cmdx.FailSilently(cmd)
			}
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	return cmd
}
