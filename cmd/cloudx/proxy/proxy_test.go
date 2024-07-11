// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy_test

import (
	"bytes"
	"context"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	. "github.com/ory/x/pointerx"
)

func TestProxy(t *testing.T) {
	t.Skip("TODO: not yet finished")

	const baseURL, pathPrefix = "http://localhost:4000", "/.ory"

	pw, err := playwright.Run()
	require.NoError(t, err)
	t.Cleanup(func() {
		t.Logf("playwright stop error %+v", pw.Stop())
	})
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: Ptr(true),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		t.Logf("browser close error %+v", browser.Close())
	})
	pwCtx, err := browser.NewContext()
	require.NoError(t, err)
	t.Cleanup(func() {
		t.Logf("context close error %+v", pwCtx.Close())
	})
	page, err := pwCtx.NewPage()
	require.NoError(t, err)
	t.Cleanup(func() {
		t.Logf("page close error %+v", page.Close())
	})

	_, password, _, wsID, pjID := testhelpers.BrowserRegistration(t, page)
	ctx := client.ContextWithOptions(context.Background(),
		client.WithConfigLocation(testhelpers.NewConfigFile(t)),
		client.WithOpenBrowserHook(testhelpers.PlaywrightAcceptConsentBrowserHook(t, page, password)),
	)

	cmd := testhelpers.Cmd(ctx)

	stdOut, stdErr := bytes.Buffer{}, bytes.Buffer{}
	cmd.ExecBackground(nil, &stdOut, &stdErr, "proxy", "https://ory-network-httpbin-ijakee5waq-ez.a.run.app/anything", "--rewrite-host", "--quiet", "--workspace", wsID, "--project", pjID)

	require.Eventually(t, func() bool {
		return strings.Contains(stdErr.String(), "access your application via the Ory Proxy")
	}, 5*time.Second, 100*time.Millisecond)

	// we don't need the authenticated context anymore
	t.Logf("pw context close error %+v", pwCtx.Close())

	pwCtx, err = browser.NewContext(playwright.BrowserNewContextOptions{
		BaseURL: Ptr(baseURL),
	})
	require.NoError(t, err)
	page, err = pwCtx.NewPage()
	require.NoError(t, err)

	assertIsLoggedIn := func() func(t *testing.T) {
		resp, err := http.Get(baseURL + "/.ory/proxy/jwks.json")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		// TODO parse keys

		return func(t *testing.T) {
			// the page behind is httpbin
			resp, err := pwCtx.Request().Get(baseURL)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.Status)

			body, err := resp.Body()
			require.NoError(t, err)

			authorizationHeader := gjson.GetBytes(body, "headers.Authorization")
			assert.Truef(t, authorizationHeader.Exists(), "full body: %s", body)
			// TODO parse and verify token
		}
	}()

	t.Run("navigation works", func(t *testing.T) {
		_, err := page.Goto(pathPrefix + "/ui/registration")
		require.NoError(t, err)
		require.NoError(t, page.Locator(`[data-testid="cta-link"]`).Click())
		require.NoError(t, page.WaitForURL(regexp.MustCompile(baseURL+pathPrefix+"/ui/login\\?.*")))
	})

	t.Run("should be able to execute registration", func(t *testing.T) {
		_, err := page.Goto(pathPrefix + "/ui/registration")
		require.NoError(t, err)
		require.NoError(t, page.Locator(`[name="traits.email"]`).Fill(testhelpers.FakeEmail()))
		require.NoError(t, page.Locator(`[name="password"]`).Fill(testhelpers.FakePassword()))
		require.NoError(t, page.Locator(`[name="method"]`).Click())

		assertIsLoggedIn(t)
	})
}
