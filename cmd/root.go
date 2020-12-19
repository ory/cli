package cmd

import (
	"context"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/kratos-client-go/client"
	"github.com/ory/kratos/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ory",
	Short: "The ORY CLI",
}

func Execute() {
	ctx := context.WithValue(context.Background(), cliclient.ClientContextKey, func(cmd *cobra.Command) *client.OryKratos {
		return remote.NewClient(cmd)
	})
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		cmdx.Fatalf(err.Error())
	}
}
