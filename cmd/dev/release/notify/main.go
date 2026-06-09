// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package notify

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "notify",
		Short: "Notify subscribers about new releases",
	}
	c.AddCommand(newSendCmd(), newDraftCmd())
	return c
}
