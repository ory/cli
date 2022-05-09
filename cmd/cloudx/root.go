package cloudx

import (
	"fmt"
	"strings"

	"github.com/ory/cli/cmd/cloudx/proxy"

	"github.com/spf13/cobra"
)

func NewRootCommand(parent *cobra.Command, project string, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: fmt.Sprintf("Run and manage Ory %s in Ory Cloud", project),
	}

	cmdName := strings.ToLower(project + " cloud")

	cmd.AddCommand(NewAuthCmd())
	cmd.AddCommand(NewLogoutCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewListCmd(parent))
	cmd.AddCommand(NewDeleteCmd(parent))
	cmd.AddCommand(NewPatchCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewImportCmd(parent))
	cmd.AddCommand(NewGetCmd(parent))
	cmd.AddCommand(proxy.NewProxyCommand(cmdName, version))
	cmd.AddCommand(proxy.NewTunnelCommand(cmdName, version))
	return cmd
}
