// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package organizations

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go/1.2"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

func NewUpdateOrganizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization id [--project=PROJECT_ID] [--domains=a.example.com,b.example.com] [--label=LABEL]",
		Args:  cobra.ExactArgs(1),
		Short: "Update the organization with the given ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := client.ProjectOrDefault(cmd, h)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			orgID := args[0]

			body := cloud.OrganizationBody{}
			if l := flagx.MustGetString(cmd, "label"); l != "" {
				body.Label = &l
			}
			if domains := flagx.MustGetStringSlice(cmd, "domains"); len(domains) > 0 {
				body.Domains = domains
			}

			organization, err := h.UpdateOrganization(projectID, orgID, body)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Organization updated successfully!")
			cmdx.PrintRow(cmd, output(*organization))
			return nil
		},
	}

	cmd.Flags().StringSliceP("domains", "d", []string{}, "A list of domains that will be used for this organization.")
	cmd.Flags().StringP("label", "l", "", "The label of the organization.")
	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())

	return cmd
}
