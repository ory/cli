package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewUpdateIdentityConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "identity-config <project-id>",
		Aliases: []string{
			"ic",
			"kratos-config",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Update Ory Cloud project's identity service configuration",
		Example: `$ ory update identity-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--file /path/to/config.json \
	--file /path/to/config.yml \
	--file https://example.org/config.yaml \
	--file base64://<json> \
	--format json

{
  "courier": {
	"smtp": {
	  "from_name": "..."
	}
	// ...
  }
}

$ ory update identity-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--file /path/to/kratos-config.yaml \
    --format yaml

courier:
  smtp:
    # ...`,
		Long: `Updates your Ory Cloud project's identity service configuration. All values
of the identity service will be overwritten. To update individual settings use the ` + "`patch`" + ` command instead.

Compared to the ` + "`update project`" + ` command, this command only updates the identity service configuration
and also only returns the identity service configuration as a result. This command is useful when you want to
import an Ory Kratos config as well, for example.

The full configuration payload can be found at:

	https://www.ory.sh/docs/reference/api#operation/updateProject.

This command expects the contents of the ` + "`/services/identity/config`" + ` key, so for example:

	{
	  "courier": {
		"smtp": {
		  "from_name": "..."
		}
		// ...
	  }
	}
`,
		RunE: runUpdate(prefixFileIdentityConfig, outputIdentityConfig),
	}

	cmd.Flags().StringSliceP("file", "f", nil, "Configuration file(s) (file://config.json, https://example.org/config.yaml, ...) to update the identity config")
	client.RegisterYesFlag(cmd.Flags())
	cmdx.RegisterNoiseFlags(cmd.Flags())
	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	return cmd
}
