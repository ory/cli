package cloudx

import (
	"github.com/ory/cli/cmd/cloudx/accountexperience"
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

func NewOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open Ory Account Experience Pages",
	}
	cmd.AddCommand(accountexperience.NewAccountExperienceOpenCmd())
	client.RegisterProjectFlag(cmd.PersistentFlags())
	client.RegisterConfigFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())

	return cmd
}
