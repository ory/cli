// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	hydra "github.com/ory/hydra/cmd"
	"github.com/ory/x/cmdx"
)

func wrapHydraCmd(newCmd func() *cobra.Command) *cobra.Command {
	c := newCmd()
	client.RegisterProjectFlag(c.Flags())
	cmdx.RegisterFormatFlags(c.Flags())
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
