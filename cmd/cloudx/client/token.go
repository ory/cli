package client

import (
	"context"

	"github.com/spf13/cobra"

	cloud "github.com/ory/client-go"
)

func RegisterAuthHelpers(cmd *cobra.Command) {
	var (
		h  *CommandHelper
		ac *AuthContext
	)
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		h, err = NewCommandHelper(cmd)
		if err != nil {
			return err
		}
		ac, err = h.EnsureContext()
		if err != nil {
			return err
		}
		cmd.SetContext(context.WithValue(h.Ctx, cloud.ContextOAuth2, oac.TokenSource(h.Ctx, ac.AccessToken)))
		h.Ctx = cmd.Context()
		return nil
	}
	cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		return h.WriteConfig(ac)
	}
}
