package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewPerformCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perform",
		Short: "Perform a flow",
	}

	cmd.AddCommand(
		oauth2.NewPerformAuthorizationCodeCmd(parent),
	)

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())

	return cmd
}
