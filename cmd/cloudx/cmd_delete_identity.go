package cloudx

import (
	"github.com/ory/kratos/cmd/identities"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

func NewDeleteIdentityCmd(parent *cobra.Command) *cobra.Command {
	cmd := identities.NewDeleteIdentityCmd(parent)
	RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
