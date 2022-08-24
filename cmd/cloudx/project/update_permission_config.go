package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewUpdatePermissionConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "permission-config <project-id>",
		Aliases: []string{
			"pc",
			"keto-config",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Update Ory Cloud Project's Permission Service Configuration",
		Example: `$ ory update permission-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--file /path/to/config.json \
	--file /path/to/config.yml \
	--file https://example.org/config.yaml \
	--file base64://<json> \
	--format json

{
  "namespaces": [
    {
      "name": "files",
      "id": 2
	},
    // ...
  ]
}

$ ory update permission-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--file /path/to/keto-config.yaml \
    --format yaml

namespaces:
  - name: files
    # ...`,
		Long: `Use this command to replace your current Ory Cloud Project's permission service configuration. All values
of the permission service will be overwritten. To update individual settings use the ` + "`patch`" + ` command instead.

Compared to the ` + "`update project`" + ` command, this command only updates the permission service configuration
and also only returns the permission service configuration as a result. This command is useful when you want to
import an Ory Keto config as well, for example.

The full configuration payload can be found at:

	https://www.ory.sh/docs/reference/api#operation/updateProject.

This command expects the contents of the ` + "`/services/permission/config`" + ` key, so for example:

    {
	  "namespaces": [
        {
          "name": "files",
          "id": 2
        },
        // ...
      ]
    }
`,
		RunE: runUpdate(prefixFilePermissionConfig, outputPermissionConfig),
	}

	cmd.Flags().StringSliceP("file", "f", nil, "Configuration file(s) (file://config.json, https://example.org/config.yaml, ...) to update the permission config")
	client.RegisterYesFlag(cmd.Flags())
	cmdx.RegisterNoiseFlags(cmd.Flags())
	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	return cmd
}
