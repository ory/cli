package auth

import (
	"fmt"
	"github.com/ory/cli/x/cloudx"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Sign out of your Ory Cloud account",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := cloudx.NewHandler(cmd)
			if err != nil {
				return err
			}
			if err := h.SignOut();err != nil {
				return err
			}
			fmt.Println("You signed out successfully.")
			return nil
		},
	}

	return cmd
}
