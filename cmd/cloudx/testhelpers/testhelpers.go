// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/randx"
)

const testProjectPattern = "ory-cy-e2e-da2f162d-af61-42dd-90dc-e3fcfa7c84a0-"
const testAccountPrefix = "dev+orycye2eda2f162daf6142dd"

func TestName() string {
	return testProjectPattern + randx.MustString(16, randx.AlphaLowerNum)
}

func FakeEmail() string {
	return fmt.Sprintf(testAccountPrefix+".%s@ory.dev", randx.MustString(16, randx.AlphaLowerNum))
}

func FakePassword() string {
	return randx.MustString(16, randx.AlphaLowerNum)
}

func FakeName() string {
	return randx.MustString(1, randx.AlphaUpper) + randx.MustString(5, randx.AlphaLower)
}

func FakeAccount() (email string, password string, name string) {
	return FakeEmail(), FakePassword(), FakeName()
}

func NewConfigFile(t testing.TB) string {
	return filepath.Join(t.TempDir(), "config.json")
}

func ReadConfig(t testing.TB, configDir string) *client.Config {
	f, err := os.ReadFile(configDir)
	require.NoError(t, err)
	var ac client.Config
	require.NoError(t, json.Unmarshal(f, &ac))
	return &ac
}

var ErrAuthFlowTriggered = fmt.Errorf("flow triggered")

func WithEmitAuthFlowTriggeredErr(ctx context.Context, t testing.TB) context.Context {
	return client.ContextWithOptions(ctx,
		client.WithConfigLocation(NewConfigFile(t)),
		client.WithOpenBrowserHook(func(uri string) error {
			return fmt.Errorf("opened browser with %s: %w", uri, ErrAuthFlowTriggered)
		}),
	)
}

func WithCleanConfigFile(ctx context.Context, t testing.TB) context.Context {
	return client.ContextWithOptions(ctx, client.WithConfigLocation(NewConfigFile(t)))
}

func WithDuplicatedConfigFile(ctx context.Context, t testing.TB, originalFile string) context.Context {
	dst, err := os.Create(NewConfigFile(t))
	require.NoError(t, err)
	src, err := os.Open(originalFile)
	require.NoError(t, err)
	_, err = io.Copy(dst, src)
	require.NoError(t, err)

	return client.ContextWithOptions(ctx, client.WithConfigLocation(dst.Name()))
}

func Cmd(ctx context.Context) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{
		New: cmd.NewRootCmd,
		Ctx: client.ContextWithClient(ctx),
	}
}

func CreateProject(ctx context.Context, t testing.TB, workspace *string) *cloud.Project {
	args := []string{"create", "project", "--name", TestName(), "--format", "json"}
	if workspace != nil {
		args = append(args, "--workspace", *workspace)
	}
	stdout, stderr, err := Cmd(ctx).Exec(nil, args...)
	require.NoError(t, err, stderr)
	p := cloud.Project{}
	require.NoError(t, json.Unmarshal([]byte(stdout), &p), stdout)
	if ap, ok := p.AdditionalProperties["AdditionalProperties"]; ok {
		// the SDK types are weird sometimes...
		p.AdditionalProperties = ap.(map[string]interface{})
	}
	return &p
}

func CreateWorkspace(ctx context.Context, t testing.TB) string {
	return strings.TrimSpace(Cmd(ctx).ExecNoErr(t, "create", "workspace", "--name", TestName(), "--quiet"))
}

func SetDefaultProject(ctx context.Context, t testing.TB, projectID string) {
	require.Equal(t, projectID, strings.TrimSpace(Cmd(ctx).ExecNoErr(t, "use", "project", projectID, "--quiet")))
}

func GetDefaultProjectID(ctx context.Context, t testing.TB) string {
	return strings.TrimSpace(Cmd(ctx).ExecNoErr(t, "use", "project", "--quiet"))
}

func MakeRandomIdentity(t testing.TB, email string) string {
	path := filepath.Join(t.TempDir(), "import.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
  "schema_id": "preset://username",
  "traits": {
    "username": "`+email+`"
  }
}`), 0600))
	return path
}

func MakeRandomClient(t testing.TB, name string) string {
	homeDir, err := os.MkdirTemp(os.TempDir(), "cloudx-*")
	require.NoError(t, err)
	path := filepath.Join(homeDir, "import.json")
	require.NoError(t, os.WriteFile(path, []byte(`[
  {
    "client_name": "`+name+`"
  }
]`), 0600))
	return path
}

func ImportIdentity(ctx context.Context, t testing.TB, project string, stdin io.Reader) string {
	email := FakeEmail()
	stdout, stderr, err := Cmd(ctx).Exec(stdin, "import", "identities", "--format", "json", "--project", project, MakeRandomIdentity(t, email))
	require.NoError(t, err, stderr)
	out := gjson.Parse(stdout)
	assert.True(t, gjson.Valid(stdout))
	assert.Equal(t, email, out.Get("traits.username").String())
	return out.Get("id").String()
}

func ListIdentities(ctx context.Context, t testing.TB, project string) gjson.Result {
	stdout, stderr, err := Cmd(ctx).Exec(nil, "list", "identities", "--format", "json", "--project", project)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}

func CreateClient(ctx context.Context, t testing.TB, project string) gjson.Result {
	stdout, stderr, err := Cmd(ctx).Exec(nil, "create", "client", "--format", "json", "--project", project)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}

func BrowserLogin(t testing.TB, page playwright.Page, email, password string) {
	_, err := page.Goto(client.CloudConsoleURL("").String() + "/login")
	require.NoError(t, err)
	require.NoError(t, page.Locator(`[data-testid="node/input/identifier"] input`).Fill(email))
	require.NoError(t, page.Locator(`[data-testid="node/input/password"] input`).Fill(password))
	require.NoError(t, page.Locator(`[type="submit"][name="method"][value="password"]`).Click())
}

func RegisterAccount(ctx context.Context, t testing.TB) (email, password, name, sessionToken string) {
	email, password, name = FakeAccount()
	c := client.NewPublicOryProjectClient()

	flow, _, err := c.FrontendAPI.CreateNativeRegistrationFlow(ctx).Execute()
	require.NoError(t, err)

	res, _, err := c.FrontendAPI.
		UpdateRegistrationFlow(ctx).
		Flow(flow.Id).
		UpdateRegistrationFlowBody(cloud.UpdateRegistrationFlowBody{UpdateRegistrationFlowWithPasswordMethod: &cloud.UpdateRegistrationFlowWithPasswordMethod{
			Method:   "password",
			Password: password,
			Traits: map[string]any{
				"email": email,
				"name":  name,
				"consent": map[string]any{
					"tos": time.Now().UTC().Format(time.RFC3339),
				},
			},
		}}).
		Execute()
	require.NoError(t, err)
	require.NotNil(t, res.SessionToken)

	return email, password, name, *res.SessionToken
}
