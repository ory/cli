// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	defaultProject, extraProject *cloud.Project
	defaultConfig                string
	defaultCmd                   *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	defaultConfig, _, _, extraProject, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
}
