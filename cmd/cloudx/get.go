package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/cli/cmd/cloudx/project"
	"github.com/ory/x/cmdx"
)

func NewGetCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a resource",
	}

	cmd.AddCommand(
		project.NewGetProjectCmd(),
		project.NewGetKratosConfigCmd(),
		project.NewGetKetoConfigCmd(),
		project.NewGetOAuth2ConfigCmd(),
		identity.NewGetIdentityCmd(parent),
		oauth2.NewGetOAuth2Client(parent),
	)

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())

	return cmd
}
