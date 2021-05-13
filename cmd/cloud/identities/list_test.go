package identities_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/x/logrusx"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/ory/kratos-client-go"
	"github.com/ory/kratos/cmd/cliclient"

	"github.com/stretchr/testify/assert"

	"github.com/ory/cli/cmd"
	"github.com/ory/x/cmdx"
	"github.com/txn2/txeh"
)

const (
	TokenKey    = "ORY_ACCESS_TOKEN"
	TokenValue  = "nCCXCGpG6S6ejFEHfbuZvpaW9Ts84Pkq"
	APIEndpoint = "https://oryapis:8080"
	ConsoleURL  = "https://api.console.ory:8080"
	kratosAdminPath    = "/api/kratos/admin/identities"
	backofficeSlugPath = "/backoffice/token/slug"
	slug               = "pedantic-shannon-6947p3gdsf"
	//testHostfile = "/tmp/hostfile"
)

var (
	ctx = context.WithValue(context.Background(), cliclient.ClientContextKey, func(cmd *cobra.Command) *kratos.APIClient {
		return remote.NewAdminClient(APIEndpoint, ConsoleURL)
	})
	slugJSON = json.RawMessage(`{"slug":"` + slug + `"}`)
	identityJSON=json.RawMessage(`[{"id":"cbd285ee-b342-4384-bb32-74bba6d937d8","schema_id":"default","schema_url":"https://`+ slug + `.projects.oryapis:8080/api/kratos/public/schemas/default","traits":{"username":"qwerty"}}]`)
)

func newCommand(t *testing.T, ctx context.Context) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{New: cmd.NewRootCmd, Ctx: ctx}
}

func fakeProjectEndpoint(t *testing.T, writer herodot.Writer) *url.URL {
	router := httprouter.New()
	router.GET(kratosAdminPath, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		writer.Write(w, r, identityJSON)
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

func TestIdentityListNoToken(t *testing.T) {
	if os.Getenv("TEST_WILL_PANIC") == "1" {
		err := os.Unsetenv("TokenKey")
		require.NoError(t, err)
		newCommand(t, ctx).ExecExpectedErr(t, "identities", "list")
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestIdentityListNoToken")
	cmd.Env = append(os.Environ(), "TEST_WILL_PANIC=1")
	out, err := cmd.CombinedOutput()
	assert.NotNil(t, err)
	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	assert.Equal(t, true, ok)
	assert.Contains(t, string(out), "Ory API Token could not be detected! Did you forget to set the environment variable \"ORY_ACCESS_TOKEN\"?")
	assert.Equal(t, "exit status 1", e.Error())
}

func TestIdentityListWithToken(t *testing.T) {
	if os.Getenv("TEST_WILL_PANIC") == "1" {
		err := os.Setenv(TokenKey, TokenValue)
		require.NoError(t, err)
		newCommand(t, ctx).ExecExpectedErr(t, "identities", "list")
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestIdentityListWithToken")
	cmd.Env = append(os.Environ(), "TEST_WILL_PANIC=1", TokenKey+"="+TokenValue)
	out, err := cmd.CombinedOutput()
	assert.NotNil(t, err)
	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	assert.Equal(t, true, ok)
	assert.Contains(t, string(out), "context deadline exceeded (Client.Timeout exceeded while awaiting headers)")
	assert.Equal(t, "exit status 1", e.Error())
}

func TestIdentityListFakeAPI(t *testing.T) {
	l := logrusx.New("ory cli", "tests")
	writer := herodot.NewJSONWriter(l)
	kratosApi := fakeProjectEndpoint(t, writer)
	slugApi := fakeSlugEndpoint(t, writer)
	err := os.Setenv(TokenKey, TokenValue)
	require.NoError(t, err)

	hosts, err := txeh.NewHostsDefault()
	require.NoError(t, err)
	hosts.AddHost("127.0.0.1", fmt.Sprintf("%s.projects.127.0.0.1", slug))
	err = hosts.Save()
	require.NoError(t, err)
	defer func() {
		hosts.RemoveHost(fmt.Sprintf("%s.projects.127.0.0.1", slug))
		hosts.Reload()
	}()

	rsp, err := http.Get(kratosApi.String() + kratosAdminPath)
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(bodyBytes), slug)

	if os.Getenv("TEST_WILL_PANIC") == "1" {
		err = hosts.Reload()
		require.NoError(t, err)
		ctx = context.WithValue(context.Background(), cliclient.ClientContextKey, func(cmd *cobra.Command) *kratos.APIClient {
			return remote.NewAdminClient(kratosApi.String(), slugApi.String())
		})

		//newCommand(t, ctx).ExecNoErr(t, "identities", "list", "-f", "json")
		// Expect command to fail to capture output in the error msg
		newCommand(t, ctx).ExecExpectedErr(t, "identities", "list", "-f", "json")
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestIdentityListFakeAPI")
	cmd.Env = append(os.Environ(), "TEST_WILL_PANIC=1", TokenKey+"="+TokenValue)
	out, err := cmd.CombinedOutput()
	assert.NotNil(t, err)
	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	assert.Equal(t, true, ok)
	b, err := identityJSON.MarshalJSON()
	require.NoError(t, err)
	assert.Contains(t, string(out), string(b))
	assert.Equal(t, "exit status 1", e.Error())
}
