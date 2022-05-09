package cloudx

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/x/cmdx"
)

func NewDeleteCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: fmt.Sprintf("Delete resources"),
	}

	cmd.AddCommand(identity.NewDeleteIdentityCmd(parent))

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
