// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	defaultConfig, defaultEmail, defaultPassword string
	defaultProject                               *cloud.Project
	defaultCmd                                   *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	defaultConfig, defaultEmail, defaultPassword, _, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
}
