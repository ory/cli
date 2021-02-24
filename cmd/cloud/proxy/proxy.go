package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/elnormous/contenttype"
	"github.com/gobuffalo/x"
	"github.com/gofrs/uuid/v3"
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/graceful"
	"github.com/ory/herodot"
	"github.com/ory/x/flagx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/jwksx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/reqlog"
	"github.com/ory/x/urlx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
	"github.com/tidwall/gjson"
	"github.com/urfave/negroni"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	PortFlag           = "port"
	AllowAnonymousFlag = "allow-anonymous"
)

func NewProxyCmd() *cobra.Command {
	proxyCmd := &cobra.Command{
		Use:   "proxy [upstream]",
		Short: "Secure Endpoint Using the Ory Reverse Proxy",
		Long: fmt.Sprintf(`This command starts a reverse proxy which can be deployed in front of your application.

All incoming requests will be checked for a valid session, and if a valid session exists, the request will be forwarded
to your application. To exclude paths from authentication - useful for public pages such as your landing page - use
the --%s flag:

	$ ory proxy -port 4000 http://localhost:3000 --%s /login --%s /dashboard

The --%s flag is currently using a string exact match - except that the host is case insensitive. Future versions will
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

`, AllowAnonymousFlag, AllowAnonymousFlag, AllowAnonymousFlag, AllowAnonymousFlag),
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

			if flagx.MustGetString(cmd, remote.FlagProject) == "" {
				return errors.New("flag --project must be set")
			}

			l := logrusx.New("ory/proxy", x.Version)

			handler := httputil.NewSingleHostReverseProxy(upstream)
			writer := herodot.NewJSONWriter(l)

			mw := negroni.New()
			mw.Use(reqlog.NewMiddlewareFromLogger(l, "ory"))

			signer, key, err := newSigner(l)
			if err != nil {
				return err
			}

			endpoint, err := getEndpointURL(cmd)
			if err != nil {
				return err
			}
			mw.UseFunc(checkOry(cmd, writer, key, signer, endpoint))
			mw.UseHandler(handler)

			addr := fmt.Sprintf(":%d", flagx.MustGetInt(cmd, PortFlag))
			server := graceful.WithDefaults(&http.Server{
				Addr:    addr,
				Handler: mw,
			})

			l.Printf("Starting the http reverse proxy on: %s", server.Addr)
			if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
				l.Fatalln("Failed to gracefully shutdown http reverse proxy")
			}

			l.Println("Http reverse proxy was shutdown gracefully")
			return nil
		},
	}

	var port int64 = 4000
	if p, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64); p != 0 {
		port = p
	}

	proxyCmd.Flags().Int(PortFlag, int(port), "The port the proxy should listen on.")
	proxyCmd.Flags().StringSliceP(AllowAnonymousFlag, AllowAnonymousFlag[:1], []string{}, "Allow one or more URLs to be accessed without authentication.")
	remote.RegisterClientFlags(proxyCmd.PersistentFlags())
	return proxyCmd
}

func newSigner(l *logrusx.Logger) (jose.Signer, *jose.JSONWebKeySet, error) {
	l.WithField("started_at", time.Now()).Info("")
	key, err := jwksx.GenerateSigningKeys(
		uuid.Must(uuid.NewV4()).String(),
		"ES256",
		0,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to generate JSON Web Key")
	}
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: key.Keys[0].Key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create signer")
	}
	l.WithField("completed_at", time.Now()).Info("ES256 JSON Web Key generation completed.")
	return sig, key, nil
}

