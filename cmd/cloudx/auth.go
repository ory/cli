// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Create a new Ory Network account or sign in to an existing account.",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			ac, err := h.GetAuthenticatedConfig(cmd.Context())
			if err != nil {
				return err
			}

			cmdx.PrintRow(cmd, ac)
			return nil
		},
	}
	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	cmd.AddCommand(NewLogoutCmd())
	return cmd
}
