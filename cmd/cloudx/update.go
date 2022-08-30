package cloudx

import (
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/project"

	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update resources",
	}
	cmd.AddCommand(
		project.NewProjectsUpdateCmd(),
		project.NewUpdateIdentityConfigCmd(),
		project.NewUpdateOAuth2ConfigCmd(),
		project.NewUpdatePermissionConfigCmd(),
	)
	client.RegisterConfigFlag(cmd.PersistentFlags())
	return cmd
}
