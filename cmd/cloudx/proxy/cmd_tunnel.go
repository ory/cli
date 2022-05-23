package proxy

import (
	"fmt"
	"github.com/ory/x/corsx"
	"net/url"

	"github.com/ory/x/stringsx"

	"github.com/spf13/cobra"

	"github.com/ory/x/flagx"
)

func NewTunnelCommand(self string, version string) *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "tunnel application-url [tunnel-url]",
		Short: fmt.Sprintf("Tunnel Ory on a subdomain of your app or a seperate port your app's domain"),
		Args:  cobra.RangeArgs(1, 2),
		Example: fmt.Sprintf(`%[1]s tunnel http://localhost:3000 --dev
%[1]s tunnel https://app.example.com \
	--allowed-cors-origins https://www.example.org \
	--allowed-cors-origins https://api.example.org \
	--allowed-cors-origins https://www.another-app.com
`, self),
		Long: fmt.Sprintf(`
This command runs an HTTP Server which is connected to Ory's APIs, in order for your application and Ory's
APIs to run on the same top level domain (for example yourapp.com, localhost). Having Ory on your domain
is required for cookies to work.

The first argument `+"`"+`application-url`+"`"+` points to the location of your application. This location
will be used as the default redirect URL for the tunnel, for example after a successful login.

    $ %[1]s tunnel https://www.example.org
    $ %[1]s tunnel http://localhost:3000

### Connecting to Ory

Before you start, you need to have a running instance of Ory. Set the environment variable
ORY_SDK_URL to the path where Ory is available. For Ory Cloud, this is the "SDK URL"
which can be found in the "API & Services" section of your Ory Cloud Console.

	$ export ORY_SDK_URL=https://playground.projects.oryapis.com

Alternatively, you can set this using the --sdk-url flag:

	$ %[1]s proxy --sdk-url https://playground.projects.oryapis.com \
		...

### Developing Locally

When developing locally we recommend to use the `+"`"+`--dev`+"`"+` flag, which uses a lax security setting:

    $ %[1]s tunnel --dev \
		http://localhost:3000

### Running on a Server

To go to production set up a custom domain (CNAME) for Ory. If you can not set up a custom
domain - for example because you are developing a staging environment - using the Ory Tunnel is an alternative.

To run on a server, you need to set the optional second argument  `+"`"+`[tunnel-url]`+"`"+`. It tells the Ory Tunnel
on which domain it will run (for example https://ory.example.org).

	$ %[1]s tunnel \
		https://www.example.org \
		https://auth.example.org \
		--cookie-domain example.org \
		--allowed-cors-origins https://www.example.org \
		--allowed-cors-origins https://api.example.org

Please note that you can not set a path in the `+"`"+`[tunnel-url]`+"`"+`!

### Ports

Per default, the tunnel listens on port 4000. If you want to listen on another port, use the
port flag:

	$ %[1]s tunnel --port 8080 \
		https://www.example.org

If your application URL is available on a non-standard HTTP/HTTPS port, you can set that port in the `+"`"+`application-url`+"`"+`:

	$ %[1]s tunnel \
		https://example.org:1234

### Cookies

We recommend setting the `+"`"+`--cookie-domain`+"`"+` value to your top level domain:

	$ %[1]s tunnel \
		--cookie-domain example.org \
		https://www.example.org \
		https://auth.example.org

### Redirects

TO use a different default redirect URL, use the `+"`"+`--default-redirect-url`+"`"+` flag:

    $ %[1]s tunnel --default-redirect-url /welcome \
		https://www.example.org
`, self),

		RunE: func(cmd *cobra.Command, args []string) error {
			port := flagx.MustGetInt(cmd, PortFlag)
			selfURLString := fmt.Sprintf("http://localhost:%d", port)
			if len(args) == 2 {
				selfURLString = args[1]
			}

			selfURL, err := url.ParseRequestURI(selfURLString)
			if err != nil {
				return err
			}

			redirectURL, err := url.Parse(stringsx.Coalesce(flagx.MustGetString(cmd, DefaultRedirectURLFlag), args[0]))
			if err != nil {
				return err
			}

			oryURL, err := getEndpointURL(cmd)
			if err != nil {
				return err
			}

			appURL, err := url.ParseRequestURI(args[0])
			if err != nil {
				return err
			}

			origins, err := corsx.NormalizeOriginStrings(append(
				flagx.MustGetStringSlice(cmd, CORSFlag), appURL.String()),
			)
			if err != nil {
				return err
			}

			conf := &config{
				port:              flagx.MustGetInt(cmd, PortFlag),
				noJWT:             true,
				noOpen:            true,
				upstream:          oryURL.String(),
				cookieDomain:      flagx.MustGetString(cmd, CookieDomainFlag),
				publicURL:         selfURL,
				oryURL:            oryURL,
				pathPrefix:        "",
				isTunnel:          true,
				defaultRedirectTo: redirectURL,
				isDev:             flagx.MustGetBool(cmd, DevFlag),
				isDebug:           flagx.MustGetBool(cmd, DebugFlag),
				corsOrigins:       origins,
			}

			return run(cmd, conf, version, "cloud")
		},
	}

	proxyCmd.Flags().String(CookieDomainFlag, "", "Set a dedicated cookie domain.")
	proxyCmd.Flags().String(ServiceURL, "", "Set the Ory SDK URL.")
	proxyCmd.Flags().Int(PortFlag, portFromEnv(), "The port the proxy should listen on.")
	proxyCmd.Flags().Bool(DevFlag, false, "Use this flag when developing locally.")
	proxyCmd.Flags().Bool(DebugFlag, false, "Use this flag to debug, for example, CORS requests.")
	proxyCmd.Flags().String(DefaultRedirectURLFlag, "", "Set the URL to redirect to per default after e.g. login or account creation.")
	proxyCmd.Flags().StringSlice(CORSFlag, []string{}, "A list of allowed CORS origins. Wildcards are allowed.")
	return proxyCmd
}
