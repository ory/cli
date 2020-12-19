package identities

import (
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/kratos/cmd/identities"
)

func init() {
	remote.RegisterClientFlags(identities.ListCmd.PersistentFlags())
}
