package cloudx

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/cli/cmd/cloudx/project"

	"github.com/ory/x/cmdx"
)

func NewListCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   fmt.Sprintf("List resources"),
	}

	cmd.AddCommand(project.NewListProjectsCmd())
	cmd.AddCommand(identity.NewListIdentityCmd(parent))

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
