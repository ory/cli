// Copyright © 2022 Ory Corp

package identity

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/kratos/cmd/identities"
	"github.com/ory/x/cmdx"
)

func NewGetIdentityCmd() *cobra.Command {
	cmd := identities.NewGetIdentityCmd()
	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
