package remote_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/herodot"
	"github.com/ory/x/logrusx"
)

const (
	kratosAdminPath    = "/api/kratos/admin"
	backofficeSlugPath = "/backoffice/token/slug"
	slug               = "pedantic-shannon-6947p3gdsf"
	TokenKey           = "ORY_ACCESS_TOKEN"
	TokenValue         = "nCCXCGpG6S6ejFEHfbuZvpaW9Ts84Pkq"
)

var (
	slugJSON = json.RawMessage(`{"slug":"` + slug + `"}`)
)

type Output struct {
	slug      string
	errorMSG  string
	kratosURL *url.URL
}

type TestingStruct struct {
	description string
	input       []string
	output      Output
}

func fakeProjectEndpoint(t *testing.T, writer herodot.Writer) *url.URL {
	router := httprouter.New()
	router.GET(kratosAdminPath+"/identities", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	router.GET(backofficeSlugPath, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		writer.Write(w, r, slugJSON)
	})
	api := httptest.NewServer(router)
	t.Cleanup(api.Close)
	parsed, err := url.ParseRequestURI(api.URL)
	require.NoError(t, err)
	return parsed
}

func TestClient(t *testing.T) {
	os.Setenv(TokenKey, TokenValue)
	l := logrusx.New("ory cli", "tests")
	writer := herodot.NewJSONWriter(l)

	kratosApi := fakeProjectEndpoint(t, writer)
	slugApi := fakeSlugEndpoint(t, writer)
	rsp, err := http.Get(slugApi.String() + backofficeSlugPath)
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(bodyBytes), slug)
	assert.NotEmpty(t, rsp)
	assert.NoError(t, err)

	FlagTest := []TestingStruct{
		{
			description: "Missing schema",
			input:       []string{"oryapis:8080", "api.console.ory:8080"},
			output: Output{
				slug:      "",
				errorMSG:  "Could not parse requested url: api.console.ory:8080",
				kratosURL: nil,
			},
		},
		{
			description: "Invalid urls",
			input:       []string{"hi/there?", "hi/there?"},
			output: Output{
				slug:      "",
				errorMSG:  "invalid URI for request",
				kratosURL: nil,
			},
		},
		{
			description: "Valid urls",
			input:       []string{kratosApi.String(), slugApi.String()},
			output: Output{
				slug:     slug,
				errorMSG: "",
				kratosURL: &url.URL{
					Scheme: kratosApi.Scheme,
					Host:   fmt.Sprintf("%s.projects.%s", slug, kratosApi.Host),
					Path:   kratosAdminPath,
				},
			},
		},
	}

	t.Run("function=GetProjectSlug", func(t *testing.T) {
		for _, tt := range FlagTest {
			t.Run(tt.description, func(t *testing.T) {
				s, err := remote.GetProjectSlug(tt.input[1])
				if err != nil {
					assert.Contains(t, err.Error(), tt.output.errorMSG)
				}
				assert.EqualValues(t, tt.output.slug, s)
			})
		}
	})

	t.Run("function=CreateKratosURL", func(t *testing.T) {
		for _, tt := range FlagTest {
			t.Run(tt.description, func(t *testing.T) {
				u, err := remote.CreateKratosURL(tt.input[0], tt.input[1])
				if err != nil {
					assert.Contains(t, err.Error(), tt.output.errorMSG)
				}
				assert.EqualValues(t, tt.output.kratosURL, u)
			})
		}
	})

}
