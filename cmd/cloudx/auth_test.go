package cloudx_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx"
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/pointerx"
)

func TestAuthenticator(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)

	t.Run("errors without config and --quiet flag", func(t *testing.T) {
		cmd := cloudx.NewRootCommand(new(cobra.Command), "", "")
		cmd.SetArgs([]string{"auth", "--" + client.ConfigFlag, configDir, "--quiet"})
		require.Error(t, cmd.Execute())
	})

	password := testhelpers.FakePassword()
	exec := testhelpers.ConfigPasswordAwareCmd(configDir, password)

	signIn := func(t *testing.T, email string) (string, string, error) {
		testhelpers.ClearConfig(t, configDir)
		var r bytes.Buffer

		_, _ = r.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
		_, _ = r.WriteString(email + "\n") // Email: FakeEmail()

		return exec.Exec(&r, "auth")
	}

	t.Run("success", func(t *testing.T) {
		email := testhelpers.FakeEmail()
		name := testhelpers.FakeName()

		// Create the account
		var r bytes.Buffer
		_, _ = r.WriteString("n\n")        // Do you already have an Ory Console account you wish to use? [y/n]: n
		_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
		_, _ = r.WriteString(name + "\n")  // Name: FakeName()
		_, _ = r.WriteString("n\n")        // Please inform me about platform and security updates? [y/n]: n
		_, _ = r.WriteString("n\n")        // I accept the Terms of Service [y/n]: n
		_, _ = r.WriteString("y\n")        // I accept the Terms of Service [y/n]: y

		stdout, stderr, err := exec.Exec(&r, "auth")
		require.NoError(t, err)

		assert.Contains(t, stderr, "You are now signed in as: "+email, "Expected to be signed in but response was:\n\t%s\n\tstderr: %s", stdout, stderr)
		assert.Contains(t, stdout, email)
		testhelpers.AssertConfig(t, configDir, email, name, false)
		testhelpers.ClearConfig(t, configDir)

		expectSignInSuccess := func(t *testing.T) {
			stdout, _, err := signIn(t, email)
			require.NoError(t, err)

			assert.Contains(t, stderr, "You are now signed in as: ", email, stdout)
			testhelpers.AssertConfig(t, configDir, email, name, false)
		}

		t.Run("sign in with valid data", func(t *testing.T) {
			expectSignInSuccess(t)
		})

		t.Run("forced to reauthenticate on session expiration", func(t *testing.T) {
			cmd := testhelpers.ConfigAwareCmd(configDir)
			expectSignInSuccess(t)
			testhelpers.ChangeAccessToken(t, configDir)
			var r bytes.Buffer
			r.WriteString("n\n") // Your CLI session has expired. Do you wish to login again as <email>?
			_, stderr, err := cmd.ExecDebug(t, &r, "list", "projects")
			require.Error(t, err)
			assert.Contains(t, stderr, "Your CLI session has expired. Do you wish to log in again as")
		})

		t.Run("user is able to reauthenticate on session expiration", func(t *testing.T) {
			cmd := testhelpers.ConfigAwareCmd(configDir)
			expectSignInSuccess(t)
			testhelpers.ChangeAccessToken(t, configDir)
			var r bytes.Buffer
			r.WriteString("y\n") // Your CLI session has expired. Do you wish to login again as <email>?
			_, stderr, err := cmd.ExecDebug(t, &r, "list", "projects")
			require.Error(t, err)
			assert.Contains(t, stderr, "Your CLI session has expired. Do you wish to log in again as")
			expectSignInSuccess(t)
		})

		t.Run("expired session with quiet flag returns error", func(t *testing.T) {
			cmd := testhelpers.ConfigAwareCmd(configDir)
			expectSignInSuccess(t)
			testhelpers.ChangeAccessToken(t, configDir)
			_, stderr, err := cmd.ExecDebug(t, nil, "list", "projects", "-q")
			require.Error(t, err)
			assert.Equal(t, "Your session has expired and you cannot reauthenticate when the --quiet flag is set", err.Error())
			assert.NotContains(t, stderr, "Your CLI session has expired. Do you wish to log in again as")
		})

		t.Run("set up 2fa", func(t *testing.T) {
			expectSignInSuccess(t)
			ac := testhelpers.ReadConfig(t, configDir)

			c, err := client.NewKratosClient()
			require.NoError(t, err)

			flow, _, err := c.V0alpha2Api.InitializeSelfServiceSettingsFlowWithoutBrowser(context.Background()).XSessionToken(ac.SessionToken).Execute()
			require.NoError(t, err)

			var secret string
			for _, node := range flow.Ui.Nodes {
				if node.Type != "text" {
					continue
				}

				attrs := node.Attributes.UiNodeTextAttributes
				if attrs.Text.Id == 1050006 {
					secret = attrs.Text.Text
				}
			}

			require.NotEmpty(t, secret)
			code, err := totp.GenerateCode(secret, time.Now())
			require.NoError(t, err)

			_, _, err = c.V0alpha2Api.SubmitSelfServiceSettingsFlow(context.Background()).XSessionToken(ac.SessionToken).Flow(flow.Id).SubmitSelfServiceSettingsFlowBody(cloud.SubmitSelfServiceSettingsFlowBody{
				SubmitSelfServiceSettingsFlowWithTotpMethodBody: &cloud.SubmitSelfServiceSettingsFlowWithTotpMethodBody{
					TotpCode: pointerx.String(code),
					Method:   "totp",
				},
			}).Execute()
			require.NoError(t, err)
			testhelpers.ClearConfig(t, configDir)

			t.Run("sign in fails because second factor is missing", func(t *testing.T) {
				testhelpers.ClearConfig(t, configDir)

				var r bytes.Buffer

				_, _ = r.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
				_, _ = r.WriteString(email + "\n") // Email: FakeEmail()

				stdout, stderr, err := exec.Exec(&r, "auth")
				require.Error(t, err, stdout)

				assert.Contains(t, stderr, "Please complete the second authentication challenge", stdout)
				_, err = os.Stat(configDir)
				assert.ErrorIs(t, err, os.ErrNotExist)
			})

			t.Run("sign in succeeds with second factor", func(t *testing.T) {
				testhelpers.ClearConfig(t, configDir)

				var r bytes.Buffer

				code, err := totp.GenerateCode(secret, time.Now())
				require.NoError(t, err)
				_, _ = r.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
				_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
				_, _ = r.WriteString(code + "\n")  // TOTP code

				stdout, stderr, err := exec.Exec(&r, "auth")
				require.NoError(t, err, stdout)

				assert.Contains(t, stderr, "Please complete the second authentication challenge", stdout)
				assert.Contains(t, stderr, "You are now signed in as: ", email, stdout)
				testhelpers.AssertConfig(t, configDir, email, name, false)
			})
		})
	})

	t.Run("retry sign up on invalid data", func(t *testing.T) {
		testhelpers.ClearConfig(t, configDir)

		var r bytes.Buffer

		_, _ = r.WriteString("n\n")                         // Do you already have an Ory Console account you wish to use? [y/n]: n
		_, _ = r.WriteString("not-an-email" + "\n")         // Email: FakeEmail()
		_, _ = r.WriteString(testhelpers.FakeName() + "\n") // Name: FakeName()
		_, _ = r.WriteString("n\n")                         // Please inform me about platform and security updates? [y/n]: n
		_, _ = r.WriteString("y\n")                         // I accept the Terms of Service [y/n]: y

		// Redo the flow
		email := testhelpers.FakeEmail()
		name := testhelpers.FakeName()
		_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
		_, _ = r.WriteString(name + "\n")  // Name: FakeName()
		_, _ = r.WriteString("y\n")        // Please inform me about platform and security updates? [y/n]: n
		_, _ = r.WriteString("y\n")        // I accept the Terms of Service [y/n]: y

		stdout, stderr, err := exec.Exec(&r, "auth", "--"+client.ConfigFlag, configDir)
		require.NoError(t, err)

		assert.Contains(t, stderr, "Your account creation attempt failed. Please try again!", stdout) // First try fails
		assert.Contains(t, stderr, "You are now signed in as: "+email, stdout)                        // Second try succeeds
		testhelpers.AssertConfig(t, configDir, email, name, true)
	})

	t.Run("sign in with invalid data", func(t *testing.T) {
		stdout, stderr, err := signIn(t, testhelpers.FakeEmail())
		require.Error(t, err, stdout)

		assert.Contains(t, stderr, "The provided credentials are invalid", stdout)
	})
}
