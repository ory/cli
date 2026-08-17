// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/randx"
)

const (
	testProjectPattern = "ory-cy-e2e-da2f162d-af61-42dd-90dc-e3fcfa7c84a0-"
	testAccountPrefix  = "dev+orycye2eda2f162daf6142dd"
)

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

func WithDuplicatedConfigFile(ctx context.Context, t testing.TB, originalFile string) (context.Context, string) {
	dst, err := os.Create(NewConfigFile(t))
	require.NoError(t, err)
	defer dst.Close()
	src, err := os.Open(originalFile)
	require.NoError(t, err)
	defer src.Close()
	_, err = io.Copy(dst, src)
	require.NoError(t, err)

	return client.ContextWithOptions(ctx, client.WithConfigLocation(dst.Name())), dst.Name()
}

func Cmd(ctx context.Context) *cmdx.CommandExecuter {
	return &cmdx.CommandExecuter{
		New: cmd.NewRootCmd,
		Ctx: client.ContextWithClient(ctx),
	}
}

func CreateProject(ctx context.Context, t testing.TB, workspace string) *cloud.Project {
	args := []string{"create", "project", "--name", TestName(), "--workspace", workspace, "--format", "json", "--environment", "dev"}
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
}`), 0o600))
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
]`), 0o600))
	return path
}

func ImportIdentity(ctx context.Context, t testing.TB, project string, stdin io.Reader) string {
	email := FakeEmail()
	args := []string{"import", "identities", "--format", "json", MakeRandomIdentity(t, email)}
	if project != "" {
		args = append(args, "--project", project)
	}
	stdout, stderr, err := Cmd(ctx).Exec(stdin, args...)
	require.NoError(t, err, stderr)
	out := gjson.Parse(stdout)
	assert.True(t, gjson.Valid(stdout))
	assert.Equal(t, email, out.Get("traits.username").String())
	return out.Get("id").String()
}

func ListIdentities(ctx context.Context, t testing.TB, project string) gjson.Result {
	args := []string{"list", "identities", "--format", "json"}
	if project != "" {
		args = append(args, "--project", project)
	}
	stdout, stderr, err := Cmd(ctx).Exec(nil, args...)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}

func ListRelationTuples(ctx context.Context, t testing.TB, project string) gjson.Result {
	args := []string{"list", "relation-tuples", "--format", "json"}
	if project != "" {
		args = append(args, "--project", project)
	}
	stdout, stderr, err := Cmd(ctx).Exec(nil, args...)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}

func ListClients(ctx context.Context, t testing.TB, project string) gjson.Result {
	args := []string{"list", "clients", "--format", "json"}
	if project != "" {
		args = append(args, "--project", project)
	}
	stdout, stderr, err := Cmd(ctx).Exec(nil, args...)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}

func CreateClient(ctx context.Context, t testing.TB, project string) gjson.Result {
	args := []string{"create", "client", "--format", "json"}
	if project != "" {
		args = append(args, "--project", project)
	}
	stdout, stderr, err := Cmd(ctx).Exec(nil, args...)
	require.NoError(t, err, stderr)
	return gjson.Parse(stdout)
}

func RegisterAccount(ctx context.Context, t testing.TB) (email, password, name string) {
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

	return email, password, name
}

func SetupPlaywright(t testing.TB) (playwright.Browser, playwright.Page, func()) {
	pw, err := playwright.Run()
	require.NoError(t, err)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless:  new(true),
		TracesDir: new(tracesDir),
	})
	require.NoError(t, err)

	page := NewPage(t, browser)

	return browser, page, func() {
		// Drain any in-flight Route handlers (registered in NewPage) before closing
		// the page; otherwise page.Close() can race the handler goroutine inside
		// playwright-go.
		t.Logf("unroute error: %+v", page.Context().UnrouteAll(playwright.BrowserContextUnrouteAllOptions{
			Behavior: playwright.UnrouteBehaviorWait,
		}))
		t.Logf("page close error: %+v", page.Close())
		t.Logf("browser close error: %+v", browser.Close())
		t.Logf("playwright stop error: %+v", pw.Stop())
	}
}

