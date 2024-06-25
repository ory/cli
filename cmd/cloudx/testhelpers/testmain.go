// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"
	"github.com/ory/x/randx"

	"github.com/ory/x/cmdx"
)

func setEnvIfUnset(key, value string) {
	if _, ok := os.LookupEnv(key); !ok {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}

func UseStaging() {
	setEnvIfUnset(client.ConsoleURLKey, "https://console.staging.ory.dev:443")
	setEnvIfUnset(client.OryAPIsURLKey, "https://staging.oryapis.dev:443")
}

func CreateDefaultAssets() (defaultConfig, defaultEmail, defaultPassword string, extraProject, defaultProject *cloud.Project, defaultCmd *cmdx.CommandExecuter) {
	UseStaging()

	t := testingT{}

	defaultConfig = NewConfigFile(t)

	defaultEmail, defaultPassword, _ = RegisterAccount(t, defaultConfig)
	extraProject = CreateProject(t, defaultConfig, nil)
	defaultProject = CreateProject(t, defaultConfig, nil)
	defaultCmd = CmdWithConfig(defaultConfig)
	return
}

func RunAgainstStaging(m *testing.M) {
	UseStaging()
	os.Exit(m.Run())
}

type testingT struct {
	testing.TB
}

func (testingT) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
	debug.PrintStack()
}

func (testingT) FailNow() {
	os.Exit(1)
}

func (testingT) TempDir() string {
	dirname := filepath.Join(os.TempDir(), randx.MustString(6, randx.AlphaLowerNum))
	if err := os.MkdirAll(dirname, 0700); err != nil {
		panic(err)
	}
	return dirname
}

func (testingT) Helper() {}

func (testingT) Name() string {
	return "TestMain"
}
