// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	extraProject, defaultProject, defaultConfig, defaultEmail, defaultPassword string
	defaultCmd                                                                 *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	defaultConfig, defaultEmail, defaultPassword, extraProject, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
}
