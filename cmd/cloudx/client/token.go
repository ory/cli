package client

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var Token *oauth2.Token

func RegisterAuthHelpers(cmd *cobra.Command) {
	var (
		h  *CommandHelper
		ac *AuthContext
	)
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		fmt.Fprintf(cmd.OutOrStderr(), "RegisterAuthHelpers.PersistentPreRunE\n")
		h, err = NewCommandHelper(cmd)
		if err != nil {
			return err
		}
		ac, err = h.EnsureContext()
		if err != nil {
			return err
		}
		Token = ac.AccessToken
		return nil
	}
	cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(cmd.OutOrStderr(), "RegisterAuthHelpers.PersistentPostRunE\n")
		return h.WriteConfig(ac)
	}
}
