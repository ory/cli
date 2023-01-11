// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewGetOAuth2ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oauth2-config [project-id]",
		Aliases: []string{"oc", "oauth2-config"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Ory OAuth2 & OpenID Connect configuration.",
		Long:    "Get the Ory OAuth2 & OpenID Connect configuration for the specified Ory Network project.",
		Example: `$ ory get oauth2-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format yaml > oauth2-config.yaml

$ ory get oauth2-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format json

{
  "oauth2": {
    "client_credentials": {
      "default_grant_allowed_scope": false
    },
    "expose_internal_errors": true,
    "grant": {
      "jwt": {
        "iat_optional": false,
        "jti_optional": false,
        "max_ttl": "720h0m0s"
      }
    }
  },
  // ...
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			id, err := getSelectedProjectId(h, args)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			project, err := h.GetProject(id)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintJSONAble(cmd, outputConfig(project.Services.Oauth2.Config))
			return nil
		},
	}

	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	return cmd
}
