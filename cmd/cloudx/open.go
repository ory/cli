// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/accountexperience"
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open Ory Account Experience Pages",
	}
	cmd.AddCommand(accountexperience.NewAccountExperienceOpenCmd())
	client.RegisterProjectFlag(cmd.PersistentFlags())
	client.RegisterConfigFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())

	return cmd
}
