// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"context"
	"fmt"
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	ctx context.Context
	// defaultProject and extraProject are read-only fixtures shared by every
	// test in this package. Tests that write project configuration must not use
	// them, because they run in parallel and would clobber each other; call
	// newProject instead.
	defaultProject, extraProject      *cloud.Project
	defaultConfig, defaultWorkspaceID string
	defaultCmd                        *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	ctx, defaultConfig, defaultWorkspaceID, extraProject, defaultProject, defaultCmd = testhelpers.CreateDefaultAssetsBrowser()
	fmt.Println("done setting up, running tests")
	m.Run()
}
