package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/cli/x"
	"github.com/ory/graceful"
	"github.com/ory/herodot"
	"github.com/ory/x/flagx"
	"github.com/ory/x/logrusx"
)

const (
	PortFlag          = "port"
	ProtectPathsFlag  = "protect-path-prefix"
	NoCertInstallFlag = "dont-install-cert"
	NoOpenFlag        = "no-open"
)

func NewProxyLocalCmd() *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "local [upstream]",
		Short: "Develop an application locally and integrate it with Ory",
		Args:  cobra.ExactArgs(1),
		Long: fmt.Sprintf(`This command starts a reverse proxy which can be deployed in front of your application. This works best on local (your computer) environments, for example when developing a React, NodeJS, Java, PHP app.

To require login before accessing paths in your application, use the --%[1]s flag:

	$ ory local proxy --port 4000 http://localhost:3000 --%[1]s /members --%[1]s /admin

The --%[1]s flag is currently using a string prefix match. Future versions will
include support for regular expressions and glob matching.

If the request is authenticated, a JSON Web Token will be sent in the HTTP Authorization Header containing the
Ory Session:

	GET / HTTP/1.1
	Host: www.example.com
	Authorization Bearer <the-json-web-token>

The JSON Web Token claims contain:

* The "sub" field which is set to the Ory Identity ID.
* The "session" field which contains the full Ory Session.

The JSON Web Token is signed using the ES256 algorithm. The public key can be found by fetching the /.ory/jwks.json path
when calling the proxy - for example http://127.0.0.1:4000/.ory/jwks.json

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

`, ProtectPathsFlag),
		/*
		   The --%s values support regular expression templating, meaning that you can use regular expressions within "<>":

		   	$ ory proxy http://localhost:3000 --allow --%s "http://localhost:3000/<(login|dashboard)>" --%s "http://localhost:3000/<([0-9]{3})>"

		   The supported Regular Expression Syntax is RE2 and documented at: https://golang.org/pkg/regexp/
		   To test your Regular Expression, head over to https://regex101.com and select "Golang" on the left.
		*/
		RunE: func(cmd *cobra.Command, args []string) error {
			upstream, err := url.ParseRequestURI(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to parse upstream URL")
			}

			c, cleanup, err := createTLSCertificate(cmd)
			if err != nil {
				return err
			}

			l := logrusx.New("ory/proxy", x.BuildVersion)

			handler := httputil.NewSingleHostReverseProxy(upstream)
			writer := herodot.NewJSONWriter(l)

			mw := negroni.New()
			// mw.Use(reqlog.NewMiddlewareFromLogger(l, "ory"))

			signer, key, err := newSigner(l)
			if err != nil {
				return errors.WithStack(err)
			}

			endpoint, err := getEndpointURL(cmd)
			if err != nil {
				return errors.WithStack(err)
			}

			mw.UseFunc(func(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {
				// Disable HSTS because it is very annoying to use in localhost.
				w.Header().Set("Strict-Transport-Security", "max-age=0;")
				n(w, r)
			})

			mw.UseFunc(checkOry(cmd, writer, l, key, signer, endpoint)) // This must be the last method before the handler
			mw.UseHandler(handler)

			addr := fmt.Sprintf(":%d", flagx.MustGetInt(cmd, PortFlag))
			server := graceful.WithDefaults(&http.Server{
				Addr:      addr,
				Handler:   mw,
				TLSConfig: &tls.Config{Certificates: []tls.Certificate{*c}},
			})

			l.Printf("Starting the https reverse proxy on: %s", server.Addr)
			proxyUrl := fmt.Sprintf("https://localhost:%d/", flagx.MustGetInt(cmd, PortFlag))
			l.Printf(`To access your application through the Ory Proxy, open:

	%s`, proxyUrl)
			if !flagx.MustGetBool(cmd, NoOpenFlag) {
				if err := exec.Command("open", proxyUrl).Run(); err != nil {
					l.WithError(err).Warn("Unable to automatically open the proxy URL in your browser. Please open it manually!")
				}
			}

			if err := graceful.Graceful(func() error {
				return server.ListenAndServeTLS("", "")
			}, func(ctx context.Context) error {
				l.Println("http reverse proxy was shutdown gracefully")
				if err := server.Shutdown(ctx); err != nil {
					return err
				}

				return cleanup()
			}); err != nil {
				l.Fatalln("Failed to gracefully shutdown https reverse proxy")
			}

			return nil
		},
	}

	var port int64 = 4000
	if p, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64); p != 0 {
		port = p
	}

	proxyCmd.Flags().Int(PortFlag, int(port), "The port the proxy should listen on.")
	proxyCmd.Flags().Bool(NoCertInstallFlag, false, "If set will not try to add the HTTPS certificate to your certificate store.")
	proxyCmd.Flags().StringSlice(ProtectPathsFlag, []string{}, "Require authentication before accessing these paths.")
	proxyCmd.Flags().Bool(NoOpenFlag, false, "Do not open the browser when the proxy starts.")
	remote.RegisterClientFlags(proxyCmd.PersistentFlags())
	return proxyCmd
}
