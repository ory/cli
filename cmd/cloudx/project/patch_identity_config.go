// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewPatchKratosConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "identity-config <project-id>",
		Aliases: []string{"ic", "kratos-config"},
		Args:    cobra.ExactArgs(1),
		Short:   "Patch the Ory Identities configuration of the defined Ory Network project.",
		Example: `$ ory patch identity-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--add '/courier/smtp={"from_name":"My new email name"}' \
	--replace '/selfservice/methods/password/enabled=false' \
	--remove '/selfservice/methods/totp/enabled' \
	--format json-pretty

{
  "selfservice": {
    "methods": {
      "password": { "enabled": false }
    }
    // ...
  }
}`,
		Long: `Patch the Ory Identities configuration of the defined Ory Network project. Only values specified in the patch will be overwritten. To replace the config use the ` + "`update`" + ` command instead.

Compared to the ` + "`patch project`" + ` command, this command updates only the Ory Identities configuration
and returns the configuration as a result. This command is useful when you want to import configuration from 
self-hosted Ory Kratos to Ory Network. Using this command allows for shorter paths when specifying the flags:

	ory patch identity-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
		--replace '/selfservice/methods/password/enabled=false'

when compared to the ` + "`patch project`" + ` command:

	ory patch project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
		--replace '/services/identity/config/selfservice/methods/password/enabled=false'

The format of the patch is a JSON-Patch document. For more details please check:

	https://www.ory.sh/docs/reference/api#operation/patchProject
	https://jsonpatch.com`,
		RunE: runPatch(
			prefixIdentityConfig,
			prefixFileIdentityConfig,
			outputIdentityConfig,
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
