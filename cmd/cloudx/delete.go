package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"
	"github.com/ory/cli/cmd/cloudx/relationtuples"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/x/cmdx"
)

func NewDeleteCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
	}

	cmd.AddCommand(identity.NewDeleteIdentityCmd(parent))
	cmd.AddCommand(oauth2.NewDeleteOAuth2Cmd(parent))
	cmd.AddCommand(relationtuples.Delete())

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
