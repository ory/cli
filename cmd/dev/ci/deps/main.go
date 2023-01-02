// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package deps

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "deps",
	Short: "Helpers for binary dependencies in Makefiles.",
}

func init() {
	Main.PersistentFlags().StringP("os", "o", "", "OS the binary should run on. Currently only 'linux' and 'darwin' are supported.")
	Main.PersistentFlags().StringP("architecture", "a", "", "Architecture the binary should run on. Currently only 'amd64' is supported.")
	Main.PersistentFlags().StringP("config", "c", "", "Path to config files.")
}
