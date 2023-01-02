// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/ory/cli/cmd/cloudx/oauth2"
	"github.com/ory/x/cmdx"
)

func NewIntrospectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "introspect",
		Short: "Introspect resources",
	}
	cmd.AddCommand(oauth2.NewIntrospectToken())

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
