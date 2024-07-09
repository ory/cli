// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"fmt"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Signs you out of your account on this computer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}
			if err := h.SignOut(cmd.Context()); err != nil {
				return err
			}
			fmt.Println("You signed out successfully.")
			return nil
		},
	}
	client.RegisterConfigFlag(cmd.PersistentFlags())
	return cmd
}
