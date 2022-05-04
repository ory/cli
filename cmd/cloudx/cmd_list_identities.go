package cloudx

import (
	"github.com/ory/kratos/cmd/identities"
	"github.com/spf13/cobra"
)

func NewListIdentityCmd(parent *cobra.Command) *cobra.Command {
	cmd := identities.NewListIdentitiesCmd(parent)
	RegisterProjectFlag(cmd.Flags())
	return cmd
}
