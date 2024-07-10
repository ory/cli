// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/randx"

	"github.com/ory/cli/cmd/cloudx/client"
)

func setEnvIfUnset(key, value string) {
	if _, ok := os.LookupEnv(key); !ok {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}

func UseStaging() {
	setEnvIfUnset(client.ConsoleURLKey, "https://console.staging.ory.dev")
	setEnvIfUnset(client.OryAPIsURLKey, "https://staging.oryapis.dev")
}

func CreateDefaultAssetsBrowser() (ctx context.Context, defaultConfig string, extraProject, defaultProject *cloud.Project, defaultCmd *cmdx.CommandExecuter) {
	UseStaging()

	t := MockTestingTForMain{}

	defaultConfig = NewConfigFile(t)

	email, password, _, sessionToken := RegisterAccount(context.Background(), t)
	ctx = client.ContextWithOptions(context.Background(),
		client.WithConfigLocation(defaultConfig),
		client.WithSessionToken(t, sessionToken))

	defaultProject = CreateProject(ctx, t, nil)
	extraProject = CreateProject(ctx, t, nil)

	_, page, cleanup := SetupPlaywright(t)
	defer cleanup()
	BrowserLogin(t, page, email, password)

	ctx = client.ContextWithOptions(context.Background(), client.WithConfigLocation(NewConfigFile(t)))
	h, err := client.NewCommandHelper(
		ctx,
		client.WithQuiet(false),
		client.WithOpenBrowserHook(PlaywrightAcceptConsentBrowserHook(t, page, password)),
	)
	require.NoError(t, err)
	require.NoError(t, h.Authenticate(ctx))

	defaultCmd = Cmd(ctx)
	return
}

func CreateDefaultAssets() (ctx context.Context, defaultConfig string, extraProject, defaultProject *cloud.Project, defaultCmd *cmdx.CommandExecuter) {
	UseStaging()

	t := MockTestingTForMain{}

	defaultConfig = NewConfigFile(t)

	_, _, _, sessionToken := RegisterAccount(context.Background(), t)
	ctx = client.ContextWithOptions(context.Background(),
		client.WithConfigLocation(defaultConfig),
		client.WithSessionToken(t, sessionToken),
		client.WithOpenBrowserHook(func(uri string) error {
			return errors.WithStack(fmt.Errorf("open browser hook not expected: %s", uri))
		}))

	defaultProject = CreateProject(ctx, t, nil)
	extraProject = CreateProject(ctx, t, nil)
	defaultCmd = Cmd(ctx)
	return
}

type MockTestingTForMain struct {
	testing.TB
}

func (MockTestingTForMain) Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
}

func (MockTestingTForMain) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
	debug.PrintStack()
}

func (MockTestingTForMain) FailNow() {
	os.Exit(1)
}

func (MockTestingTForMain) TempDir() string {
	dirname := filepath.Join(os.TempDir(), randx.MustString(6, randx.AlphaLowerNum))
	if err := os.MkdirAll(dirname, 0700); err != nil {
		panic(err)
	}
	return dirname
}

func (MockTestingTForMain) Helper() {}

func (MockTestingTForMain) Name() string {
	return "TestMain"
}
