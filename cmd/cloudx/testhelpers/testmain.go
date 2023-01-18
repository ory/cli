// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

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
	setEnvIfUnset("ORY_CLOUD_CONSOLE_URL", "https://console.staging.ory.dev:443")
	setEnvIfUnset("ORY_CLOUD_ORYAPIS_URL", "https://staging.oryapis.dev:443")
}

func CreateDefaultAssets() (defaultConfig, defaultEmail, defaultPassword, extraProject, defaultProject string, defaultCmd *cmdx.CommandExecuter) {
	UseStaging()

	t := testingT{}

	defaultConfig = NewConfigDir(t)

	defaultEmail, defaultPassword = RegisterAccount(t, defaultConfig)
	extraProject = CreateProject(t, defaultConfig)
	defaultProject = CreateProject(t, defaultConfig)
	defaultCmd = ConfigAwareCmd(defaultConfig)
	return
}

func RunAgainstStaging(m *testing.M) {
	UseStaging()
	os.Exit(m.Run())
}

type testingT struct{}

func (testingT) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
	debug.PrintStack()
}

func (testingT) FailNow() {
	os.Exit(1)
}
