package proxy

import (
	"fmt"
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/x/flagx"
	"github.com/spf13/cobra"
)

func NewProxyProductionCmd() *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "production [upstream] [domain]",
		Short: "Run an Ory integrated application in a remote environment",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(`This command starts a reverse proxy which can be deployed in front of your application.

This command is targeted for remote, hosted, internet-facing applications. If you wish to develop an application locally,
please use "ory proxy local" instead.

To require authentication before accessing paths in your application, use the --%[1]s flag:

	$ ory proxy remote --port 4000 --%[1]s /members --%[1]s /admin \
		http://localhost:3000 \
		your-domain.com

%[2]s`, ProtectPathsFlag, jwtHelp),
		/*
		   The --%s values support regular expression templating, meaning that you can use regular expressions within "<>":

		   	$ ory proxy http://localhost:3000 --allow --%s "http://localhost:3000/<(login|dashboard)>" --%s "http://localhost:3000/<([0-9]{3})>"

		   The supported Regular Expression Syntax is RE2 and documented at: https://golang.org/pkg/regexp/
		   To test your Regular Expression, head over to https://regex101.com and select "Golang" on the left.
		*/
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := &config{
				port:                flagx.MustGetInt(cmd, PortFlag),
				protectPathPrefixes: flagx.MustGetStringSlice(cmd, ProtectPathsFlag),
				noCert:              true,
				noOpen:              true,
				apiEndpoint:         flagx.MustGetString(cmd, remote.FlagAPIEndpoint),
				consoleEndpoint:     flagx.MustGetString(cmd, remote.FlagConsoleAPI),
				upstream:            args[0],
				domain:              args[1],
				isLocal:             false,
			}

			return run(cmd, conf)
		},
	}

	proxyCmd.Flags().Int(PortFlag, portFromEnv(), "The port the proxy should listen on.")
	proxyCmd.Flags().StringSlice(ProtectPathsFlag, []string{}, "Require authentication before accessing these paths.")
	remote.RegisterClientFlags(proxyCmd.PersistentFlags())
	return proxyCmd
}
