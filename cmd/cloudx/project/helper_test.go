// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

type execFunc = func(stdin io.Reader, args ...string) (string, string, error)

// newProject creates a project for the exclusive use of the calling test.
//
// Tests that write project configuration must not share a project. They run in
// parallel, and several of them write the same keys with different values: for
// example TestPatchProject replaces
// /services/identity/config/selfservice/flows/error with an object, while
// TestPatchKratosConfig sets /selfservice/flows/error/ui_url — the same key
// underneath. Sharing a project makes them clobber each other and fail
// intermittently in CI.
//
// defaultProject and extraProject are therefore read-only fixtures: only tests
// that do not modify project configuration may use them.
//
// Subtests of one test function do share the project this returns, so each must
// still perform a single operation and assert on the response to that operation
// rather than re-reading the project.
//
// The project goes into a workspace of its own because the development-project
// quota is per workspace: adding these to defaultWorkspaceID fails with
// "the quota for the feature 'Development Projects' has been exceeded". This is
// the same reason TestUpdateProject and TestListProject create a workspace
// before their projects.
func newProject(t *testing.T) string {
	t.Helper()
	return testhelpers.CreateProject(ctx, t, testhelpers.CreateWorkspace(ctx, t)).Id
}

func runWithProjectAsDefault(ctx context.Context, t *testing.T, projectID string, test func(t *testing.T, exec execFunc)) {
	t.Run("project passed as default", func(t *testing.T) {
		ctx, _ := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		testhelpers.SetDefaultProject(ctx, t, projectID)

		test(t, testhelpers.Cmd(ctx).Exec)

		// make sure, the default wasn't changed implicitly
		assert.Equal(t, projectID, testhelpers.GetDefaultProjectID(ctx, t))
	})
}

func runWithProjectAsArgument(ctx context.Context, t *testing.T, projectID string, test func(t *testing.T, exec execFunc)) {
	t.Run("project passed as argument", func(t *testing.T) {
		ctx, _ := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		selectedProject := testhelpers.GetDefaultProjectID(ctx, t)
		require.NotEqual(t, selectedProject, projectID, "to ensure correct isolation, please use another project than the 'default' selected")

		cmd := testhelpers.Cmd(ctx)
		test(t, func(stdin io.Reader, args ...string) (string, string, error) {
			return cmd.Exec(stdin, append(args, projectID)...)
		})

		// make sure, the default wasn't changed implicitly
		assert.Equal(t, selectedProject, testhelpers.GetDefaultProjectID(ctx, t))
	})
}

func runWithProjectAsFlag(ctx context.Context, t *testing.T, projectID string, test func(t *testing.T, exec execFunc)) {
	t.Run("project passed as flag", func(t *testing.T) {
		ctx, _ := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		selectedProject := testhelpers.GetDefaultProjectID(ctx, t)
		require.NotEqual(t, selectedProject, projectID, "to ensure correct isolation, please use another project than the 'default' selected")

		cmd := testhelpers.Cmd(ctx)
		test(t, func(stdin io.Reader, args ...string) (string, string, error) {
			return cmd.Exec(stdin, append(args, "--project", projectID)...)
		})

		// make sure, the default wasn't changed implicitly
		assert.Equal(t, selectedProject, testhelpers.GetDefaultProjectID(ctx, t))
	})
}
