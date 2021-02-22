package cmd

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"

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
		return remote.NewAdminClient(cmd)
	})
	ctx = context.WithValue(ctx, cliclient.HTTPClientContextKey, func(cmd *cobra.Command) *http.Client {
		return remote.NewHTTPClient(cmd)
	})
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}
