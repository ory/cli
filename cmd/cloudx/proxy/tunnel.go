// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"fmt"
	"net/url"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"

	"github.com/spf13/cobra"
)

func NewTunnelCommand() *cobra.Command {
	conf := config{
		isTunnel: true,
		noJWT:    true,
		open:     false,
	}

	cmd := &cobra.Command{
		Use:   "tunnel <application-url> [<tunnel-url>]",
		Short: "Mirror Ory APIs on your local machine for local development and testing",
		Args:  cobra.RangeArgs(1, 2),
		Example: `{{.CommandPath}} http://localhost:3000
`,
		Long: fmt.Sprintf(`The Ory Tunnel mirrors Ory APIs on your local machine, allowing seamless development and testing. This setup is required for features such as CORS and cookie support, making it possible for Ory and your application to share the same top-level domain during development. To use the tunnel, authentication via `+"`ORY_PROJECT_API_KEY`"+` or browser-based sign-in is required.

The Ory Tunnel command connects your application and Ory's APIs through a local HTTP server. This enables both to run on the same domain or subdomain (for example, yourapp.com, localhost), which is required for cookies to function correctly.

The first argument, `+"`application-url`"+`, points to the location of your application and will be used as the default redirect URL after successful operations like login.

Example usage:

		$ {{.CommandPath}} --project <project-id-or-slug> https://www.example.org
		$ %[1]s=<project-id-or-slug> {{.CommandPath}} http://localhost:3000

### Connecting to Ory

Before using the Ory Tunnel, ensure that you have a running Ory Network project. You can create a new project with the following command:

		$ {{.Root.Name}} create project --name "Command Line Project"

Once your project is ready, pass the project's slug to the tunnel command:

		$ {{.CommandPath}} --project <project-id-or-slug> ...
		$ %[1]s=<project-id-or-slug> {{.CommandPath}} ...

### Connecting in automated environments

To connect the Ory Tunnel in automated environments, create a Project API Key for your project, set it as an environment variable, and use the `+"`--quiet`"+` flag:

		$ %[2]s=<project-api-key> {{.CommandPath}} -q ...

This will prevent the browser window from opening.

### Local development

For local development, use:

		$ {{.CommandPath}} --project <project-id-or-slug> http://localhost:3000

### CORS

You can restrict the CORS domains using the `+"`--allowed-cors-origins`"+` flag:

		$ {{.CommandPath}} http://localhost:3000 https://app.example.com \
			--allowed-cors-origins https://www.example.org \
			--allowed-cors-origins https://api.example.org \
			--allowed-cors-origins https://www.another-app.com

Per default, CORS is enabled for all origins.

### Running behind a gateway (development only)

Important: The Ory Tunnel is designed for development purposes only and should not be used in production environments.

If you need to run the tunnel behind a gateway during development, you can specify the optional second argument, tunnel-url, to define the domain where the Ory Tunnel will run (for example, https://ory.example.org).

Example:

		$ {{.CommandPath}} --project <project-id-or-slug> \
			https://www.example.org \
			https://auth.example.org \
			--cookie-domain example.org

Note: You cannot set a path in the `+"`tunnel-url`"+`.

### Ports

By default, the tunnel listens on port 4000. To change the port, use the --port flag:

		$ {{.CommandPath}} --port 8080 --project <project-id-or-slug> https://www.example.org

If your application runs on a non-standard HTTP or HTTPS port, include the port in the `+"`application-url`"+`:

		$ {{.CommandPath}} --project <project-id-or-slug> https://example.org:1234

### Cookies

For cookie support, set the `+"`--cookie-domain`"+` flag to your top-level domain:

		$ {{.CommandPath}} --project <project-id-or-slug> \
			--cookie-domain example.org \
			https://www.example.org \
			https://auth.example.org

### Redirects

To specify a custom redirect URL, use the `+"`--default-redirect-url`"+` flag:

$ {{.CommandPath}} --project <project-id-or-slug> \
	--default-redirect-url /welcome \
	https://www.example.org`, client.ProjectKey, client.ProjectAPIKey),

		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}
			selfURLString := fmt.Sprintf("http://localhost:%d", conf.port)
			if len(args) == 2 {
				selfURLString = args[1]
			}

			selfURL, err := url.ParseRequestURI(selfURLString)
			if err != nil {
				return err
			}
			conf.publicURL = selfURL

			appURL, err := url.ParseRequestURI(args[0])
			if err != nil {
				return err
			}
			if conf.defaultRedirectTo.String() == "" {
				conf.defaultRedirectTo.URL = *appURL
			}

			return runReverseProxy(cmd.Context(), h, cmd.ErrOrStderr(), &conf, "tunnel")
		},
	}

	registerConfigFlags(&conf, cmd.Flags())
	client.RegisterConfigFlag(cmd.Flags())
	client.RegisterYesFlag(cmd.Flags())
	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	cmdx.RegisterNoiseFlags(cmd.Flags())

	return cmd
}
