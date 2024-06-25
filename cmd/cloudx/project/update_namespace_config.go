// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/osx"
)

func NewUpdateNamespaceConfigCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use: "opl",
		Aliases: []string{
			"namespaces-config",
		},
		Args:  cobra.NoArgs,
		Short: "Update the Ory Permission Language file in Ory Network",
		Example: `$ {{ .CommandPath }} --file /path/to/namespace_config.ts

class Example implements Namespace {}
`,
		Long: "Update the Ory Permission Language file in Ory Network. Legacy namespace definitions will be overwritten.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			data, err := osx.ReadFileFromAllSources(file)
			if err != nil {
				return err
			}
			patch := fmt.Sprintf(`/services/permission/config/namespaces={"location": "base64://%s"}`,
				base64.StdEncoding.EncodeToString(data))

			projectID, err := h.ProjectID()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			p, err := h.PatchProject(ctx, projectID, nil, nil, []string{patch}, nil)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintJSONAble(cmd, outputConfig(p.Project.Services.Permission.Config))

			return h.PrintUpdateProjectWarnings(p)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "",
		"Configuration file (file://namespace_config.ts, https://example.org/namespace_config.ts, ...) to update the Ory Permission Language config")
	client.RegisterYesFlag(cmd.Flags())
	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())

	return cmd
}
