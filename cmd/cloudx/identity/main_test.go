// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity_test

import (
	"context"
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	ctx            context.Context
	defaultProject *cloud.Project
	defaultCmd     *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	ctx, _, _, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()
	m.Run()
}
