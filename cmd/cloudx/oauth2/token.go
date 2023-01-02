// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra/cmd"
)

func NewDeleteAccessTokens() *cobra.Command {
	return wrapHydraCmd(hydra.NewDeleteAccessTokensCmd)
}

func NewRevokeToken() *cobra.Command {
	return wrapHydraCmd(hydra.NewRevokeTokenCmd)
}

func NewIntrospectToken() *cobra.Command {
	return wrapHydraCmd(hydra.NewIntrospectTokenCmd)
}
