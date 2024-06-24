// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/oauth2"
	"github.com/ory/x/cmdx"
)

func NewRevokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke resources",
	}
	cmd.AddCommand(oauth2.NewRevokeToken())

	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	client.RegisterAuthHelpers(cmd)
	return cmd
}
