// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build nodev
// +build nodev

package cmd

import (
	"github.com/spf13/cobra"
)

var devCommands []*cobra.Command
