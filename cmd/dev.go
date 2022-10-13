// Copyright © 2022 Ory Corp

//go:build !nodev
// +build !nodev

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev"
)

var devCommands = []*cobra.Command{dev.Main}
