// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package notify

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/newsletter"
	"github.com/ory/x/flagx"
)

func newSendCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "send [list-id]",
		Short: "Send the release notification Mailchimp campaign",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			newsletter.SendCampaign(args[0], flagx.MustGetBool(cmd, "dry"))
			return nil
		},
	}
	c.Flags().Bool("dry", false, "Do not actually send the campaign")
	return c
}
