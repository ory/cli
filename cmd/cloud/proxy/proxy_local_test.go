package proxy_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/cli/cmd"
	"github.com/ory/cli/cmd/cloud/proxy"
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/herodot"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/logrusx"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	anonymous = "anonymous"
)

var (
	session = json.RawMessage(`{
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
	  "email": "foo@bar"
	}
  }
}`)

	insecureClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
)

func newCommand(t *testing.T, ctx context.Context) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{New: cmd.NewRootCmd, Ctx: ctx}
}

func newUpstream(t *testing.T, keyURL string, writer herodot.Writer) *httptest.Server {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) == 1 {
			_, _ = w.Write([]byte(anonymous))
			return
		}

		res, err := insecureClient.Get(keyURL)
		if err != nil {
			writer.WriteError(w, r, errors.WithStack(err))
			return
		}
		defer res.Body.Close()

		var keys jose.JSONWebKeySet
		require.NoError(t, json.NewDecoder(res.Body).Decode(&keys))

		tok, err := jwt.ParseSigned(parts[1])
		if err != nil {
			writer.WriteError(w, r, errors.WithStack(err))
			return
		}

		var claims map[string]interface{}
		if err := tok.Claims(keys.Keys[0].Key, &claims); err != nil {
			writer.WriteError(w, r, errors.WithStack(err))
			return
		}

		writer.Write(w, r, claims)
	}))

	t.Cleanup(upstream.Close)
	return upstream
}

func cloudAPi(t *testing.T, writer herodot.Writer) *url.URL {
	router := httprouter.New()
	router.GET("/api/kratos/public/sessions/whoami", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.Header.Get("Authorization") != "ory" {
			writer.WriteError(w, r, errors.WithStack(herodot.ErrUnauthorized))
			return
		}

		writer.Write(w, r, session)
	})

	router.GET("/api/kratos/public/self-service/login/browser", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		_, _ = w.Write([]byte("login"))
	})

	api := httptest.NewServer(router)
	t.Cleanup(api.Close)
	parsed, err := url.ParseRequestURI(api.URL)
	require.NoError(t, err)
	return parsed
}

func getRequest(t *testing.T, c *http.Client, href string) ([]byte, *http.Response) {
	res, err := c.Get(href)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	return body, res
}

func assertAnonymous(t *testing.T, c *http.Client, href string) {
	body, _ := getRequest(t, c, href)
	assert.EqualValues(t, string(body), anonymous, "%s", body)
}

func TestProxy(t *testing.T) {
	port, err := freeport.GetFreePort()
	require.NoError(t, err)
	proxyURL := fmt.Sprintf("https://127.0.0.1:%d", port)

	l := logrusx.New("ory cli", "tests")
	writer := herodot.NewJSONWriter(l)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	upstream := newUpstream(t, proxyURL+"/.ory/jwks.json", writer)
	cloudApi := cloudAPi(t, writer)
	ctx = context.WithValue(ctx, remote.FlagAPIEndpoint, cloudApi)

	go func() {
		stdout, stderr, err := newCommand(t, ctx).Exec(os.Stdin, "proxy", "local", upstream.URL, "--"+proxy.NoCertInstallFlag, "--"+proxy.NoOpenFlag, "--"+proxy.PortFlag, fmt.Sprintf("%d", port))
		assert.ErrorIs(t, err, context.Canceled)
		t.Logf("stdout:\n%s", stdout)
		t.Logf("stderr:\n%s", stderr)
	}()

	var tries int
	for {
		time.Sleep(time.Second)

		tries++
		if tries > 30 {
			t.Fatal("Proxy did not come alive")
		}

		res, err := insecureClient.Get(proxyURL + "/.ory/jwks.json")
		if err != nil {
			t.Logf("Proxy is not yet alive: %s", err)
			continue
		}
		_ = res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Logf("Proxy is not yet alive: %d is not 200", res.StatusCode)
			continue
		}

		break
	}

	t.Run("allows anonymous paths", func(t *testing.T) {
		assertAnonymous(t, insecureClient, proxyURL+"/public/1")
		assertAnonymous(t, insecureClient, proxyURL+"/public/2")
	})

	t.Run("responds with 302 redirect to login if HTML request", func(t *testing.T) {
		for k, accept := range []string{
			"text/html",
			"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"application/xhtml+xml",
		} {
			t.Run("case="+strconv.Itoa(k), func(t *testing.T) {
				req, _ := http.NewRequest("GET", proxyURL+"/.ory/api/kratos/public/self-service/login/browser?return_to="+proxyURL+"/private/1", nil)
				if accept != "" {
					req.Header.Set("Accept", accept)
				}

				res, err := insecureClient.Do(req)
				require.NoError(t, err)
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)

				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.EqualValues(t, proxyURL+"/.ory/api/kratos/public/self-service/login/browser?return_to="+proxyURL+"/private/1", res.Request.URL.String())
				assert.EqualValues(t, "login", string(body))
			})
		}
	})

	t.Run("responds forwards hosted pages", func(t *testing.T) {
		for k, accept := range []string{
			"text/html",
			"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"application/xhtml+xml",
		} {
			t.Run("case="+strconv.Itoa(k), func(t *testing.T) {
				req, _ := http.NewRequest("GET", proxyURL+"/.ory/api/kratos/public/self-service/login/browser?return_to="+proxyURL+"/private/2", nil)
				if accept != "" {
					req.Header.Set("Accept", accept)
				}

				res, err := insecureClient.Do(req)
				require.NoError(t, err)
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)

				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.EqualValues(t, proxyURL+"/.ory/api/kratos/public/self-service/login/browser?return_to="+proxyURL+"/private/2", res.Request.URL.String())
				assert.EqualValues(t, "login", string(body))
			})
		}
	})
}
