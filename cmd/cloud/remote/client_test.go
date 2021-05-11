package remote_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
)

func newCommand(t *testing.T, ctx context.Context) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{New: cmd.NewRootCmd, Ctx: ctx}
}

func TestIdentityList(t *testing.T) {
	t.Run("should fail when no env is provided", func(t *testing.T) {
		ctx := context.Background()
		stdout, stderr, err := newCommand(t, ctx).Exec(os.Stdin, "identities", "list")
		t.Logf("stdout:\n%s", stdout)
		t.Logf("stderr:\n%s", stderr)
		require.NoError(t, err)
		assert.Contains(t, stderr, "Ory API Token could not be detected! Did you forget to set the environment variable")
	})
	t.Run("should accept a valid token", func(t *testing.T) {
		ctx := context.Background()
		os.Setenv(TokenKey, TokenValue)
		stdout, stderr, err := newCommand(t, ctx).Exec(os.Stdin, "identities", "list")
		require.NoError(t, err)
		t.Logf("stdout:\n%s", stdout)
		t.Logf("stderr:\n%s", stderr)
		require.NoError(t, err)
		assert.NotContains(t, stderr, "Ory API Token could not be detected! Did you forget to set the environment variable")
	})
	t.Run("should fail when APIs are not accessible", func(t *testing.T) {

	})
	t.Run("should return expected values for fake APIs", func(t *testing.T) {
		//ctx = context.WithValue(ctx, remote.FlagAPIEndpoint, APIEndpoint)
		//ctx = context.WithValue(ctx, remote.FlagConsoleURL, ConsoleURL)
	})
}
