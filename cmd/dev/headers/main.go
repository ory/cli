// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package headers

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "headers",
		Short: "Adds language-specific headers to files",
	}
	c.AddCommand(newCopyrightCmd(), newCpCmd())
	return c
}
