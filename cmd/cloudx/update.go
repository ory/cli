// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/eventstreams"
	"github.com/ory/cli/cmd/cloudx/oauth2"
	"github.com/ory/cli/cmd/cloudx/organizations"
	"github.com/ory/cli/cmd/cloudx/project"
	"github.com/ory/x/cmdx"
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update resources",
	}
	cmd.AddCommand(
		project.NewProjectsUpdateCmd(),
		project.NewUpdateIdentityConfigCmd(),
		project.NewUpdateOAuth2ConfigCmd(),
		project.NewUpdatePermissionConfigCmd(),
		project.NewUpdateNamespaceConfigCmd(),
		oauth2.NewUpdateOAuth2Client(),
		organizations.NewUpdateOrganizationCmd(),
		eventstreams.NewUpdateEventStreamCmd(),
	)

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	client.RegisterAuthHelpers(cmd)

	return cmd
}
