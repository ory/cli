package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/x/cmdx"
)

func NewImportCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import resources",
	}

	cmd.AddCommand(identity.NewImportIdentityCmd(parent))
	cmd.AddCommand(oauth2.NewImportOAuth2Cmd(parent))

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
