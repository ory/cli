// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package migration

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "migration",
	Short: "Helpers for working with gobuffalo/pop SQL migration",
}
