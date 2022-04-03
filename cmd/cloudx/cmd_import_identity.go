package cloudx

import (
	"github.com/ory/kratos/cmd/identities"
	"github.com/spf13/cobra"
)

func NewImportIdentityCmd(parent *cobra.Command) *cobra.Command {
	cmd := identities.NewImportIdentitiesCmd(parent)
	RegisterProjectFlag(cmd.Flags())
	return cmd
}
