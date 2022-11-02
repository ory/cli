// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewPatchKetoConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "permission-config <project-id>",
		Aliases: []string{"pc", "keto-config"},
		Args:    cobra.ExactArgs(1),
		Short:   "Patch an Ory Permissions config",
		Example: `$ ory patch permission-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--add '/namespaces=[{"name":"files", "id": 2}]' \
	--replace '/namespaces/2/name="directories"' \
	--remove '/limit/max_read_depth' \
	--format json-pretty

{
  "namespaces": [
    {
      "name": "files",
      "id": 2
    },
    {
      "name": "directories",
      "id": 3
    },
    // ...
  ]
}`,
		Long: `Patches an Ory Permissions configuration. Only values
specified in the patch will be overwritten. To replace the config use the ` + "`update`" + ` command instead.

Compared to the ` + "`patch project`" + ` command, this command only updates the permission service configuration
and also only returns the permission service configuration as a result. This command is useful when you want to
import an Ory Keto config as well, for example. This allows for shorter paths when specifying the flags

	ory patch permission-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
		--replace '/limit/max_read_depth=5'

when compared to the ` + "`patch project`" + ` command:

	ory patch project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
		--replace '/services/permission/config/limit/max_read_depth=5'

The format of the patch is a JSON-Patch document. For more details please check:

	https://www.ory.sh/docs/reference/api#operation/patchProject
	https://jsonpatch.com`,
		RunE: runPatch(
			prefixPermissionConfig,
			prefixFilePermissionConfig,
			outputPermissionConfig,
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
