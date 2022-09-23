package oauth2

import (
	"github.com/spf13/cobra"

	cmd2 "github.com/ory/hydra/cmd"

	"github.com/ory/cli/cmd/cloudx/client"
)

func NewImportOAuth2Cmd(parent *cobra.Command) *cobra.Command {
	cmd := cmd2.NewImportClientCmd(parent)
	client.RegisterProjectFlag(cmd.Flags())
	return cmd
}
