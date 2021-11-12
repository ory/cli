package auth

import (
	"github.com/ory/cli/x/cloudx"
	"github.com/spf13/cobra"
)

func NewMainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "auth",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := cloudx.NewHandler(cmd)
			if err != nil {
				return err
			}
			if _, err = h.Authenticate(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(NewLogoutCmd())
	cloudx.RegisterFlags(cmd.PersistentFlags())
	return cmd
}
