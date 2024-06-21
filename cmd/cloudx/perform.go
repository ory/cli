// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewPerformCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perform",
		Short: "Perform a flow",
	}

	cmd.AddCommand(
		oauth2.NewPerformAuthorizationCode(),
		oauth2.NewPerformClientCredentials(),
	)

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	client.RegisterAuthHelpers(cmd)

	return cmd
}
