// Copyright © 2022 Ory Corp

package oauth2

import (
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra/cmd"
)

func NewPerformAuthorizationCode() *cobra.Command {
	return wrapHydraCmd(hydra.NewPerformAuthorizationCodeCmd)
}

func NewPerformClientCredentials() *cobra.Command {
	return wrapHydraCmd(hydra.NewPerformClientCredentialsCmd)
}
