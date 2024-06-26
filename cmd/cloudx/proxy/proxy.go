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

func NewProxyCommand() *cobra.Command {
	conf := config{
		pathPrefix: "/.ory",
	}

	proxyCmd := &cobra.Command{
		Use:   "proxy <application-url> [<publish-url>]",
		Short: "Run your app and Ory on the same domain using a reverse proxy",
		Args:  cobra.RangeArgs(1, 2),
		Example: `{{.CommandPath}} http://localhost:3000 --dev
{{.CommandPath}} proxy http://localhost:3000 https://app.example.com \
	--allowed-cors-origins https://www.example.org \
	--allowed-cors-origins https://api.example.org \
	--allowed-cors-origins https://www.another-app.com
`,
		Long: `Allows running your app and Ory on the same domain by starting a reverse proxy that runs in front of your application.

The first argument ` + "`<application-url>`" + ` points to the location of your application. The Ory Proxy
will pass all traffic through to this URL.

    $ {{.CommandPath}} --project <project-id-or-slug> https://www.example.org
    $ ORY_PROJECT=<project-id-or-slug> {{.CommandPath}} proxy http://localhost:3000

### Connecting to Ory

Before you start, you need to have an Ory Network project. You can create one with the following command:

	$ {{.Root.Name}} create project --name "Command Line Project" --use
	$ {{.CommandPath}} ...

### Developing Locally

When developing locally we recommend to use the ` + "`--dev`" + ` flag, which uses a lax security setting:

	$ {{.CommandPath}} --dev \
		--project <project-id-or-slug> \
		http://localhost:3000

The first argument ` + "`<application-url>`" + ` points to the location of your application. If you are
running the proxy and your app on the same host, this could be localhost. All traffic arriving at the
Ory Proxy will be passed through to this URL.

The second argument ` + "`<publish-url>`" + ` is optional and only needed when going to production.
It refers to the public URL of your application (e.g. https://www.example.org).

If ` + "`<publish-url>`" + ` is not set, it will default to the
host and port the proxy listens on.

### Running behind a Gateway

To go to production set up a custom domain (CNAME) for Ory.

You must set the ` + "`<publish-url>`" + ` if you are using the Ory Proxy behind a gateway:

	$ {{.CommandPath}} \
		--project <project-id-or-slug> \
		http://localhost:3000 \
		https://gateway.local:5000

Please note that you can not set a path in the ` + "`<publish-url>`" + `!

### Ports

Per default, the proxy listens on port 4000. If you want to listen on another port, use the
port flag:

	$ {{.CommandPath}} --port 8080 --project <project-id-or-slug> \
		http://localhost:3000

### Multiple Domains

If the proxy runs on a subdomain, and you want Ory's cookies (e.g. the session cookie) to
be available on all of your domain, you can use the ` + "`--cookie-domain`" + ` flag to customize the cookie
domain. You will also need to allow your subdomains in the CORS headers:

	$ {{.CommandPath}} --project <project-id-or-slug> \
		--cookie-domain gateway.local \
		--allowed-cors-origins https://www.gateway.local \
		--allowed-cors-origins https://api.gateway.local \
		http://127.0.0.1:3000 \
		https://ory.gateway.local

### Redirects

Per default all default redirects will go to to ` + "`[<publish-url>]`" + `. You can change this behavior using
the ` + "`--default-redirect-url`" + ` flag:

    $ {{.CommandPath}} --project <project-id-or-slug> \
		--default-redirect-url /welcome \
		http://127.0.0.1:3000 \
		https://ory.example.org

Now, all redirects happening e.g. after login will point to ` + "`/welcome`" + ` instead of ` + "`/`" + ` unless you
have specified custom redirects in your Ory configuration or in the flow's ` + "`?return_to=`" + ` query parameter.

### JSON Web Token

If the request is not authenticated, the HTTP Authorization Header will be empty:

	GET / HTTP/1.1
	Host: localhost:3000

If the request was authenticated, a JSON Web Token can be sent in the HTTP Authorization Header containing the
Ory Session:

	GET / HTTP/1.1
	Host: localhost:3000
	Authorization: Bearer the-json-web-token

The JSON Web Token claims contain:

* The "sub" field which is set to the Ory Identity ID.
* The "session" field which contains the full Ory Session.

The JSON Web Token is signed using the ES256 algorithm. The public key can be found by fetching the /.ory/jwks.json path
when calling the proxy - for example: ` + "`http://127.0.0.1:4000/.ory/jwks.json`" + `

An example payload of the JSON Web Token is:

	{
	  "id": "821f5a53-a0b3-41fa-9c62-764560fa4406",
	  "active": true,
	  "expires_at": "2021-02-25T09:25:37.929792Z",
	  "authenticated_at": "2021-02-24T09:25:37.931774Z",
	  "issued_at": "2021-02-24T09:25:37.929813Z",
	  "identity": {
		"id": "18aafd3e-b00c-4b19-81c8-351e38705126",
		"schema_id": "default",
		"schema_url": "https://example.projects.oryapis.com/api/kratos/public/schemas/default",
		"traits": {
		  "email": "foo@bar",
		  // ... your other identity traits
		}
	  }
	}
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			conf.upstream = args[0]

			selfURLString := fmt.Sprintf("http://localhost:%d", conf.port)
			if len(args) == 2 {
				selfURLString = args[1]
			}

			var err error
			conf.publicURL, err = url.ParseRequestURI(selfURLString)
			if err != nil {
				return err
			}

			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			return runReverseProxy(cmd.Context(), h, cmd.ErrOrStderr(), &conf, "proxy")
		},
	}

	flags := proxyCmd.Flags()
	registerConfigFlags(&conf, flags)

	client.RegisterConfigFlag(flags)
	client.RegisterProjectFlag(flags)
	client.RegisterWorkspaceFlag(flags)
	client.RegisterYesFlag(flags)
	cmdx.RegisterNoiseFlags(flags)

	proxyCmd.Root().Name()
	return proxyCmd
}
