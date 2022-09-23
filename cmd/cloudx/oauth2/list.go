package oauth2

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd"

	"github.com/ory/cli/cmd/cloudx/client"
)

func NewListOAuth2Cmd(parent *cobra.Command) *cobra.Command {
	c := cmd.NewListClientsCmd(parent)
	client.RegisterProjectFlag(c.Flags())
	return c
}
