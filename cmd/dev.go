// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build !nodev
// +build !nodev

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev"
)

func newDevCommands() []*cobra.Command {
	return []*cobra.Command{dev.NewCommand()}
}
