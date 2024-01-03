// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/corsx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/stringsx"
)

func NewProxyCommand(self string, version string) *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "proxy application-url [publish-url]",
		Short: "Run your app and Ory on the same domain using a reverse proxy",
		Args:  cobra.RangeArgs(1, 2),
		Example: fmt.Sprintf(`%[1]s proxy http://localhost:3000 --dev
%[1]s proxy http://localhost:3000 https://app.example.com \
	--allowed-cors-origins https://www.example.org \
	--allowed-cors-origins https://api.example.org \
	--allowed-cors-origins https://www.another-app.com
`, self),
		Long: fmt.Sprintf(`Allows running your app and Ory on the same domain by starting a reverse proxy that runs in front of your application.
This proxy works both in development and in production, for example when deploying a
React, NodeJS, Java, PHP, ... app to a server / the cloud or when developing it locally
on your machine.

The first argument `+"`"+`application-url`+"`"+` points to the location of your application. The Ory Proxy
will pass all traffic through to this URL.

    $ %[1]s proxy --project <your-project-slug> https://www.example.org
    $ ORY_PROJECT_SLUG=<your-project-slug> %[1]s proxy http://localhost:3000

### Connecting to Ory

Before you start, you need to have a running Ory Network project. You can create one with the following command:

	$ %[1]s create project --name "Command Line Project"

Pass the project's slug as a flag to the proxy command:

	$ %[1]s proxy --project <your-project-slug> ...
	$ ORY_PROJECT_SLUG=<your-project-slug> %[1]s proxy ...

When using the `+"`"+`ORY_SDK_URL`+"`"+` or `+"`"+`ORY_KRATOS_URL`+"`"+` to point to a custom domain on the project instead of the `+"`"+`ORY_PROJECT_SLUG`+"`"+` environment variable,
take care that the project has not set the custom UI base URL on this domain. This will cause the browser to always redirect to the custom UI base URL instead
of the configured `+"`"+`application-url`+"`"+`.

### Developing Locally

When developing locally we recommend to use the `+"`"+`--dev`+"`"+` flag, which uses a lax security setting:

	$ %[1]s proxy --dev \
		--project <your-project-slug> \
		http://localhost:3000

The first argument `+"`"+`application-url`+"`"+` points to the location of your application. If you are
running the proxy and your app on the same host, this could be localhost. All traffic arriving at the
Ory Proxy will be passed through to this URL.

The second argument `+"`"+`[publish-url]`+"`"+` is optional and only needed when going to production.
It refers to the public URL of your application (e.g. https://www.example.org).

If `+"`"+`[publish-url]`+"`"+` is not set, it will default to the default
host and port this proxy listens on:

	http://localhost:4000

### Running on a Server

To go to production set up a custom domain (CNAME) for Ory. If you can not set up a custom
domain - for example because you are developing a staging environment - using the Ory Proxy is an alternative.

You must set the `+"`"+`[publish-url]`+"`"+` if you are not using the Ory Proxy in locally or in
development:

	$ %[1]s proxy \
		--project <your-project-slug> \
		http://localhost:3000 \
		https://example.org

Please note that you can not set a path in the `+"`"+`[publish-url]`+"`"+`!

### Ports

Per default, the proxy listens on port 4000. If you want to listen on another port, use the
port flag:

	$ %[1]s proxy --port 8080  --project <your-project-slug> \
		http://localhost:3000 \
		https://example.org

If your public URL is available on a non-standard HTTP/HTTPS port, you can set that port in the `+"`"+`[publish-url]`+"`"+`:

	$ %[1]s proxy --project <your-project-slug> \
		http://localhost:3000 \
		https://example.org:1234

### Multiple Domains

If this proxy runs on a subdomain, and you want Ory's cookies (e.g. the session cookie) to
be available on all of your domain, you can use the following CLI flag to customize the cookie
domain. You will also need to allow your subdomains in the CORS headers:

	$ %[1]s proxy --project <your-project-slug> \
		--cookie-domain example.org \
		--allowed-cors-origins https://www.example.org \
		--allowed-cors-origins https://api.example.org \
		http://127.0.0.1:3000 \
		https://ory.example.org

### Redirects

Per default all default redirects will go to to `+"`"+`[publish-url]`+"`"+`. You can change this behavior using
the `+"`"+`--default-redirect-url`+"`"+` flag:

    $ %[1]s --project <your-project-slug> \
		--default-redirect-url /welcome \
		http://127.0.0.1:3000 \
		https://ory.example.org

Now, all redirects happening e.g. after login will point to `+"`"+`/welcome`+"`"+` instead of `+"`"+`/`+"`"+` unless you
have specified custom redirects in your Ory configuration or in the flow's `+"`"+`?return_to=`+"`"+` query parameter.

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
when calling the proxy - for example: `+"`"+`http://127.0.0.1:4000/.ory/jwks.json`+"`"+`

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

			redirectURL, err := url.ParseRequestURI(stringsx.Coalesce(flagx.MustGetString(cmd, DefaultRedirectURLFlag), selfURLString))
			if err != nil {
				return err
			}

			projectSlugId := getProjectSlugId(cmd)

			oryURL, err := getEndpointURL(cmd, projectSlugId)
			if err != nil {
				return err
			}

			origins, err := corsx.NormalizeOriginStrings(append(
				flagx.MustGetStringSlice(cmd, CORSFlag), selfURL.String()),
			)
			if err != nil {
				return err
			}

			conf := &ProxyConfig{
				Port:              flagx.MustGetInt(cmd, PortFlag),
				NoJWT:             flagx.MustGetBool(cmd, WithoutJWTFlag),
				NoOpen:            !flagx.MustGetBool(cmd, OpenFlag),
				Upstream:          args[0],
				CookieDomain:      flagx.MustGetString(cmd, CookieDomainFlag),
				PublicURL:         selfURL,
				OryURL:            oryURL,
				PathPrefix:        "/.ory",
				DefaultRedirectTo: redirectURL,
				IsDev:             flagx.MustGetBool(cmd, DevFlag),
				IsDebug:           flagx.MustGetBool(cmd, DebugFlag),
				RewriteHost:       flagx.MustGetBool(cmd, RewriteHostFlag),
				CorsOrigins:       origins,
				ProjectSlugId:     projectSlugId,
			}

			return Run(cmd, conf, version, "cloud")
		},
	}

	proxyCmd.Flags().Bool(OpenFlag, false, "Open the browser when the proxy starts.")
	proxyCmd.Flags().String(CookieDomainFlag, "", "Set a dedicated cookie domain.")
	proxyCmd.Flags().StringP(ProjectFlag, ProjectFlag[:0], "", "The slug of your Ory Network project.")
	proxyCmd.Flags().Int(PortFlag, portFromEnv(), "The port the proxy should listen on.")
	proxyCmd.Flags().Bool(WithoutJWTFlag, false, "Do not create a JWT from the Ory Session. Useful if you need fast start up times of the Ory Proxy.")
	proxyCmd.Flags().String(DefaultRedirectURLFlag, "", "Set the URL to redirect to per default after e.g. login or account creation.")
	proxyCmd.Flags().StringSlice(CORSFlag, []string{}, "A list of allowed CORS origins. Wildcards are allowed.")
	proxyCmd.Flags().Bool(DevFlag, false, "Use this flag when developing locally.")
	proxyCmd.Flags().Bool(DebugFlag, false, "Use this flag to debug, for example, CORS requests.")
	proxyCmd.Flags().Bool(RewriteHostFlag, false, "Use this flag to rewrite the host header to the upstream host.")

	client.RegisterConfigFlag(proxyCmd.PersistentFlags())
	client.RegisterYesFlag(proxyCmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(proxyCmd.PersistentFlags())

	return proxyCmd
}

const (
	envVarSlug   = "ORY_PROJECT_SLUG"
	envVarSDK    = "ORY_SDK_URL"
	envVarKratos = "ORY_KRATOS_URL"
)

func getProjectSlugId(cmd *cobra.Command) string {
	if slug := os.Getenv(envVarSlug); len(slug) > 0 {
		return slug
	} else if slug := flagx.MustGetString(cmd, ProjectFlag); len(slug) > 0 {
		return slug
	}
	return ""
}

func getEndpointURL(cmd *cobra.Command, slug string) (*url.URL, error) {
	var target string
	if slug != "" {
		target = fmt.Sprintf("https://%s.projects.oryapis.com/", slug)
	} else if url := stringsx.Coalesce(os.Getenv(envVarSDK), os.Getenv(envVarKratos)); len(url) > 0 {
		target = url
	}

	if len(target) == 0 {
		return nil, errors.Errorf("Please provide your project slug or project id using the --%s flag or the %s environment variable.", ProjectFlag, envVarSlug)
	}

	upstream, err := url.ParseRequestURI(target)
	if err != nil {
		return nil, errors.Errorf("Unable to parse \"%s\" as an URL: %s", target, err)
	}

	if err = printDeprecations(cmd, target); err != nil {
		return nil, errors.WithStack(err)
	}

	return upstream, nil
}

func printDeprecations(cmd *cobra.Command, target string) error {
	if deprecated := stringsx.Coalesce(os.Getenv(envVarSDK), os.Getenv(envVarKratos)); len(deprecated) > 0 {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "It is recommended to use the --%s flag or the %s environment variable for better developer experience. Environment variables %s and %s will continue to work!\n", ProjectFlag, envVarSlug, envVarSDK, envVarKratos)
	}

	found := map[string]string{}
	for k, s := range map[string]string{
		envVarSlug:         os.Getenv(envVarSlug),
		envVarSDK:          os.Getenv(envVarSDK),
		envVarKratos:       os.Getenv(envVarKratos),
		"--" + ProjectFlag: flagx.MustGetString(cmd, ProjectFlag),
	} {
		if len(s) > 0 {
			found[k] = s
		}
	}

	if len(found) > 1 {
		var values []string
		for k, v := range found {
			values = append(values, fmt.Sprintf("%s=%s", k, v))
		}
		sort.Strings(values)

		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Attention! We found multiple sources for the project slug. Please clean up environment variables and flags to ensure that the correct value is being used. Found values:\n\n\t%s\n\n", strings.Join(values, "\n\t"))
		return errors.New("Could not decide which source to use for endpoint configuration")
	}

	return nil
}