func checkOry(cmd *cobra.Command, writer herodot.Writer, keys *jose.JSONWebKeySet, sig jose.Signer, endpoint *url.URL) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	allowBypass := flagx.MustGetStringSlice(cmd, AllowAnonymousFlag)
	hc := httpx.NewResilientClientLatencyToleranceHigh(http.DefaultTransport)

	var publicKeys jose.JSONWebKeySet
	for _, key := range keys.Keys {
		publicKeys.Keys = append(publicKeys.Keys, key.Public())
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.URL.Path == "/.ory/jwks.json" {
			writer.Write(w, r, publicKeys)
			return
		}

		for _, allow := range allowBypass {
			if allow == urlFromRequest(r).Path {
				next(w, r)
				return
			}
		}

		accepted, _, err := contenttype.GetAcceptableMediaType(r, []contenttype.MediaType{
			contenttype.NewMediaType("text/html"), // default offer
			contenttype.NewMediaType("application/json"),
		})
		if err != nil {
			writer.WriteError(w, r, err)
			return
		}

		isJsonRequest := accepted.Type+"/"+accepted.Subtype == "application/json"

		target, err := getEndpointURL(cmd)
		if err != nil {
			writer.WriteError(w, r, errors.WithStack(err))
			return
		}

		session, err := checkSession(cmd, hc, r, target)
		if err != nil || !gjson.GetBytes(session, "active").Bool(){
			if isJsonRequest {
				innerErr := herodot.ErrUnauthorized.WithReasonf("The provided credentials are expired, malformed, missing, or otherwise invalid.")
				if err != nil {
					innerErr.Wrap(err)
				}

				writer.WriteError(w, r, errors.WithStack(innerErr))
				return
			}
			http.Redirect(w, r, urlx.AppendPaths(endpoint, "api", "kratos", "public", "self-service", "login", "browser").String(), http.StatusFound)
			return
		}

		if !gjson.GetBytes(session, "active").Bool() {
			if isJsonRequest {
				writer.WriteError(w, r, errors.WithStack(herodot.ErrUnauthorized.WithReasonf("The provided credentials are expired, malformed, missing, or otherwise invalid.")))
				return
			}
			http.Redirect(w, r, urlx.AppendPaths(endpoint, "api", "kratos", "public", "self-service", "login", "browser").String(), http.StatusFound)
			return
		}

		now := time.Now().UTC()
		raw, err := jwt.Signed(sig).Claims(&
			jwt.Claims{
				Issuer:    target.String(),
				Subject:   gjson.GetBytes(session, "identity.id").String(),
				Expiry:    jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
				ID:        uuid.Must(uuid.NewV4()).String(),
			}).Claims(map[string]interface{}{"session": session}).CompactSerialize()
		if err != nil {
			writer.WriteError(w, r, err)
			return
		}

		r.Header.Set("Authorization", "Bearer "+raw)
		r.Header.Del("Cookie")
		next(w, r)
	}
}

func checkSession(cmd *cobra.Command, c *http.Client, r *http.Request, target *url.URL) (json.RawMessage, error) {
	target.Path = filepath.Join(target.Path, "api", "kratos", "public", "sessions", "whoami")
	req, err := http.NewRequest("GET", target.String(), nil)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError)
	}

	req.Header.Set("Cookie", r.Header.Get("Cookie"))
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("X-Session-Token", r.Header.Get("X-Session-Token"))
	req.Header.Set("X-Request-Id", r.Header.Get("X-Request-Id"))
	req.Header.Set("Accept", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to call session checker: %s", err).WithWrap(err))
	}
	defer res.Body.Close()

	var body json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to decode session to JSON: %s", err).WithWrap(err))
	}

	return body, nil
}

func getEndpointURL(cmd *cobra.Command) (*url.URL, error) {
	if target, ok := cmd.Context().Value(remote.FlagEndpoint).(*url.URL); ok {
		return target, nil
	}

	upstream, err := url.ParseRequestURI(flagx.MustGetString(cmd, remote.FlagEndpoint))
	if err != nil {
		return nil, err
	}

	upstream.Host = fmt.Sprintf("%s.projects.%s", flagx.MustGetString(cmd, remote.FlagProject), upstream.Host)
	return upstream, nil
}

func urlFromRequest(r *http.Request) *url.URL {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return &url.URL{
		Scheme:   scheme,
		Host:     r.Host,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}
}
