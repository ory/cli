// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/pflag"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/gofrs/uuid/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/tidwall/gjson"
	"github.com/urfave/negroni"

	"github.com/ory/graceful"
	"github.com/ory/herodot"
	"github.com/ory/x/corsx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/jwksx"
	"github.com/ory/x/proxy"
	"github.com/ory/x/urlx"
)

const (
	PortFlag                  = "port"
	OpenFlag                  = "open"
	DevFlag                   = "dev"
	DebugFlag                 = "debug"
	WithoutJWTFlag            = "no-jwt"
	CookieDomainFlag          = "cookie-domain"
	DefaultRedirectURLFlag    = "default-redirect-url"
	CORSFlag                  = "allowed-cors-origins"
	AdditionalCORSHeadersFlag = "additional-cors-headers"
	RewriteHostFlag           = "rewrite-host"
)

type config struct {
	port                  int
	open                  bool
	noJWT                 bool
	upstream              string
	cookieDomain          string
	publicURL             *url.URL
	pathPrefix            string
	defaultRedirectTo     cmdx.URL
	isTunnel              bool
	isDebug               bool
	isDev                 bool
	corsOrigins           []string
	additionalCorsHeaders []string

	// rewriteHost means the host header will be rewritten to the upstream host.
	// This is useful in cases where upstream resolves requests based on Host.
	rewriteHost bool
}

func registerProxyConfigFlags(conf *config, flags *pflag.FlagSet) {
	flags.BoolVar(&conf.open, OpenFlag, false, "Open the browser when the proxy starts.")
	flags.BoolVar(&conf.noJWT, WithoutJWTFlag, false, "Do not create a JWT from the Ory Session. Useful if you need fast start up times of the Ory Proxy.")
}

func registerConfigFlags(conf *config, flags *pflag.FlagSet) {
	flags.StringVar(&conf.cookieDomain, CookieDomainFlag, "", "Set a dedicated cookie domain.")
	flags.IntVar(&conf.port, PortFlag, portFromEnv(), "The port the proxy should listen on.")
	flags.Var(&conf.defaultRedirectTo, DefaultRedirectURLFlag, "Set the URL to redirect to per default after e.g. login or account creation.")
	flags.StringSliceVar(&conf.corsOrigins, CORSFlag, []string{}, "A list of allowed CORS origins. Wildcards are allowed.")
	flags.StringSliceVar(&conf.additionalCorsHeaders, AdditionalCORSHeadersFlag, []string{}, "A list of additional CORS headers to allow. Wildcards are allowed.")
	flags.BoolVar(&conf.isDev, DevFlag, true, "This flag is deprecated as the command is only supposed to be used during development.")
	flags.BoolVar(&conf.isDebug, DebugFlag, false, "Use this flag to debug, for example, CORS requests.")
	flags.BoolVar(&conf.rewriteHost, RewriteHostFlag, false, "Use this flag to rewrite the host header to the upstream host.")
}

func portFromEnv() int {
	port := 4000
	if p, err := strconv.ParseInt(os.Getenv("PORT"), 10, 0); err == nil {
		port = int(p)
	}
	return port
}

type responseStatusCatcher struct {
	http.ResponseWriter
	status int
}

