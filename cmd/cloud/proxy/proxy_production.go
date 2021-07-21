package proxy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/x/flagx"
)

func NewProxyProductionCmd() *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "production [upstream] [host]",
		Short: "Run an application in production mode with Ory integration",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(`This command starts a reverse proxy which can be deployed in front of your application. This command works for remote environments,
for example when deploying a React, NodeJS, Java, PHP, ... app to a server / the cloud.

	$ ory proxy remote --port 4000 \
		http://localhost:3000 \
		example.org

If you want to expose the application / proxy at a specific port, append the port to the domain name:

	$ ory proxy remote --port 4000 \
		http://127.0.0.1:3000 \
		example.org:8080

%s
`, jwtHelp),
		/*
		   The --%s values support regular expression templating, meaning that you can use regular expressions within "<>":

		   	$ ory proxy http://localhost:3000 --allow --%s "http://localhost:3000/<(login|dashboard)>" --%s "http://localhost:3000/<([0-9]{3})>"

		   The supported Regular Expression Syntax is RE2 and documented at: https://golang.org/pkg/regexp/
		   To test your Regular Expression, head over to https://regex101.com and select "Golang" on the left.
		*/
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := &config{
				port:            flagx.MustGetInt(cmd, PortFlag),
				noCert:          true,
				noOpen:          true,
				apiEndpoint:     flagx.MustGetString(cmd, remote.FlagAPIEndpoint),
				consoleEndpoint: flagx.MustGetString(cmd, remote.FlagConsoleAPI),
				isLocal:         false,
				upstream:        args[0],
				hostPort:        args[1],
			}

			return run(cmd, conf)
		},
	}

	proxyCmd.Flags().Int(PortFlag, portFromEnv(), "The port the proxy should listen on.")
	remote.RegisterClientFlags(proxyCmd.PersistentFlags())
	return proxyCmd
}