func NewPage(t testing.TB, browser playwright.Browser) playwright.Page {
	opts := playwright.BrowserNewPageOptions{
		BaseURL: new(client.CloudConsoleURL("").String()),
	}

	// The browser talks to Ory Network directly, so it needs the same rate-limit
	// exemption the SDK clients get — the header configured on those does not
	// reach it. Without this, `go test ./...` runs the browser login of every
	// cloudx package at once from a single CI egress IP and the login endpoint
	// starts answering 429.
	if name, value, ok := client.RateLimitHeader(); ok {
		opts.ExtraHttpHeaders = map[string]string{name: value}
	}

	page, err := browser.NewPage(opts)
	require.NoError(t, err)

	for _, route := range []string{
		"doubleclick.net",
		"google-analytics.com",
		"googletagmanager.com",
		"hs-analytics.net",
		"hs-banner.com",
		"hs-scripts.com",
		"hsadspixel.net",
		"hubapi.com",
		"hubapi.com",
		"licdn.com",
		"linkedin.com",
		"eu.posthog.com",
		"r.stripe.com",
		"segment.com",
		"sentry.io",
		"sst.ory.sh",
		"www.google.com/pagead",
		"app.termly.io",
		"static.reo.dev",
		"api.reo.dev",
	} {
		require.NoError(t, page.Context().Route(func(actual string) bool {
			return strings.Contains(actual, route)
		}, func(r playwright.Route) {
			// Best-effort: the page may already be closed during teardown, in
			// which case Abort returns "target closed" for a request we no
			// longer care about. Failing the test here races test completion.
			_ = r.Abort()
		}))
	}
	return page
}

// stopTracing writes the trace of the login flow and strips the rate-limit
// header value out of it.
//
// A trace records complete request headers — this repository is public and CI
// uploads the traces as a build artifact, where GitHub's secret masking does not
// reach. The header that exempts CI from Ory Network's rate limits therefore
// must not survive into one. Playwright offers no redaction option, so the
// archive is rewritten after it has been written.
//
// If it cannot be rewritten the trace is deleted: losing a diagnostic is the
// cheaper failure by far.
func stopTracing(t testing.TB, page playwright.Page) {
	path := filepath.Join(tracesDir, fmt.Sprintf("%s.zip", t.Name()))
	if err := page.Context().Tracing().Stop(path); err != nil {
		t.Logf("tracing stop error: %+v", err)
		return
	}

	_, secret, ok := client.RateLimitHeader()
	if !ok {
		return
	}
	if err := redactInZip(path, secret); err != nil {
		t.Logf("could not redact %s, removing it: %+v", path, err)
		require.NoError(t, os.Remove(path))
	}
}

const redactedPlaceholder = "[redacted]"

// redactInZip rewrites every entry of the zip archive at path, replacing each
// occurrence of secret with a placeholder.
//
// The JSON-escaped spelling is replaced as well, because the trace stores
// headers as JSON string values and a secret containing a quote or backslash
// would otherwise appear there in a form the raw comparison does not match.
func redactInZip(path, secret string) error {
	if secret == "" {
		return nil
	}

	needles := [][]byte{[]byte(secret)}
	if escaped, err := json.Marshal(secret); err == nil {
		if inner := escaped[1 : len(escaped)-1]; !bytes.Equal(inner, []byte(secret)) {
			needles = append(needles, inner)
		}
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "trace-*.zip")
	if err != nil {
		_ = r.Close()
		return err
	}
	defer os.Remove(tmp.Name()) // no-op once the rename below succeeded

	err = func() error {
		w := zip.NewWriter(tmp)
		for _, f := range r.File {
			src, err := f.Open()
			if err != nil {
				return err
			}
			content, err := io.ReadAll(src)
			_ = src.Close()
			if err != nil {
				return err
			}
			for _, needle := range needles {
				content = bytes.ReplaceAll(content, needle, []byte(redactedPlaceholder))
			}
			dst, err := w.Create(f.Name)
			if err != nil {
				return err
			}
			if _, err := dst.Write(content); err != nil {
				return err
			}
		}
		return w.Close()
	}()
	_ = r.Close()
	_ = tmp.Close()
	if err != nil {
		return err
	}

	return os.Rename(tmp.Name(), path)
}

