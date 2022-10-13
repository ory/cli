// Copyright © 2022 Ory Corp

package migration

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "migration",
	Short: "Helpers for working with gobuffalo/pop SQL migration",
}
