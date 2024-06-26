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
	}

	cmd := &cobra.Command{
		Use:   "tunnel <application-url> [<tunnel-url>]",
		Short: "Tunnel Ory on a subdomain of your app or a separate port your app's domain",
		Args:  cobra.RangeArgs(1, 2),
		Example: `{{.CommandPath}} http://localhost:3000 --dev
{{.CommandPath}} https://app.example.com \
	--allowed-cors-origins https://www.example.org \
	--allowed-cors-origins https://api.example.org \
	--allowed-cors-origins https://www.another-app.com
`,
		Long: fmt.Sprintf(`Tunnels Ory APIs on a subdomain or separate port of your app. This command runs an HTTP Server which is connected to Ory's APIs, in order for your application and Ory's
APIs to run on the same top level domain (for example yourapp.com, localhost). Having Ory on your domain
is required for cookies to work.

The first argument `+"`application-url`"+` points to the location of your application. This location
will be used as the default redirect URL for the tunnel, for example after a successful login.

$ {{.CommandPath}} --project <project-id-or-slug> https://www.example.org
$ %[1]s=<project-id-or-slug> {{.CommandPath}} http://localhost:3000

### Connecting to Ory

Before you start, you need to have a running Ory Network project. You can create one with the following command:

	$ {{.Root.Name}} create project --name "Command Line Project"

Pass the project's slug as a flag to the tunnel command:

	$ {{.CommandPath}} --project <project-id-or-slug> ...
	$ %[1]s=<project-id-or-slug> {{.CommandPath}} tunnel ...

### Developing Locally

When developing locally we recommend to use the `+"`--dev`"+` flag, which uses a lax security setting:

	$ {{.CommandPath}} --dev --project <project-id-or-slug> \
		http://localhost:3000

### Running behind a Gateway

To go to production set up a custom domain (CNAME) for Ory.

If you need to run the tunnel behind a gateway, you have to set the optional second argument `+"`tunnel-url`"+`. It tells the Ory Tunnel
on which domain it will run (for example https://ory.example.org).

	$ {{.CommandPath}} --project <project-id-or-slug> \
		https://www.example.org \
		https://auth.example.org \
		--cookie-domain example.org \
		--allowed-cors-origins https://www.example.org \
		--allowed-cors-origins https://api.example.org

Please note that you can not set a path in the `+"`[tunnel-url]`"+`!

### Ports

Per default, the tunnel listens on port 4000. If you want to listen on another port, use the
port flag:

	$ {{.CommandPath}} --port 8080 --project <project-id-or-slug> \
		https://www.example.org

If your application URL is available on a non-standard HTTP/HTTPS port, you can set that port in the `+"`application-url`"+`:

	$ {{.CommandPath}} --project <project-id-or-slug> \
		https://example.org:1234

### Cookies

We recommend setting the `+"`--cookie-domain`"+` value to your top level domain:

	$ {{.CommandPath}} --project <project-id-or-slug> \
		--cookie-domain example.org \
		https://www.example.org \
		https://auth.example.org

### Redirects

To use a different default redirect URL, use the `+"`--default-redirect-url`"+` flag:

	$ {{.CommandPath}} tunnel --project <project-id-or-slug> \
		--default-redirect-url /welcome \
		https://www.example.org
`, client.ProjectKey),

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

			return runReverseProxy(cmd.Context(), h, cmd.ErrOrStderr(), &conf, "tunnel")
		},
	}

	registerConfigFlags(&conf, cmd.Flags())
	client.RegisterConfigFlag(cmd.Flags())
	client.RegisterYesFlag(cmd.Flags())
	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	cmdx.RegisterNoiseFlags(cmd.Flags())

	cmdx.EnableUsageTemplating(cmd)
	return cmd
}
