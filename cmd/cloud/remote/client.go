package remote

import (
	"fmt"
	"os"

	"github.com/ory/kratos-client-go/client"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/stringsx"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FlagProject   = "project"
	projectEnvKey = "ORY_CLOUD_PROJECT_ID"
)

func NewClient(cmd *cobra.Command) *client.OryKratos {
	project := stringsx.Coalesce(flagx.MustGetString(cmd, FlagProject), os.Getenv(projectEnvKey))
	if project == "" {
		cmdx.Fatalf("You have to set the Ory Cloud Project ID, try --help for details.")
	}

	return client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     fmt.Sprintf("%s.tenants.staging.oryapis.dev", project),
		BasePath: "/api/kratos/admin/",
		Schemes:  []string{"https"},
	})
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.StringP(FlagProject, FlagProject[:1], "", fmt.Sprintf("Set your ORY Cloud Project ID. Alternatively set using the %s environmental variable.", projectEnvKey))
}
