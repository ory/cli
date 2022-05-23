//go:build nodev
// +build nodev

package cmd

import (
	"github.com/spf13/cobra"
)

var devCommands []*cobra.Command