// submitPasswordForm submits the filled-in login form and fails immediately if
// the server refused the request outright.
//
// This catches the infrastructure-level refusals — 429 when a CI run drives more
// logins from one IP than Ory Network allows, or a 5xx — and it exists because
// the consent screen is the next thing the caller waits for. On a refusal the
// page never leaves the form, so an unchecked failure surfaces 30 seconds later
// as a missing `Allow` button, pointing at the consent screen instead of at the
// reason it never rendered.
//
// A rejected *credential* is not covered here: Ory Network answers that with a
// 303 back to the login page, which is indistinguishable from success at this
// point. The consent wait reports where the browser ended up, which is what
// separates the two.
func submitPasswordForm(t testing.TB, page playwright.Page) {
	// The login flow issues exactly one request to this path, and it is this
	// submission — the flow itself is created by the login UI, not the browser.
	isLoginSubmission := func(url string) bool {
		return strings.Contains(url, "/self-service/login")
	}

	resp, err := page.ExpectResponse(isLoginSubmission, func() error {
		return page.Locator(`[type="submit"][name="method"][value="password"]`).Click()
	})
	require.NoError(t, err, "the login form was never submitted")

	if resp.Status() < http.StatusBadRequest {
		return
	}

	body, _ := resp.Text()
	require.FailNowf(t, "the login form was refused",
		"POST %s\n%d %s\n%s\n\nThe consent screen never renders after this, so waiting for it would only time out.",
		resp.URL(), resp.Status(), resp.StatusText(), body)
}

func PlaywrightAcceptConsentBrowserHook(t testing.TB, page playwright.Page, email, password string) func(uri string) error {
	return func(uri string) error {
		t.Logf("open browser with %s", uri)

		require.NoError(t, page.Context().Tracing().Start(playwright.TracingStartOptions{
			Screenshots: new(true),
			Snapshots:   new(true),
		}))
		defer func() {
			r := recover()
			stopTracing(t, page)
			if r != nil {
				panic(r)
			}
		}()

		_, err := page.Goto(uri)
		require.NoError(t, err)

		if err := page.Locator(`[data-testid="node/input/identifier"] input`).WaitFor(); err == nil {
			// we need to log in first
			t.Logf("logging in")
			require.NoError(t, page.Locator(`[data-testid="node/input/identifier"] input`).Fill(email))
			require.NoError(t, page.Locator(`[data-testid="node/input/password"] input`).Fill(password))
		} else {
			// reconfirm password
			t.Logf("reconfirming password")
			require.NoError(t, page.Locator(`[data-testid="node/input/password"] input`).Fill(password))
		}
		submitPasswordForm(t, page)

		// we wait here for the button +1s because there is some console bug that can lead to form submissions before the form action is correctly set
		if err := page.Locator(`button:has-text("Allow")`).WaitFor(); err != nil {
			require.FailNowf(t, "the consent screen did not render", "%s\n\nThe browser ended up at %s. Still being on the login page means the credentials or the flow were rejected, rather than the consent screen itself being broken.",
				err, page.URL())
		}
		time.Sleep(time.Second)

		// accept consent
		require.NoError(t, page.Locator(`button:has-text("Allow")`).Click())

		t.Logf("consent successful")

		return nil
	}
}

var tracesDir string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dirs := strings.Split(cwd, string(os.PathSeparator))
	for i := range dirs {
		if dirs[i] == "cloudx" {
			dirs = dirs[:i-1]
			break
		}
	}
	tracesDir = string(os.PathSeparator) + filepath.Join(append(dirs, "playwright-traces")...)
}
