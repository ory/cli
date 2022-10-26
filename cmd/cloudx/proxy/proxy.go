// Copyright © 2022 Ory Corp

package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"

	"github.com/gofrs/uuid/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
	"github.com/tidwall/gjson"
	"github.com/urfave/negroni"

	"github.com/ory/graceful"
	"github.com/ory/herodot"
	"github.com/ory/x/corsx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/jwksx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/proxy"
	"github.com/ory/x/urlx"
)

const (
	PortFlag               = "port"
	OpenFlag               = "open"
	DevFlag                = "dev"
	DebugFlag              = "debug"
	WithoutJWTFlag         = "no-jwt"
	CookieDomainFlag       = "cookie-domain"
	DefaultRedirectURLFlag = "default-redirect-url"
	ProjectFlag            = "project"
	CORSFlag               = "allowed-cors-origins"
	RewriteHostFlag        = "rewrite-host"
)

type config struct {
	port              int
	noOpen            bool
	noJWT             bool
	upstream          string
	cookieDomain      string
	publicURL         *url.URL
	oryURL            *url.URL
	pathPrefix        string
	defaultRedirectTo *url.URL
	isTunnel          bool
	isDebug           bool
	isDev             bool
	corsOrigins       []string

	// rewriteHost means the host header will be rewritten to the upstream host.
	// This is useful in cases where upstream resolves requests based on Host.
	rewriteHost bool
}

func portFromEnv() int {
	var port int64 = 4000
	if p, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64); p != 0 {
		port = p
	}
	return int(port)
}

var errNoApiKeyAvailable = errors.New("no api key available")

func noop() {}

func getAPIKey(conf *config, l *logrusx.Logger, h *client.CommandHelper) (apiKey string, cleanup func(), err error) {
	oryURLParts := strings.Split(conf.oryURL.Hostname(), ".")
	if len(oryURLParts) < 2 {
		l.Warnf("The Ory Network URL `%s` does not appear to be a a valid Ory Network URL. It should be in the format of `https://<project-slug>.projects.oryapis.com`. Skipping API key creation.", conf.oryURL)
		return "", noop, errNoApiKeyAvailable
	}

	if ak := client.GetProjectAPIKeyFromEnvironment(); len(ak) > 0 {
		return ak, noop, nil
	}

	if oryURLParts[0] == "playground" {
		l.Warnf("The Ory Proxy / Ory Tunnel does not support Social Sign In for the playground project.")
		return "", noop, errNoApiKeyAvailable
	}

	// For all other projects, except the playground, we should to authenticate.
	_, valid, err := h.HasValidContext()
	if errors.Is(err, client.ErrNoConfigQuiet) {
		l.Warn("Because you are not authenticated, the Ory CLI can not configure your project automatically. You can still use the Ory Proxy / Ory Tunnel, but complex flows such as Social Sign In will not work. Remove the `--quiet` flag or run `ory auth login` to authenticate.")
		return "", noop, errNoApiKeyAvailable
	} else if err != nil {
		return "", noop, err
	}

	if !valid {
		ok, err := cmdx.AskScannerForConfirmation("To support complex flows such as Social Sign In, the Ory CLI can configure your project automatically. To do so, you need to be signed in. Do you want to sign in?", h.Stdin, h.VerboseErrWriter)
		if err != nil {
			return "", noop, err
		}

		if !ok {
			l.Warn("Because you are not authenticated, the Ory CLI can not configure your project automatically. You can still use the Ory Proxy / Ory Tunnel, but complex flows such as Social Sign In will not work.")
			return "", noop, errNoApiKeyAvailable
		}

		if _, err := h.EnsureContext(); err != nil {
			return "", noop, err
		}
	}

	slug := oryURLParts[0]
	ak, err := h.CreateAPIKey(slug, "Ory CLI Proxy / Tunnel - Temporary API Key")
	if err != nil {
		l.WithError(err).Errorf("Unable to create API key. Do you have the required permissions to use the Ory CLI with project `%s`?", slug)
		return "", noop, errors.Wrapf(err, "unable to create API key for project %s", slug)
	}

	if !ak.HasValue() {
		return "", noop, errNoApiKeyAvailable
	}

	return *ak.Value, func() {
		if err := h.DeleteAPIKey(slug, ak.Id); err != nil {
			l.WithError(err).Errorf("Unable to clean up API Key automatically. Please remove it up manually in the Ory Console.")
		}
	}, nil
}

