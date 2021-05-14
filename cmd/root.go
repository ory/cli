package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ory/x/flagx"

	"github.com/pkg/errors"

	"github.com/ory/cli/cmd/cloud/identities"
	"github.com/ory/cli/cmd/cloud/proxy"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloud/remote"
	kratos "github.com/ory/kratos-client-go"
	"github.com/ory/kratos/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "ory",
		Short: "The ORY CLI",
	}

	c.AddCommand(
		append(
			devCommand,
			identities.Main,
			proxy.Main,
			versionCmd,
		)...,
	)

	return c
}

func Execute() {
	ctx := context.WithValue(context.Background(), cliclient.ClientContextKey, func(cmd *cobra.Command) *kratos.APIClient {
		return remote.NewAdminClient(flagx.MustGetString(cmd, remote.FlagAPIEndpoint), flagx.MustGetString(cmd, remote.FlagConsoleAPI))
	})

	rootCmd := NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}
