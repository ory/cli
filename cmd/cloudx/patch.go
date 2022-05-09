package cloudx

import (
	"fmt"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/project"

	"github.com/spf13/cobra"
)

func NewPatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patch",
		Short: fmt.Sprintf("Patch resources"),
	}
	client.RegisterConfigFlag(cmd.PersistentFlags())
	cmd.AddCommand(project.NewProjectsPatchCmd())
	cmd.AddCommand(project.NewPatchKratosConfigCmd())
	return cmd
}
