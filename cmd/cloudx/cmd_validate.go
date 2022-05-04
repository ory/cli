package cloudx

import (
	"github.com/ory/kratos/cmd/identities"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate resources",
	}

	cmd.AddCommand(identities.NewValidateIdentityCmd())

	RegisterConfigFlag(cmd.PersistentFlags())
	RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