func (r *responseStatusCatcher) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func runReverseProxy(ctx context.Context, h *client.CommandHelper, stdErr io.Writer, conf *config, name string) error {
	signer, key, err := newJWTSigner()
	if err != nil {
		return err
	}

	apiKey, removeAPIKey, err := h.TemporaryAPIKey(ctx, fmt.Sprintf("Ory %s temporary API key - %s", name, h.UserName(ctx)))
	if err != nil {
		return err
	}
	defer func() {
		if err := removeAPIKey(); err != nil {
			_, _ = fmt.Fprintf(stdErr, "unable to remove temporary API key, please remove it manually: %s\n", err)
		}
	}()

	project, err := h.GetSelectedProject(ctx)
	if err != nil {
		return err
	}
	oryURL := client.CloudAPIsURL(project.Slug + ".projects")
	oryURL.Host = strings.TrimSuffix(oryURL.Host, ":443")

	writer := herodot.NewJSONWriter(&errorLogger{Writer: stdErr})
	mw := negroni.New()

	mw.UseFunc(func(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {
		start := time.Now()
		_, _ = fmt.Fprintf(stdErr, "%s [%s]\n", r.Method, r.URL)
		statusCatcher := &responseStatusCatcher{ResponseWriter: w}
		n(statusCatcher, r)
		_, _ = fmt.Fprintf(stdErr, "=> %d %s [%s] took %s\n", statusCatcher.status, r.Method, r.URL, time.Since(start))
	})

	mw.UseFunc(func(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {
		// Disable HSTS because it is very annoying to use on localhost.
		w.Header().Set("Strict-Transport-Security", "max-age=0;")
		n(w, r)
	})

	if !conf.noJWT {
		mw.UseFunc(sessionToJWTMiddleware(conf, writer, key, signer, oryURL)) // This must be the last method before the handler
	}

	var upstream *url.URL
	if conf.isTunnel {
		upstream = oryURL
	} else {
		upstream, err = url.ParseRequestURI(conf.upstream)
		if err != nil {
			return errors.Wrap(err, "unable to parse upstream URL")
		}
	}

	mw.UseHandler(proxy.New(
		func(ctx context.Context, r *http.Request) (context.Context, *proxy.HostConfig, error) {
			if conf.isTunnel || strings.HasPrefix(r.URL.Path, conf.pathPrefix) {
				return ctx, &proxy.HostConfig{
					CookieDomain:   conf.cookieDomain,
					UpstreamHost:   oryURL.Host,
					UpstreamScheme: oryURL.Scheme,
					TargetHost:     oryURL.Host,
					PathPrefix:     conf.pathPrefix,
				}, nil
			}

			return ctx, &proxy.HostConfig{
				CookieDomain:   conf.cookieDomain,
				UpstreamHost:   upstream.Host,
				UpstreamScheme: upstream.Scheme,
				TargetHost:     upstream.Host,
				PathPrefix:     "",
			}, nil
		},
		proxy.WithReqMiddleware(func(r *httputil.ProxyRequest, c *proxy.HostConfig, body []byte) ([]byte, error) {
			if r.Out.URL.Host == oryURL.Host {
				r.Out.URL.Path = strings.TrimPrefix(r.Out.URL.Path, conf.pathPrefix)
				r.Out.Host = oryURL.Host
			} else if conf.rewriteHost {
				r.Out.Header.Set("X-Forwarded-Host", r.In.Host)
				r.Out.Host = c.UpstreamHost
			}

			publicURL := conf.publicURL
			if conf.pathPrefix != "" {
				publicURL = urlx.AppendPaths(publicURL, conf.pathPrefix)
			}

			r.Out.Header.Set("Ory-No-Custom-Domain-Redirect", "true")
			r.Out.Header.Set("Ory-Base-URL-Rewrite", publicURL.String())
			if len(apiKey) > 0 {
				r.Out.Header.Set("Ory-Base-URL-Rewrite-Token", apiKey)
			}

			return body, nil
		}),
		proxy.WithRespMiddleware(func(resp *http.Response, config *proxy.HostConfig, body []byte) ([]byte, error) {
			l, err := resp.Location()
			if err == nil {
				// Redirect to main page if path is the default ui welcome page.
				if l.Path == filepath.Join(conf.pathPrefix, "/ui/welcome") {
					resp.Header.Set("Location", conf.defaultRedirectTo.String())
				}
			}

			return body, nil
		}),
	))

	cleanup := func() error {
		return nil
	}

	var originFunc func(r *http.Request, origin string) bool
	if len(conf.corsOrigins) == 0 {
		originFunc = func(r *http.Request, origin string) bool {
			return true
		}
	}

	corsOrigins, err := corsx.NormalizeOriginStrings(append(conf.corsOrigins, conf.publicURL.String()))
	if err != nil {
		return err
	}
	addr := fmt.Sprintf(":%d", conf.port)
	ch := cors.New(cors.Options{
		AllowedOrigins:         corsOrigins,
		AllowOriginRequestFunc: originFunc,
		AllowedMethods:         corsx.CORSDefaultAllowedMethods,
		AllowedHeaders:         append(corsx.CORSRequestHeadersSafelist, append(corsx.CORSRequestHeadersExtended, conf.additionalCorsHeaders...)...),
		ExposedHeaders:         corsx.CORSResponseHeadersSafelist,
		MaxAge:                 0,
		AllowCredentials:       true,
		OptionsPassthrough:     false,
		Debug:                  conf.isDebug,
	})

	server := graceful.WithDefaults(&http.Server{
		Addr:    addr,
		Handler: ch.Handler(mw),
	})

	if conf.isTunnel {
		_, _ = fmt.Fprintf(stdErr, `To access Ory's APIs, use URL

	export ORY_SDK_URL=%[1]s # Linux / macOS
	set ORY_SDK_URL=%[1]s # Windows CMD
	$env:ORY_SDK_URL = "%[1]s" # Windows PowerShell

and configure your SDKs to point to it, for example in JavaScript:

	import { FrontendApi, Configuration } from '@ory/client-fetch'
	const ory = new FrontendApi(new Configuration({
	  basePath: '%[1]s',
	  baseOptions: {
		withCredentials: true
	  },
	  headers: {
	    Accept: "application/json"
      }
	}))

`, conf.publicURL.String())
	} else {
		_, _ = fmt.Fprintf(stdErr, `To access your application via the Ory Proxy, open:

	%s

`, conf.publicURL.String())
	}

	if conf.open {
		if err := browser.OpenURL(conf.publicURL.String()); err != nil {
			_, _ = fmt.Fprintln(stdErr, "Unable to automatically open the proxy URL in your browser. Please open it manually!")
		}
	}

	if err := graceful.Graceful(func() error {
		return server.ListenAndServe()
	}, func(ctx context.Context) error {
		_, _ = fmt.Fprintln(stdErr, "http server was shutdown gracefully")
		if err := server.Shutdown(ctx); err != nil {
			return err
		}

		return cleanup()
	}); err != nil {
		_, _ = fmt.Fprintf(stdErr, "Failed to gracefully shutdown http server: %s\n", err)
	}

	return nil
}

func newJWTSigner() (jose.Signer, *jose.JSONWebKeySet, error) {
	key, err := jwksx.GenerateSigningKeys(
		uuid.Must(uuid.NewV4()).String(),
		"ES256",
		0,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to generate JSON Web Key")
	}
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: key.Keys[0].Key}, (&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", key.Keys[0].KeyID))
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create signer")
	}
	return sig, key, nil
}

type errorLogger struct {
	io.Writer
}

func (e *errorLogger) ReportError(r *http.Request, _ int, err error, _ ...any) {
	_, _ = fmt.Fprintf(e.Writer, "encountered error on %s: %s\n", r.URL, err)
}

func sessionToJWTMiddleware(conf *config, writer herodot.Writer, keys *jose.JSONWebKeySet, sig jose.Signer, endpoint *url.URL) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	hc := httpx.NewResilientClient(httpx.ResilientClientWithMaxRetry(5), httpx.ResilientClientWithMaxRetryWait(time.Millisecond*5), httpx.ResilientClientWithConnectionTimeout(time.Second*30))

	var publicKeys jose.JSONWebKeySet
	for _, key := range keys.Keys {
		publicKeys.Keys = append(publicKeys.Keys, key.Public())
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		switch r.URL.Path {
		case path.Join(conf.pathPrefix, "/jwks.json"), path.Join(conf.pathPrefix, "/proxy/jwks.json"):
			writer.Write(w, r, publicKeys)
			return
		}

		session, err := checkSession(hc, r, endpoint)
		if err != nil || !gjson.GetBytes(session, "active").Bool() {
			next(w, r)
			return
		}

		if conf.noJWT || (len(conf.pathPrefix) > 0 && strings.HasPrefix(r.URL.Path, conf.pathPrefix)) {
			next(w, r)
			return
		}

		now := time.Now().UTC()
		raw, err := jwt.Signed(sig).Claims(&jwt.Claims{
			Issuer:    endpoint.String(),
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
		next(w, r)
	}
}

func checkSession(c *retryablehttp.Client, r *http.Request, target *url.URL) (json.RawMessage, error) {
	target = urlx.Copy(target)
	target.Path = path.Join(target.Path, "sessions", "whoami")
	req, err := retryablehttp.NewRequest("GET", target.String(), nil)
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
