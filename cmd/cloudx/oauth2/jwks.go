// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	hydra "github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
)

func wrapHydraCmd(newCmd func() *cobra.Command) *cobra.Command {
	c := newCmd()
	client.RegisterProjectFlag(c.Flags())
	client.RegisterWorkspaceFlag(c.Flags())
	cmdx.RegisterFormatFlags(c.Flags())
	// The wrapped Hydra commands resolve their endpoint through
	// cmdx.NewClient, which reads this flag. Register it here rather than
	// pulling in Kratos' cliclient, which drags the whole Kratos server in.
	c.Flags().StringP(cmdx.FlagEndpoint, cmdx.FlagEndpoint[:1], "", "The API URL this command should target. Alternatively set using the ORY_SDK_URL environmental variable.")
	return c
}

func NewGetJWK() *cobra.Command {
	return wrapHydraCmd(hydra.NewGetJWKSCmd)
}

func NewImportJWK() *cobra.Command {
	return wrapHydraCmd(hydra.NewKeysImportCmd)
}

func NewCreateJWK() *cobra.Command {
	return wrapHydraCmd(hydra.NewCreateJWKSCmd)
}

func NewDeleteJWKs() *cobra.Command {
	return wrapHydraCmd(hydra.NewDeleteJWKSCommand)
}
