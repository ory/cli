package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/cli/buildinfo"
	"github.com/ory/x/cloudx"
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
			cloudx.NewAuthCmd(),
			cloudx.NewCreateCmd(),
			cloudx.NewGetCmd(),
			cloudx.NewListCmd(),
			cloudx.NewPatchCmd(),
			cloudx.NewUpdateCmd(),
			cloudx.NewProxyCommand("ory", buildinfo.Version),
			cloudx.NewTunnelCommand("ory", buildinfo.Version),
			versionCmd,
		)...,
	)

	return c
}

func Execute() {
	ctx := context.Background()

	rootCmd := NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}
