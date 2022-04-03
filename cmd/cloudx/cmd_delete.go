package cloudx

import (
	"fmt"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: fmt.Sprintf("Delete resources"),
	}

	cmd.AddCommand(NewDeleteIdentityCmd(parent))

	RegisterConfigFlag(cmd.PersistentFlags())
	RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
