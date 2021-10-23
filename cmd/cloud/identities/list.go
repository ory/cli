package identities

import (
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/kratos/cmd/identities"
)

func init() {
	listCmd := identities.NewListCmd()
	remote.RegisterClientFlags(listCmd.PersistentFlags())
}
