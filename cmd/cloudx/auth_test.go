// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd"
	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/pointerx"
)

func TestAuthenticator(t *testing.T) {
	configDir := testhelpers.NewConfigFile(t)

	t.Run("errors without config and --quiet flag", func(t *testing.T) {
		c := cmd.NewRootCmd()
		c.SetArgs([]string{"auth", "--" + client.FlagConfig, configDir, "--quiet"})
		require.Error(t, c.Execute())
	})

	password := testhelpers.FakePassword()
	cmd := testhelpers.CmdWithConfigPassword(configDir, password)

	signIn := func(t *testing.T, email string) (string, string, error) {
		testhelpers.ClearConfig(t, configDir)
		var r bytes.Buffer

		_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
		_, _ = r.WriteString(email + "\n") // Email: FakeEmail()

		return cmd.Exec(&r, "auth")
	}

	t.Run("success", func(t *testing.T) {
		email := testhelpers.FakeEmail()
		name := testhelpers.FakeName()

		// Create the account
		r := testhelpers.RegistrationBuffer(name, email)
		stdout, stderr, err := cmd.Exec(r, "auth")
		require.NoError(t, err)

		assert.Contains(t, stderr, "You are now signed in as: "+email, "Expected to be signed in but response was:\n\t%s\n\tstderr: %s", stdout, stderr)
		assert.Contains(t, stdout, email)
		testhelpers.AssertConfig(t, configDir, email, name)
		testhelpers.ClearConfig(t, configDir)

		expectSignInSuccess := func(t *testing.T) {
			stdout, _, err := signIn(t, email)
			require.NoError(t, err)

			assert.Contains(t, stderr, "You are now signed in as: ", email, stdout)
			testhelpers.AssertConfig(t, configDir, email, name)
		}

		t.Run("sign in with valid data", func(t *testing.T) {
			expectSignInSuccess(t)
		})

		t.Run("forced to reauthenticate on session expiration", func(t *testing.T) {
			cmd := testhelpers.CmdWithConfig(configDir)
			expectSignInSuccess(t)
			testhelpers.ChangeAccessToken(t, configDir)
			var r bytes.Buffer
			r.WriteString("n\n") // Your CLI session has expired. Do you wish to login again as <email>?
			_, stderr, err := cmd.Exec(&r, "list", "projects")
			require.Error(t, err)
			assert.Contains(t, stderr, "Your session has expired or has otherwise become invalid. Please re-authenticate to continue.")
		})

		t.Run("user is able to reauthenticate on session expiration", func(t *testing.T) {
			cmd := testhelpers.CmdWithConfig(configDir)
			expectSignInSuccess(t)
			testhelpers.ChangeAccessToken(t, configDir)
			var r bytes.Buffer
			r.WriteString("y\n") // Your CLI session has expired. Do you wish to login again as <email>?
			_, stderr, err := cmd.Exec(&r, "list", "projects")
			require.Error(t, err)
			assert.Contains(t, stderr, "Your session has expired or has otherwise become invalid. Please re-authenticate to continue.")
			expectSignInSuccess(t)
		})

		t.Run("expired session with quiet flag returns error", func(t *testing.T) {
			cmd := testhelpers.CmdWithConfig(configDir)
			expectSignInSuccess(t)
			testhelpers.ChangeAccessToken(t, configDir)
			_, stderr, err := cmd.Exec(nil, "list", "projects", "-q")
			require.Error(t, err)
			assert.Equal(t, "please run `ory auth` to initialize your configuration or remove the `--quiet` flag", err.Error())
			assert.NotContains(t, stderr, "Your session has expired or has otherwise become invalid. Please re-authenticate to continue.")
		})

		t.Run("set up 2fa", func(t *testing.T) {
			expectSignInSuccess(t)
			ac := testhelpers.ReadConfig(t, configDir)

			c, err := client.NewOryProjectClient()
			require.NoError(t, err)

			flow, _, err := c.FrontendAPI.CreateNativeSettingsFlow(context.Background()).XSessionToken(ac.SessionToken).Execute()
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

			_, _, err = c.FrontendAPI.UpdateSettingsFlow(context.Background()).XSessionToken(ac.SessionToken).Flow(flow.Id).UpdateSettingsFlowBody(cloud.UpdateSettingsFlowBody{
				UpdateSettingsFlowWithTotpMethod: &cloud.UpdateSettingsFlowWithTotpMethod{
					TotpCode: pointerx.Ptr(code),
					Method:   "totp",
				},
			}).Execute()
			require.NoError(t, err)
			testhelpers.ClearConfig(t, configDir)

			t.Run("sign in fails because second factor is missing", func(t *testing.T) {
				t.Skip("TODO")

				testhelpers.ClearConfig(t, configDir)

				var r bytes.Buffer

				_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
				_, _ = r.WriteString(email + "\n") // Email: FakeEmail()

				stdout, stderr, err := cmd.Exec(&r, "auth")
				require.Error(t, err, stdout)

				assert.Contains(t, stderr, "Please complete the second authentication challenge", stdout)
				_, err = os.Stat(configDir)
				assert.ErrorIs(t, err, os.ErrNotExist)
			})

			t.Run("sign in succeeds with second factor", func(t *testing.T) {
				t.Skip("TODO")

				testhelpers.ClearConfig(t, configDir)

				var r bytes.Buffer

				code, err := totp.GenerateCode(secret, time.Now())
				require.NoError(t, err)
				_, _ = r.WriteString("y\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: y
				_, _ = r.WriteString(email + "\n") // Email: FakeEmail()
				_, _ = r.WriteString(code + "\n")  // TOTP code

				stdout, stderr, err := cmd.Exec(&r, "auth")
				require.NoError(t, err, stdout)

				assert.Contains(t, stderr, "Please complete the second authentication challenge", stdout)
				assert.Contains(t, stderr, "You are now signed in as: ", email, stdout)
				testhelpers.AssertConfig(t, configDir, email, name)
			})
		})
	})

	t.Run("retry sign up on invalid data", func(t *testing.T) {
		testhelpers.ClearConfig(t, configDir)

		r0 := testhelpers.RegistrationBuffer(testhelpers.FakeName(), "not-an-email")

		// Redo the flow
		email := testhelpers.FakeEmail()
		name := testhelpers.FakeName()
		r1 := testhelpers.RegistrationBuffer(name, email)
		// on retry, we need to skip "Do you want to sign in to an existing Ory Network account? [y/n]: "
		_, _ = r1.ReadString('\n')

		stdout, stderr, err := cmd.Exec(io.MultiReader(r0, r1), "auth", "--"+client.FlagConfig, configDir)
		require.NoError(t, err)

		assert.Contains(t, stderr, "Your account creation attempt failed. Please try again!", stdout) // First try fails
		assert.Contains(t, stderr, "You are now signed in as: "+email, stdout)                        // Second try succeeds
		testhelpers.AssertConfig(t, configDir, email, name)
	})

	t.Run("sign in with invalid data", func(t *testing.T) {
		stdout, stderr, err := signIn(t, testhelpers.FakeEmail())
		require.Error(t, err, stdout)

		assert.Contains(t, stderr, "The provided credentials are invalid", stdout)
	})
}
