package remote_test

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/herodot"
	"github.com/ory/x/logrusx"
)

var (
	slug           = "pedantic-shannon-6947p3gdsf"
	slugJSON       = json.RawMessage(`{"slug":"` + slug + `"}`)
	insecureClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
)

func fakeProjectEndpoint(t *testing.T, writer herodot.Writer) *url.URL {
	router := httprouter.New()
	router.GET("/api/kratos/admin/identities", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		writer.Write(w, r, []byte("[]"))
	})
	api := httptest.NewServer(router)
	t.Cleanup(api.Close)
	parsed, err := url.ParseRequestURI(api.URL)
	require.NoError(t, err)
	return parsed
}

func fakeSlugEndpoint(t *testing.T, writer herodot.Writer) *url.URL {
	router := httprouter.New()
	router.GET("/backoffice/token/slug", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		writer.Write(w, r, slugJSON)
	})
	api := httptest.NewServer(router)
	t.Cleanup(api.Close)
	parsed, err := url.ParseRequestURI(api.URL)
	require.NoError(t, err)
	return parsed
}

func TestClient(t *testing.T) {
	l := logrusx.New("ory cli", "tests")
	writer := herodot.NewJSONWriter(l)

	kratosApi := fakeProjectEndpoint(t, writer)
	slugApi := fakeSlugEndpoint(t, writer)
	rsp, err := http.Get(slugApi.String() + "/backoffice/token/slug")
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(bodyBytes), slug)
	assert.NotEmpty(t, rsp)
	assert.NoError(t, err)

	var FlagTest = []struct {
		description string
		in          []string
		out         []string
	}{
		{"Missing schema", []string{"oryapis:8080", "api.console.ory:8080"}, []string{"", "Could not parse requested url: api.console.ory:8080"}},
		{"Invalid urls", []string{"hi/there?", "hi/there?"}, []string{"", "invalid URI for request"}},
		{"Valid urls", []string{kratosApi.String(), slugApi.String()}, []string{slug, ""}},
	}

	t.Run("function=GetProjectSlug", func(t *testing.T) {
		for _, tt := range FlagTest {
			t.Run(tt.description, func(t *testing.T) {
				s, err := remote.GetProjectSlug(tt.in[1])
				if err != nil {
					assert.Contains(t, err.Error(), tt.out[1])
				}
				assert.EqualValues(t, tt.out[0], s)
			})
		}
	})

}
