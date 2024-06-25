// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package organizations

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"

	"github.com/ory/x/cmdx"
)

func NewUpdateOrganizationCmd() *cobra.Command {
	var (
		domains []string
		label   string
	)

	cmd := &cobra.Command{
		Use:   "organization <id> [--project=PROJECT_ID] [--domains=a.example.com,b.example.com] [--label=LABEL]",
		Args:  cobra.ExactArgs(1),
		Short: "Update the organization with the given ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			projectID, err := h.ProjectID()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			orgID := args[0]

			body := cloud.OrganizationBody{}
			if label != "" {
				body.Label = &label
			}
			if len(domains) > 0 {
				body.Domains = domains
			}

			organization, err := h.UpdateOrganization(cmd.Context(), projectID, orgID, body)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Organization updated successfully!")
			cmdx.PrintRow(cmd, output(*organization))
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&domains, "domains", "d", []string{}, "A list of domains that will be used for this organization.")
	cmd.Flags().StringVarP(&label, "label", "l", "", "The label of the organization.")
	client.RegisterProjectFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())

	return cmd
}
