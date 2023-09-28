// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package organizations

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

func NewCreateOrganizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization label [--project=PROJECT_ID] [--domains=a.example.com,b.example.com]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a new Ory Network organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := client.ProjectOrDefault(cmd, h)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			label := args[0]
			domains := flagx.MustGetStringSlice(cmd, "domains")

			organization, err := h.CreateOrganization(projectID, cloud.OrganizationBody{
				Label:   &label,
				Domains: domains,
			})
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Organization created successfully!")
			cmdx.PrintRow(cmd, output(*organization))
			return nil
		},
	}

	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	cmd.Flags().StringSliceP("domains", "d", []string{}, "A list of domains that will be used for this organization.")

	return cmd
}
