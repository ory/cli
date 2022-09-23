package oauth2

import (
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra/cmd"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewGetOAuth2Client(parent *cobra.Command) *cobra.Command {
	cmd := hydra.NewGetClientsCmd(parent)
	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
