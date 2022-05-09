package identity

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/kratos/cmd/identities"
)

func NewImportIdentityCmd(parent *cobra.Command) *cobra.Command {
	cmd := identities.NewImportIdentitiesCmd(parent)
	client.RegisterProjectFlag(cmd.Flags())
	return cmd
}
