// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/kratos/cmd/identities"
)

func NewListIdentityCmd() *cobra.Command {
	cmd := identities.NewListIdentitiesCmd()
	client.RegisterProjectFlag(cmd.Flags())
	return cmd
}
