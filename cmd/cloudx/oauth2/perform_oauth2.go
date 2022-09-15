package oauth2

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/hydra/cmd"
	"github.com/ory/x/cmdx"
)

func NewPerformAuthorizationCodeCmd(parent *cobra.Command) *cobra.Command {
	c := cmd.NewPerformAuthorizationCodeCmd(parent)
	client.RegisterProjectFlag(c.Flags())
	cmdx.RegisterFormatFlags(c.Flags())
	return c
}
