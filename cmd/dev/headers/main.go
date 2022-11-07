// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package headers

import "github.com/spf13/cobra"

var Main = &cobra.Command{
	Use:   "headers",
	Short: "Adds language-specific headers to files",
}
