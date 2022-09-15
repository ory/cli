package oauth2

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewDeleteOAuth2Cmd(parent *cobra.Command) *cobra.Command {
	c := cmd.NewDeleteClientCmd(parent)
	client.RegisterProjectFlag(c.Flags())
	cmdx.RegisterFormatFlags(c.Flags())
	return c
}
