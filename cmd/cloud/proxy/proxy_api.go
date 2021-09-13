package proxy

import (
	"fmt"

	"github.com/ory/x/urlx"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/x/flagx"
)

func NewProxyAPICmd() *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "api",
		Short: "Proxy Ory's APIs.",
		Args:  cobra.ExactArgs(0),
		Long: fmt.Sprintf(`This command starts a proxy for Ory's APIs without reverse proxying anything else.

	$ ory proxy api --port 4000`),
		RunE: func(cmd *cobra.Command, args []string) error {
			port := flagx.MustGetInt(cmd, PortFlag)
			proto := "http"
			isHTTP := flagx.MustGetBool(cmd, WithoutHTTPSFlag)
			if !isHTTP {
				proto = "https"
			}
			conf := &config{
				noUpstream:      true,
				port:            flagx.MustGetInt(cmd, PortFlag),
				noCert:          true,
				noOpen:          true,
				apiEndpoint:     flagx.MustGetString(cmd, remote.FlagAPIEndpoint),
				consoleEndpoint: flagx.MustGetString(cmd, remote.FlagConsoleAPI),
				noJWT:           true,
				noHTTPS:         isHTTP,
				isLocal:         false,
				upstream:        fmt.Sprintf("%s://localhost:%d", proto, port),
				hostPort:        fmt.Sprintf("localhost:%d", port),
				selfURL:         urlx.ParseOrPanic(fmt.Sprintf("%s://localhost:%d", proto, port)),
			}

			return run(cmd, conf)
		},
	}

	proxyCmd.Flags().Int(PortFlag, portFromEnv(), "The port the proxy should listen on.")
	proxyCmd.Flags().Bool(WithoutHTTPSFlag, false, "Run the proxy without HTTPS. Useful if you have TLS termination or are handling HTTPS otherwise.")
	remote.RegisterClientFlags(proxyCmd.PersistentFlags())
	return proxyCmd
}
