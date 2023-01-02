// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/ory/x/cmdx"
)

func NewPatchOAuth2ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oauth2-config <project-id>",
		Aliases: []string{"oc", "hydra-config"},
		Args:    cobra.ExactArgs(1),
		Short:   "Patch the Ory OAuth2 & OpenID Connect configuration of the specified Ory Network project.",
		Example: `$ ory patch oauth2-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--replace '/strategies/access_token="jwt"' \
	--add '/ttl/login_consent_request="1h"' \
	--remove '/strategies/scope' \
	--format json-pretty

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
}
`,
		Long: `Patch the Ory OAuth2 & OpenID Connect configuration of the specified Ory Network project. Only values
specified in the patch will be overwritten. To replace the config use the ` + "`update`" + ` command instead.

Compared to the ` + "`patch project`" + ` command, this command updates only the Ory OAuth2 & OpenID Connect 
configuration and returns configuration as a result. This command is useful when you want to import configuration 
from self-hosted Ory Hydra to Ory Network. Using this command allows for shorter paths when specifying the flags:

	ory patch oauth2-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
		--replace '/strategies/access_token="jwt"'

when compared to the ` + "`patch project`" + ` command:

	ory patch project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
		--replace '/strategies/access_token="jwt"'

The format of the patch is a JSON-Patch document. For more details please check:

	https://www.ory.sh/docs/reference/api#operation/patchProject
	https://jsonpatch.com`,
		RunE: runPatch(
			prefixOAuth2Config,
			prefixFileOAuth2Config,
			outputOAuth2Config,
		),
	}

	cmd.Flags().StringSliceP("file", "f", nil, "Configuration file(s) (file://config.json, https://example.org/config.yaml, ...) to update the project")
	cmd.Flags().StringArray("replace", nil, "Replace a specific key in the configuration")
	cmd.Flags().StringArray("add", nil, "Add a specific key to the configuration")
	cmd.Flags().StringArray("remove", nil, "Remove a specific key from the configuration")
	client.RegisterYesFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
