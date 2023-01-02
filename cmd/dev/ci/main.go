// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package ci

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/ci/deps"
	"github.com/ory/cli/cmd/dev/ci/github"
	"github.com/ory/cli/cmd/dev/ci/monorepo"
	"github.com/ory/cli/cmd/dev/ci/orbs"
)

var Main = &cobra.Command{
	Use:   "ci",
	Short: "Continuous Integration helpers",
}

func init() {
	Main.AddCommand(orbs.Main)
	Main.AddCommand(github.Main)
	Main.AddCommand(monorepo.Main)
	Main.AddCommand(deps.Main)
}
