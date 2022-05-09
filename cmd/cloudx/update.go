package cloudx

import (
	"fmt"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/project"

	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: fmt.Sprintf("Update resources"),
	}
	cmd.AddCommand(project.NewProjectsUpdateCmd())
	cmd.AddCommand(project.NewUpdateIdentityConfigCmd())
	client.RegisterConfigFlag(cmd.PersistentFlags())
	return cmd
}
