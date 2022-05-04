package cloudx

import (
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

func NewImportCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import resources",
	}

	cmd.AddCommand(NewImportIdentityCmd(parent))

	RegisterConfigFlag(cmd.PersistentFlags())
	RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
