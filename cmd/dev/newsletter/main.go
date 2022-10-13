// Copyright © 2022 Ory Corp

package newsletter

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "newsletter",
	Short: "Draft and send release newsletters using Mailchimp",
}
