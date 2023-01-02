// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/relationtuples"
	"github.com/ory/x/cmdx"
)

func NewIsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is",
		Short: "Assert the state of Ory Network resources",
	}
	cmd.AddCommand(relationtuples.NewAllowedCmd())

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())

	return cmd
}
