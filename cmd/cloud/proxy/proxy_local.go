package proxy

import (
	"fmt"

	"github.com/ory/x/urlx"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/x/flagx"
)

func NewProxyLocalCmd() *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "local [upstream]",
		Short: "Develop an application locally and integrate it with Ory",
		Args:  cobra.ExactArgs(1),
		Long: fmt.Sprintf(`This command starts a reverse proxy which can be deployed in front of your application. This works best on local (your computer) environments, for example when developing a React, NodeJS, Java, PHP app.

	$ ory proxy local --port 4000 \
		http://localhost:3000

%s
`, jwtHelp),
		/*
		   The --%s values support regular expression templating, meaning that you can use regular expressions within "<>":

		   	$ ory proxy http://localhost:3000 --allow --%s "http://localhost:3000/<(login|dashboard)>" --%s "http://localhost:3000/<([0-9]{3})>"

		   The supported Regular Expression Syntax is RE2 and documented at: https://golang.org/pkg/regexp/
		   To test your Regular Expression, head over to https://regex101.com and select "Golang" on the left.
		*/
		RunE: func(cmd *cobra.Command, args []string) error {
			port := flagx.MustGetInt(cmd, PortFlag)
			proto := "http"
			isHTTP := flagx.MustGetBool(cmd, WithoutHTTPSFlag)
			if !isHTTP {
				proto = "https"
			}
			conf := &config{
				port:            flagx.MustGetInt(cmd, PortFlag),
				noCert:          flagx.MustGetBool(cmd, NoCertInstallFlag),
				noOpen:          flagx.MustGetBool(cmd, NoOpenFlag),
				apiEndpoint:     flagx.MustGetString(cmd, remote.FlagAPIEndpoint),
				consoleEndpoint: flagx.MustGetString(cmd, remote.FlagConsoleAPI),
				noJWT:           flagx.MustGetBool(cmd, WithoutJWTFlag),
				noHTTPS:         isHTTP,
				isLocal:         true,
				upstream:        args[0],
				hostPort:        fmt.Sprintf("localhost:%d", port),
				selfURL:         urlx.ParseOrPanic(fmt.Sprintf("%s://localhost:%d", proto, port)),
			}

			return run(cmd, conf)
		},
	}

	proxyCmd.Flags().Int(PortFlag, portFromEnv(), "The port the proxy should listen on.")
	proxyCmd.Flags().Bool(NoCertInstallFlag, false, "If set will not try to add the HTTPS certificate to your certificate store.")
	proxyCmd.Flags().Bool(NoOpenFlag, false, "Do not open the browser when the proxy starts.")
	proxyCmd.Flags().Bool(WithoutJWTFlag, false, "Do not create a JWT from the Ory Kratos Session. Useful if you need fast start up times of the Ory Proxy.")
	proxyCmd.Flags().Bool(WithoutHTTPSFlag, false, "Run the proxy without HTTPS. Useful if you have TLS termination or are handling HTTPS otherwise.")
	remote.RegisterClientFlags(proxyCmd.PersistentFlags())
	return proxyCmd
}
