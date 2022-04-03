package cloudx

import (
	"github.com/ory/kratos/cmd/identities"
	"github.com/spf13/cobra"
)

func NewDeleteIdentityCmd(parent *cobra.Command) *cobra.Command {
	cmd := identities.NewDeleteIdentityCmd(parent)
	RegisterProjectFlag(cmd.Flags())
	return cmd
}
