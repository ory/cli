// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx"
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/proxy"
	"github.com/ory/kratos/cmd/jsonnet"
	"github.com/ory/x/cmdx"
)

func NewRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "ory",
		Short: "The Ory CLI",
	}

	c.AddCommand(devCommands...)
	c.AddCommand(
		cloudx.NewAuthCmd(),
		cloudx.NewCreateCmd(),
		jsonnet.NewFormatCmd(),
		jsonnet.NewLintCmd(),
		cloudx.NewDeleteCmd(),
		cloudx.NewGetCmd(),
		cloudx.NewUseCmd(),
		cloudx.NewListCmd(),
		cloudx.NewImportCmd(),
		cloudx.NewOpenCmd(),
		cloudx.NewPatchCmd(),
		cloudx.NewParseCmd(),
		cloudx.NewPerformCmd(),
		proxy.NewProxyCommand(),
		proxy.NewTunnelCommand(),
		cloudx.NewUpdateCmd(),
		cloudx.NewValidateCmd(),
		cloudx.NewRevokeCmd(),
		cloudx.NewIntrospectCmd(),
		cloudx.NewIsCmd(),
		versionCmd,
	)
	cmdx.EnableUsageTemplating(c)

	return c
}

func Execute() {
	ctx := client.ContextWithClient(context.Background())
	rootCmd := NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}
