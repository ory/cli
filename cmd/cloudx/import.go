// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/x/cmdx"
)

func NewImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import resources",
	}

	cmd.AddCommand(
		identity.NewImportIdentityCmd(),
		oauth2.NewImportOAuth2Client(),
		oauth2.NewImportJWK(),
	)

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	client.RegisterAuthHelpers(cmd)
	return cmd
}
