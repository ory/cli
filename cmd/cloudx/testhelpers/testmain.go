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

func CreateDefaultAssetsBrowser() (ctx context.Context, defaultConfig, defaultWorkspaceID string, extraProject, defaultProject *cloud.Project, defaultCmd *cmdx.CommandExecuter) {
	UseStaging()

	t := MockTestingTForMain{}

	defaultConfig = NewConfigFile(t)

	email, password, _, _ := RegisterAccount(context.Background(), t)

	_, page, cleanup := SetupPlaywright(t)
	defer cleanup()

	ctx = client.ContextWithOptions(context.Background(), client.WithConfigLocation(defaultConfig))
	h, err := client.NewCommandHelper(
		ctx,
		client.WithQuiet(false),
		client.WithOpenBrowserHook(PlaywrightAcceptConsentBrowserHook(t, page, email, password)),
	)
	require.NoError(t, err)
	require.NoError(t, h.Authenticate(ctx))
	// we don't need playwright anymore
	cleanup()
	fmt.Println("authenticated, creating default assets")

	defaultWorkspaceID = CreateWorkspace(ctx, t)
	defaultProject = CreateProject(ctx, t, defaultWorkspaceID)
	extraProject = CreateProject(ctx, t, defaultWorkspaceID)

	defaultCmd = Cmd(ctx)
	return
}

type MockTestingTForMain struct {
	testing.TB
}

func (MockTestingTForMain) Log(args ...interface{}) {
	fmt.Println(args...)
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