func run(cmd *cobra.Command, conf *config, version string, name string) error {
	h, err := client.NewCommandHelper(cmd)
	if err != nil {
		return err
	}

	upstream, err := url.ParseRequestURI(conf.upstream)
	if err != nil {
		return errors.Wrap(err, "unable to parse upstream URL")
	}

	l := logrusx.New("ory/"+strings.ToLower(name), version)
	writer := herodot.NewJSONWriter(l)
	mw := negroni.New()

	signer, key, err := newSigner(l, conf)
	if err != nil {
		return errors.WithStack(err)
	}

	apiKey, removeAPIKey, err := getAPIKey(conf, l, h)
	if errors.Is(err, errNoApiKeyAvailable) {
		// Do nothing - no API key is available and social sign in will not work.
	} else if err != nil {
		return err
	}
	defer removeAPIKey()

	mw.UseFunc(func(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {
		// Disable HSTS because it is very annoying to use in localhost.
		w.Header().Set("Strict-Transport-Security", "max-age=0;")
		n(w, r)
	})

	mw.UseFunc(checkOry(conf, l, writer, key, signer, conf.oryURL)) // This must be the last method before the handler

	mw.UseHandler(proxy.New(
		func(_ context.Context, r *http.Request) (*proxy.HostConfig, error) {
			if conf.isTunnel || strings.HasPrefix(r.URL.Path, conf.pathPrefix) {
				return &proxy.HostConfig{
					CookieDomain:   conf.cookieDomain,
					UpstreamHost:   conf.oryURL.Host,
					UpstreamScheme: conf.oryURL.Scheme,
					TargetHost:     conf.oryURL.Host,
					PathPrefix:     conf.pathPrefix,
				}, nil
			}

			return &proxy.HostConfig{
				CookieDomain:   conf.cookieDomain,
				UpstreamHost:   upstream.Host,
				UpstreamScheme: upstream.Scheme,
				TargetHost:     upstream.Host,
				PathPrefix:     "",
			}, nil
		},
		proxy.WithReqMiddleware(func(r *http.Request, c *proxy.HostConfig, body []byte) ([]byte, error) {
			if r.URL.Host == conf.oryURL.Host {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, conf.pathPrefix)
				r.Host = conf.oryURL.Host
			} else if conf.rewriteHost {
				r.Header.Set("X-Forwarded-Host", r.Host)
				r.Host = c.UpstreamHost
			}

			publicURL := conf.publicURL
			if conf.pathPrefix != "" {
				publicURL = urlx.AppendPaths(publicURL, conf.pathPrefix)
			}

			r.Header.Set("Ory-No-Custom-Domain-Redirect", "true")
			r.Header.Set("Ory-Base-URL-Rewrite", publicURL.String())
			if len(apiKey) > 0 {
				r.Header.Set("Ory-Base-URL-Rewrite-Token", apiKey)
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
	if conf.isDev {
		originFunc = func(r *http.Request, origin string) bool {
			return true
		}
	}

	proto := "http"
	addr := fmt.Sprintf(":%d", conf.port)
	ch := cors.New(cors.Options{
		AllowedOrigins:         conf.corsOrigins,
		AllowOriginRequestFunc: originFunc,
		AllowedMethods:         corsx.CORSDefaultAllowedMethods,
		AllowedHeaders:         append(corsx.CORSRequestHeadersSafelist, corsx.CORSRequestHeadersExtended...),
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
		_, _ = fmt.Fprintf(os.Stderr, `To access Ory's APIs, use URL

	%s

and configure your SDKs to point to it, for example in JavaScript:

	import { V0alpha2Api, Configuration } from '@ory/client'
	const ory = new V0alpha2Api(new Configuration({
	  basePath: 'http://localhost:4000',
	  baseOptions: {
		withCredentials: true
	  }
	}))

`, conf.publicURL.String())
	} else {
		_, _ = fmt.Fprintf(os.Stderr, `To access your application via the Ory Proxy, open:

	%s
`, conf.publicURL.String())
	}

	if !conf.noOpen {
		// #nosec G204 - this is ok
		if err := exec.Command("open", conf.publicURL.String()).Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to automatically open the proxy URL in your browser. Please open it manually!")
		}
	}

	if err := graceful.Graceful(func() error {
		return server.ListenAndServe()
	}, func(ctx context.Context) error {
		_, _ = fmt.Fprintf(os.Stderr, "http server was shutdown gracefully\n")
		if err := server.Shutdown(ctx); err != nil {
			return err
		}

		return cleanup()
	}); err != nil {
		l.Fatalf("Failed to gracefully shutdown %s server because: %s\n", proto, err)
	}

	return nil
}

func newSigner(l *logrusx.Logger, conf *config) (jose.Signer, *jose.JSONWebKeySet, error) {
	if conf.noJWT {
		return nil, &jose.JSONWebKeySet{}, nil
	}

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

func checkOry(conf *config, _ *logrusx.Logger, writer herodot.Writer, keys *jose.JSONWebKeySet, sig jose.Signer, endpoint *url.URL) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	hc := httpx.NewResilientClient(httpx.ResilientClientWithMaxRetry(5), httpx.ResilientClientWithMaxRetryWait(time.Millisecond*5), httpx.ResilientClientWithConnectionTimeout(time.Second*2))

	var publicKeys jose.JSONWebKeySet
	for _, key := range keys.Keys {
		publicKeys.Keys = append(publicKeys.Keys, key.Public())
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if !conf.noJWT && r.URL.Path == filepath.Join(conf.pathPrefix, "/proxy/jwks.json") {
			writer.Write(w, r, publicKeys)
			return
		}

		switch r.URL.Path {
		case filepath.Join(conf.pathPrefix, "/jwks.json"):
			writer.Write(w, r, publicKeys)
			return
		}

		session, err := checkSession(hc, r, endpoint)
		r.Header.Del("Authorization")
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
	target.Path = filepath.Join(target.Path, "api", "kratos", "public", "sessions", "whoami")
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
