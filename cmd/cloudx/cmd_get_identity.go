package cloudx

import (
	"github.com/ory/kratos/cmd/identities"
	"github.com/spf13/cobra"
)

func NewGetIdentityCmd(parent *cobra.Command) *cobra.Command {
	cmd := identities.NewGetIdentityCmd(parent)
	RegisterProjectFlag(cmd.Flags())
	return cmd
}
