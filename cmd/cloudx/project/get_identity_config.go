package project

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewGetKratosConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "identity-config <project-id>",
		Aliases: []string{"ic", "kratos-config"},
		Args:    cobra.ExactArgs(1),
		Short:   "Get an Ory Identities configuration",
		Long:    "Get an Ory Identities configuration.",
		Example: `$ ory get kratos-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format yaml > kratos-config.yaml

$ ory get kratos-config ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 --format json

{
  "selfservice": {
	"methods": {
	  "password": { "enabled": false }
	}
	// ...
  }
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			project, err := h.GetProject(args[0])
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintJSONAble(cmd, outputConfig(project.Services.Identity.Config))
			return nil
		},
	}

	cmdx.RegisterJSONFormatFlags(cmd.Flags())
	return cmd
}
