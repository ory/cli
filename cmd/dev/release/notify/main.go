package notify

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "notify",
	Short: "Notify subscribers about new releases",
}
