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

func runWithProjectAsDefault(ctx context.Context, t *testing.T, projectID string, test func(t *testing.T, exec execFunc)) {
	t.Run("project passed as default", func(t *testing.T) {
		ctx := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
		testhelpers.SetDefaultProject(ctx, t, projectID)

		test(t, testhelpers.Cmd(ctx).Exec)

		// make sure, the default wasn't changed implicitly
		assert.Equal(t, projectID, testhelpers.GetDefaultProjectID(ctx, t))
	})
}

func runWithProjectAsArgument(ctx context.Context, t *testing.T, projectID string, test func(t *testing.T, exec execFunc)) {
	t.Run("project passed as argument", func(t *testing.T) {
		ctx := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
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
		ctx := testhelpers.WithDuplicatedConfigFile(ctx, t, defaultConfig)
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
