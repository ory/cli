package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/elnormous/contenttype"
	"github.com/gofrs/uuid/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/herodot"
	"github.com/ory/x/flagx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/jwksx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/tlsx"
	"github.com/ory/x/urlx"
	"github.com/pkg/errors"
	"github.com/smallstep/truststore"
	"github.com/spf13/cobra"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

func checkOry(cmd *cobra.Command, writer herodot.Writer, l *logrusx.Logger, keys *jose.JSONWebKeySet, sig jose.Signer, endpoint *url.URL) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	protectPaths := flagx.MustGetStringSlice(cmd, ProtectPathsFlag)
	hc := httpx.NewResilientClient(httpx.ResilientClientWithMaxRetry(5), httpx.ResilientClientWithMaxRetryWait(time.Millisecond*5), httpx.ResilientClientWithConnectionTimeout(time.Second*2))

	var publicKeys jose.JSONWebKeySet
	for _, key := range keys.Keys {
		publicKeys.Keys = append(publicKeys.Keys, key.Public())
	}

	oryUpstream := httputil.NewSingleHostReverseProxy(endpoint)

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		originalHost := r.Host

		if r.URL.Path == "/.ory/jwks.json" {
			writer.Write(w, r, publicKeys)
			return
		}

		// We proxy ory things
		if strings.HasPrefix(r.URL.Path, "/.ory/api/kratos/public") {
			q := r.URL.Query()
			q.Set("alias", originalHost)

			r.URL.Path = strings.ReplaceAll(r.URL.Path, "/.ory/api/kratos/public", "/api/kratos/public")
			r.Host = endpoint.Host
			r.URL.RawQuery = q.Encode()

			l.WithRequest(r).
				WithField("forwarding_path", r.URL.String()).
				WithField("forwarding_host", r.Host).
				Debug("Forwarding request to Ory.")
			oryUpstream.ServeHTTP(w, r)
			return
		}

		var shouldProtect bool
		for _, protect := range protectPaths {
			if strings.HasPrefix(urlFromRequest(r).Path, protect) {
				shouldProtect = true
				break
			}
		}

		if !shouldProtect {
			next(w, r)
			return
		}

		var isJsonRequest bool
		accepted, _, err := contenttype.GetAcceptableMediaType(r, []contenttype.MediaType{
			contenttype.NewMediaType("text/html"), // default offer
			contenttype.NewMediaType("application/json"),
		})
		if err != nil {
			isJsonRequest = false
		} else {
			isJsonRequest = accepted.Type+"/"+accepted.Subtype == "application/json"
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		returnToLogin := func() {
			http.Redirect(w, r, fmt.Sprintf("/.ory/api/kratos/public/self-service/login/browser?return_to=%s://%s%s", scheme, r.Host, r.URL.Path), http.StatusFound)
		}

		session, err := checkSession(hc, r, endpoint)
		if err != nil || !gjson.GetBytes(session, "active").Bool() {
			if isJsonRequest {
				innerErr := herodot.ErrUnauthorized.WithReasonf("The provided credentials are expired, malformed, missing, or otherwise invalid.")
				if err != nil {
					innerErr.Wrap(err)
				}

				writer.WriteError(w, r, errors.WithStack(innerErr))
				return
			}
			returnToLogin()
			return
		}

		if !gjson.GetBytes(session, "active").Bool() {
			if isJsonRequest {
				writer.WriteError(w, r, errors.WithStack(herodot.ErrUnauthorized.WithReasonf("The provided credentials are expired, malformed, missing, or otherwise invalid.")))
				return
			}
			returnToLogin()
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
		r.Header.Del("Cookie")

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

func getEndpointURL(cmd *cobra.Command) (*url.URL, error) {
	if target, ok := cmd.Context().Value(remote.FlagAPIEndpoint).(*url.URL); ok {
		return target, nil
	}

	upstream, err := url.ParseRequestURI(flagx.MustGetString(cmd, remote.FlagAPIEndpoint))
	if err != nil {
		return nil, err
	}
	project, err := remote.GetProjectSlug(flagx.MustGetString(cmd, remote.FlagConsoleAPI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	upstream.Host = fmt.Sprintf("%s.projects.%s", project, upstream.Host)

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

func createTLSCertificate(cmd *cobra.Command) (*tls.Certificate, func() error, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)

	c, err := tlsx.CreateSelfSignedCertificate(key)
	if err != nil {
		return nil, nil, err
	}

	block, err := tlsx.PEMBlockForKey(key)
	if err != nil {
		return nil, nil, err
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})
	pemKey := pem.EncodeToMemory(block)
	cert, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		return nil, nil, err
	}

	const passwordMessage = "To modify your operating system certificate store, you might might be prompted for your password now:"

	if flagx.MustGetBool(cmd, NoCertInstallFlag) {
		return &cert, func() error {
			return nil
		}, nil
	}

	_, _ = fmt.Fprintln(os.Stdout, "Trying to install temporary TLS (HTTPS) certificate for localhost on your operating system. This allows to access the proxy using HTTPS.")
	_, _ = fmt.Fprintln(os.Stdout, passwordMessage)
	opts := []truststore.Option{
		truststore.WithFirefox(),
		truststore.WithJava(),
	}

	if err := truststore.Install(c, opts...); err != nil {
		return nil, nil, err
	}

	return &cert, func() error {
		_, _ = fmt.Fprintln(os.Stdout, passwordMessage)
		return truststore.Uninstall(c, opts...)
	}, nil
}
