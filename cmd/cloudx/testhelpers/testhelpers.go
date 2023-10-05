// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

func TestProjectName() string {
	return testProjectPattern + randx.MustString(16, randx.AlphaLowerNum)
}

func FakeEmail() string {
	return fmt.Sprintf(testAccountPrefix+".%s@ory.dev", randx.MustString(16, randx.AlphaLowerNum))
}

func FakePassword() string {
	return randx.MustString(16, randx.AlphaLowerNum)
}

func FakeName() string {
	return randx.MustString(16, randx.AlphaLowerNum)
}

func NewConfigDir(t require.TestingT) string {
	homeDir, err := os.MkdirTemp(os.TempDir(), "cloudx-*")
	require.NoError(t, err)
	return filepath.Join(homeDir, "config.json")
}

func ReadConfig(t require.TestingT, configDir string) *client.AuthContext {
	f, err := os.ReadFile(configDir)
	require.NoError(t, err)
	var ac client.AuthContext
	require.NoError(t, json.Unmarshal(f, &ac))
	return &ac
}

func ClearConfig(t require.TestingT, configDir string) {
	require.NoError(t, os.RemoveAll(configDir))
}

func AssertConfig(t *testing.T, configDir string, email string, name string, newsletter bool) {
	ac := ReadConfig(t, configDir)
	assert.Equal(t, email, ac.IdentityTraits.Email)
	assert.Equal(t, client.Version, ac.Version)
	assert.NotEmpty(t, ac.SessionToken)

	c, err := client.NewKratosClient()
	require.NoError(t, err)

	res, _, err := c.FrontendAPI.ToSession(context.Background()).XSessionToken(ac.SessionToken).Execute()
	require.NoError(t, err)

	traits, err := json.Marshal(res.Identity.Traits)
	require.NoError(t, err)

	assertx.EqualAsJSONExcept(t, json.RawMessage(`{
  "email": "`+email+`",
  "name": "`+name+`",
  "consent": {
    "newsletter": `+fmt.Sprintf("%v", newsletter)+`,
    "tos": ""
  }
}`), json.RawMessage(traits), []string{"consent.tos"})
	assert.NotEmpty(t, gjson.GetBytes(traits, "consent.tos").String())
}

func ConfigAwareCmd(configDir string) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{
		New:            cmd.NewRootCmd,
		Ctx:            client.ContextWithClient(context.Background()),
		PersistentArgs: []string{"--" + client.ConfigFlag, configDir},
	}
}

func ConfigPasswordAwareCmd(configDir, password string) *cmdx.CommandExecuter {
	ctx := client.ContextWithClient(context.WithValue(context.Background(), client.PasswordReader{}, func() ([]byte, error) {
		return []byte(password), nil
	}))
	return &cmdx.CommandExecuter{
		New:            cmd.NewRootCmd,
		Ctx:            ctx,
		PersistentArgs: []string{"--" + client.ConfigFlag, configDir},
	}
}

func ChangeAccessToken(t require.TestingT, configDir string) {
	ac := ReadConfig(t, configDir)
	ac.SessionToken = "12341234"
	data, err := json.Marshal(ac)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configDir, data, 0644))
}

