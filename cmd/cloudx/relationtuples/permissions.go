// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package relationtuples

import (
	"github.com/spf13/cobra"

	"github.com/ory/keto/cmd/check"
)

func NewAllowedCmd() *cobra.Command {
	cmd := check.NewCheckCmd()
	wrapForOryCLI(cmd)
	cmd.Use = "allowed <subject> <relation> <namespace> <object>"

	return cmd
}
