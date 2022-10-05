package oauth2

import (
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra/cmd"
)

func NewCreateOAuth2Client() *cobra.Command {
	return wrapHydraCmd(hydra.NewCreateClientsCommand)
}

func NewDeleteOAuth2Client() *cobra.Command {
	return wrapHydraCmd(hydra.NewDeleteClientCmd)
}

func NewGetOAuth2Client() *cobra.Command {
	return wrapHydraCmd(hydra.NewGetClientsCmd)
}

func NewImportOAuth2Client() *cobra.Command {
	return wrapHydraCmd(hydra.NewImportClientCmd)
}

func NewListOAuth2Clients() *cobra.Command {
	return wrapHydraCmd(hydra.NewListClientsCmd)
}

func NewUpdateOAuth2Client() *cobra.Command {
	return wrapHydraCmd(hydra.NewUpdateClientCmd)
}
