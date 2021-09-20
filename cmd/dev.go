//go:build !nodev
// +build !nodev

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev"
)

var devCommand = []*cobra.Command{dev.Main}
