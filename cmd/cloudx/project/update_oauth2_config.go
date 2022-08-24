package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/ory/x/cmdx"
)

func NewUpdateOAuth2ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "oauth2-config <project-id>",
		Aliases: []string{
			"oc",
			"hydra-config",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Update Ory Cloud Project's OAuth2 Service Configuration",
		Example: `$ ory update oauth2-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--file /path/to/config.json \
	--file /path/to/config.yml \
	--file https://example.org/config.yaml \
	--file base64://<json> \
	--format json

{
  "oauth2": {
    "pkce": {
      "enabled": true
    }
  },
  // ...
}

$ ory update oauth2-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--file /path/to/keto-config.yaml \
    --format yaml

oauth2:
  pkce:
    enabled: true
# ...`,
		Long: `Use this command to replace your current Ory Cloud Project's OAuth2 service configuration. All values
of the OAuth2 service will be overwritten. To update individual settings use the ` + "`patch`" + ` command instead.

Compared to the ` + "`update project`" + ` command, this command only updates the OAuth2 service configuration
and also only returns the OAuth2 service configuration as a result. This command is useful when you want to
import an Ory Keto config as well, for example.

The full configuration payload can be found at:

	https://www.ory.sh/docs/reference/api#operation/updateProject.

This command expects the contents of the ` + "`/services/oauth2/config`" + ` key, so for example:

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
		RunE: runUpdate(prefixFileOAuth2Config, outputOAuth2Config),
	}

	cmd.Flags().StringSliceP("file", "f", nil, "Configuration file(s) (file://config.json, https://example.org/config.yaml, ...) to update the oAuth2 config")
	client.RegisterYesFlag(cmd.Flags())
	cmdx.RegisterNoiseFlags(cmd.Flags())
	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	return cmd
}