func RegisterAccount(t require.TestingT, configDir string) (email, password string) {
	password = FakePassword()
	email = FakeEmail()
	name := FakeName()

	// Create the account
	var r bytes.Buffer
	_, _ = r.WriteString("n\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: n
	_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
	_, _ = r.WriteString(name + "\n")  // Name: FakeName()
	_, _ = r.WriteString("n\n")        // Subscribe to the Ory Security Newsletter to get platform and security updates? [y/n]: n
	_, _ = r.WriteString("n\n")        // I accept the Terms of Service [y/n]: n
	_, _ = r.WriteString("y\n")        // I accept the Terms of Service [y/n]: y

	exec := cmdx.CommandExecuter{
		New: cmd.NewRootCmd,
		Ctx: context.WithValue(context.Background(), client.PasswordReader{}, func() ([]byte, error) {
			return []byte(password), nil
		}),
		PersistentArgs: []string{"--" + client.ConfigFlag, configDir},
	}

	stdout, stderr, err := exec.Exec(&r, "auth")
	require.NoError(t, err)

	assert.Contains(t, stderr, "You are now signed in as: "+email, stdout)
	if t, ok := t.(*testing.T); ok {
		AssertConfig(t, configDir, email, name, false)
	}
	return email, password
}

func WithReAuth(t require.TestingT, email, password string) (*cmdx.CommandExecuter, *bytes.Buffer) {
	configDir := NewConfigDir(t)
	cmd := ConfigPasswordAwareCmd(configDir, password)
	// Create the account
	var r bytes.Buffer
	r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
	r.WriteString(email + "\n") // Email FakeEmail()
	return cmd, &r
}

func CreateProject(t require.TestingT, configDir string) string {
	cmd := ConfigAwareCmd(configDir)
	name := TestProjectName()
	stdout, stderr, err := cmd.Exec(nil, "create", "project", "--name", name, "--format", "json")
	require.NoError(t, err, "stdout: %s\nstderr: %s", stderr)
	id := gjson.Get(stdout, "id").String()
	return id
}

func CreateAndUseProject(t require.TestingT, configDir string) string {
	cmd := ConfigAwareCmd(configDir)
	name := TestProjectName()
	stdout, stderr, err := cmd.Exec(nil, "create", "project", "--name", name, "--use-project", "--format", "json")
	require.NoError(t, err, "stdout: %s\nstderr: %s", stderr)
	ac := ReadConfig(t, configDir)
	id := gjson.Get(stdout, "id").String()
	assert.Equal(t, ac.SelectedProject.String(), id)
	return id
}

func SetDefaultProject(t require.TestingT, configDir string, projectId string) {
	cmd := ConfigAwareCmd(configDir)
	stdout, stderr, err := cmd.Exec(nil, "use", "project", projectId, "--format", "json")
	require.NoError(t, err, "stdout: %s\nstderr: %s", stderr)
	ac := ReadConfig(t, configDir)
	id := gjson.Get(stdout, "id").String()
	assert.Equal(t, ac.SelectedProject.String(), id)
	assert.Equal(t, projectId, id)
}

func GetDefaultProject(t require.TestingT, configDir string) string {
	cmd := ConfigAwareCmd(configDir)
	stdout, stderr, err := cmd.Exec(nil, "use", "project", "--format", "json")
	require.NoError(t, err, "stdout: %s\nstderr: %s", stderr)
	ac := ReadConfig(t, configDir)
	id := gjson.Get(stdout, "id").String()
	assert.Equal(t, ac.SelectedProject.String(), id)
	return id
}

func MakeRandomIdentity(t require.TestingT, email string) string {
	homeDir, err := os.MkdirTemp(os.TempDir(), "cloudx-*")
	require.NoError(t, err)
	path := filepath.Join(homeDir, "import.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
  "schema_id": "preset://username",
  "traits": {
    "username": "`+email+`"
  }
}`), 0600))
	return path
}

func MakeRandomClient(t require.TestingT, name string) string {
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

func ImportIdentity(t require.TestingT, cmd *cmdx.CommandExecuter, project string, stdin *bytes.Buffer) string {
	email := FakeEmail()
	stdout, stderr, err := cmd.Exec(stdin, "import", "identities", "--format", "json", "--project", project, MakeRandomIdentity(t, email))
	require.NoError(t, err, stderr)
	out := gjson.Parse(stdout)
	assert.True(t, gjson.Valid(stdout))
	assert.Equal(t, email, out.Get("traits.username").String())
	return out.Get("id").String()
}

func CreateClient(t require.TestingT, cmd *cmdx.CommandExecuter, project string) gjson.Result {
	stdout, stderr, err := cmd.Exec(nil, "create", "client", "--format", "json", "--project", project)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}
