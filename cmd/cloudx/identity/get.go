package identity

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/kratos/cmd/identities"
	"github.com/ory/x/cmdx"
)

func NewGetIdentityCmd(parent *cobra.Command) *cobra.Command {
	cmd := identities.NewGetIdentityCmd(parent)
	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
