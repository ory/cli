// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package organizations_test

import (
	"context"
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	defaultProject *cloud.Project
	defaultCmd     *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	_, _, _, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
}

func TestNoUnauthenticated(t *testing.T) {
	t.Parallel()

	ctx := testhelpers.WithCleanConfigFile(context.Background(), t)
	cmd := testhelpers.Cmd(ctx)

	cases := []struct {
		verb string
		noun string
		arg  string
	}{
		{verb: "create", noun: "organization", arg: "my-org"},
		{verb: "ls", noun: "organizations"},
		{verb: "list", noun: "organizations"},
		{verb: "delete", noun: "organization", arg: "some-uuid"},
	}

	for _, tc := range cases {
		t.Run("verb="+tc.verb, func(t *testing.T) {
			args := []string{tc.verb, tc.noun, "--quiet", "--project", defaultProject.Id}
			if tc.arg != "" {
				args = append(args, tc.arg)
			}
			_, _, err := cmd.Exec(nil, args...)
			require.ErrorIsf(t, err, client.ErrNoConfigQuiet, "got error: %v", err)
		})
	}
}

func TestCRUD(t *testing.T) {
	t.Parallel()

	// List organizations: Empty
	out := defaultCmd.ExecNoErr(t, "list", "organizations", "--format=json", "--project", defaultProject.Id)
	assert.Equal(t, "[]\n", out)
}
