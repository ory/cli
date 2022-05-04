package cmd

import (
	"context"
	"fmt"
	"github.com/ory/kratos/cmd/jsonnet"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/cli/buildinfo"
	"github.com/ory/cli/cmd/cloudx"
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
			jsonnet.NewFormatCmd(),
			jsonnet.NewLintCmd(),
			cloudx.NewDeleteCmd(c),
			cloudx.NewGetCmd(c),
			cloudx.NewListCmd(c),
			cloudx.NewImportCmd(c),
			cloudx.NewPatchCmd(),
			cloudx.NewProxyCommand("ory", buildinfo.Version),
			cloudx.NewTunnelCommand("ory", buildinfo.Version),
			cloudx.NewUpdateCmd(),
			cloudx.NewValidateCmd(),
			versionCmd,
		)...,
	)

	return c
}

func Execute() {
	ctx := cloudx.ContextWithClient(context.Background())
	rootCmd := NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}
