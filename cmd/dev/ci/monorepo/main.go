package monorepo

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "monorepo",
	Short: "Helpers for CircleCI Monorepo Support",
}
