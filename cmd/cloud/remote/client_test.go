package remote_test

import (
	"context"
	"encoding/json"
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
)

//.bin/cli identities list --api-endpoint oryapis:8080 --console-url console.ory:8080
//2021/05/11 12:36:28 [DEBUG] GET https://api.console.ory:8080/backoffice/token/slug
//2021/05/11 12:36:28 [DEBUG] GET https://unruffled-dijkstra-1qjqn90055.projects.oryapis:8080/api/kratos/admin/identities
//ID	VERIFIED ADDRESS 1	RECOVERY ADDRESS 1	SCHEMA ID	SCHEMA URL

const (
	TokenKey    = "ORY_ACCESS_TOKEN"
	TokenValue  = "nCCXCGpG6S6ejFEHfbuZvpaW9Ts84Pkq"
	APIEndpoint = "oryapis:8080"
	ConsoleURL  = "console.ory:8080"
)

var (
	slug = json.RawMessage(`{
  "slug": "unruffled-dijkstra-1qjqn90055"
}`)
	ctx = context.WithValue(context.Background(), cliclient.ClientContextKey, func(cmd *cobra.Command) *kratos.APIClient {
		return remote.NewAdminClient(cmd)
	})
)

func newCommand(t *testing.T, ctx context.Context) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{New: cmd.NewRootCmd, Ctx: ctx}
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
	cmd.Env = append(os.Environ(), "TEST_WILL_PANIC=1")
	out, err := cmd.CombinedOutput()
	assert.NotNil(t, err)
	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	assert.Equal(t, true, ok)
	assert.Contains(t, string(out), "Could not retrieve valid project slug from https://console.ory.sh")
	assert.Equal(t, "exit status 1", e.Error())
}
