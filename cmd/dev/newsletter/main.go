// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package newsletter

import (
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{
	Use:   "newsletter",
	Short: "Draft and send release newsletters using Mailchimp",
}
