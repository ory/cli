// Copyright © 2022 Ory Corp

package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"
	"github.com/ory/x/cmdx"
)

func NewIntrospectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "introspect",
		Short: "Introspect resources",
	}
	cmd.AddCommand(oauth2.NewIntrospectToken())

	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
