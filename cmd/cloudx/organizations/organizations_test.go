// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package organizations_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	project, defaultEmail, defaultPassword string
	defaultCmd                             *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	_, defaultEmail, defaultPassword, _, project, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
}

func TestNoUnauthenticated(t *testing.T) {
	t.Parallel()
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
			configDir := testhelpers.NewConfigDir(t)
			cmd := testhelpers.ConfigAwareCmd(configDir)
			args := []string{tc.verb, tc.noun, "--quiet", "--project", project}
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

	defaultCmd.ExecNoErr(t, "use", project)

	// List organizations: Empty
	out := defaultCmd.ExecNoErr(t, "list", "organizations", "--format=json")
	assert.Equal(t, "[]\n", out)
}
