// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/assertx"
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

func ClearConfig(t testing.TB, configDir string) {
	require.NoError(t, os.RemoveAll(configDir))
}

func AssertConfig(t testing.TB, configDir string, email string, name string) {
	ac := ReadConfig(t, configDir)
	assert.Equal(t, email, ac.IdentityTraits.Email)
	assert.Equal(t, client.ConfigVersion, ac.Version)
	assert.NotEmpty(t, ac.SessionToken)

	c, err := client.NewOryProjectClient()
	require.NoError(t, err)

	res, _, err := c.FrontendAPI.ToSession(context.Background()).XSessionToken(ac.SessionToken).Execute()
	require.NoError(t, err)

	traits, err := json.Marshal(res.Identity.Traits)
	require.NoError(t, err)

	assertx.EqualAsJSONExcept(t, json.RawMessage(`{
  "email": "`+email+`",
  "name": "`+name+`"
}`), json.RawMessage(traits), []string{"consent", "details"})
	assert.NotEmpty(t, gjson.GetBytes(traits, "consent.tos").String())
}

func CmdWithConfig(configDir string) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{
		New:            cmd.NewRootCmd,
		Ctx:            client.ContextWithClient(context.Background()),
		PersistentArgs: []string{"--" + client.FlagConfig, configDir},
	}
}

func CmdWithConfigPassword(configPath, password string) *cmdx.CommandExecuter {
	c := CmdWithConfig(configPath)
	c.Ctx = client.ContextWithOptions(c.Ctx, client.WithPasswordReader(func() ([]byte, error) {
		return []byte(password), nil
	}))
	return c
}

func ChangeAccessToken(t testing.TB, configDir string) {
	ac := ReadConfig(t, configDir)
	ac.SessionToken = "12341234"
	data, err := json.Marshal(ac)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configDir, data, 0644))
}

func RegisterAccount(t testing.TB, configPath string) (email, password, name string) {
	email, password, name = FakeAccount()

	stdout, stderr, err := CmdWithConfigPassword(configPath, password).Exec(RegistrationBuffer(name, email), "auth")
	require.NoError(t, err)
	require.Contains(t, stderr, "You are now signed in as: "+email, stdout)

	AssertConfig(t, configPath, email, name)

	return email, password, name
}

func RegistrationBuffer(name string, email string) *bytes.Buffer {
	var r bytes.Buffer
	_, _ = r.WriteString("n\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: n
	_, _ = r.WriteString(email + "\n") // Work email: FakeEmail()
	// Password is read through the password reader
	_, _ = r.WriteString(name + "\n") // Name: FakeName()
	_, _ = r.WriteString("n\n")       // Please inform me about platform and security updates:  [y/n]: n
	_, _ = r.WriteString("y\n")       // I accept the Terms of Service https://www.ory.sh/ptos:  [y/n]: y
	return &r
}

func LoginBuffer(email string) io.Reader {
	var r bytes.Buffer
	_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
	_, _ = r.WriteString(email + "\n") // Email FakeEmail()
	return &r
}

func WithReAuth(t testing.TB, email, password string) (*cmdx.CommandExecuter, io.Reader) {
	return CmdWithConfigPassword(NewConfigFile(t), password), LoginBuffer(email)
}

func CreateProject(t testing.TB, configDir string, workspace *string) *cloud.Project {
	args := []string{"create", "project", "--name", TestName(), "--format", "json"}
	if workspace != nil {
		args = append(args, "--workspace", *workspace)
	}
	stdout, stderr, err := CmdWithConfig(configDir).Exec(nil, args...)
	require.NoError(t, err, stderr)
	p := cloud.Project{}
	require.NoError(t, json.Unmarshal([]byte(stdout), &p), stdout)
	if ap, ok := p.AdditionalProperties["AdditionalProperties"]; ok {
		// the SDK types are weird sometimes...
		p.AdditionalProperties = ap.(map[string]interface{})
	}
	return &p
}

func CreateWorkspace(t testing.TB, configDir string) string {
	return strings.TrimSpace(CmdWithConfig(configDir).ExecNoErr(t, "create", "workspace", "--name", TestName(), "--quiet"))
}

func SetDefaultProject(t testing.TB, configPath string, projectID string) {
	require.Equal(t, projectID, strings.TrimSpace(CmdWithConfig(configPath).ExecNoErr(t, "use", "project", projectID, "--quiet")))
}

func GetDefaultProjectID(t testing.TB, configDir string) string {
	return strings.TrimSpace(CmdWithConfig(configDir).ExecNoErr(t, "use", "project", "--quiet"))
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

func ImportIdentity(t testing.TB, cmd *cmdx.CommandExecuter, project string, stdin io.Reader) string {
	email := FakeEmail()
	stdout, stderr, err := cmd.Exec(stdin, "import", "identities", "--format", "json", "--project", project, MakeRandomIdentity(t, email))
	require.NoError(t, err, stderr)
	out := gjson.Parse(stdout)
	assert.True(t, gjson.Valid(stdout))
	assert.Equal(t, email, out.Get("traits.username").String())
	return out.Get("id").String()
}

func CreateClient(t testing.TB, cmd *cmdx.CommandExecuter, project string) gjson.Result {
	stdout, stderr, err := cmd.Exec(nil, "create", "client", "--format", "json", "--project", project)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}
