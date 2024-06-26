// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestCreateProject(t *testing.T) {
	parseOutput := func(stdout string) (id string, slug string, name string) {
		id = gjson.Get(stdout, "id").String()
		slug = gjson.Get(stdout, "slug").String()
		name = gjson.Get(stdout, "name").String()
		return
	}
	assertResult := func(t *testing.T, configDir string, stdout string, expectedName string) {
		_, slug, name := parseOutput(stdout)
		assert.NotEmpty(t, slug, stdout)
		assert.Equal(t, expectedName, name, stdout)
	}

	t.Run("is able to create a project", func(t *testing.T) {
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject.Id)

		name := testhelpers.TestName()
		stdout, _, err := defaultCmd.Exec(nil, "create", "project", "--name", name, "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)

		assert.Equal(t, defaultProject.Id, testhelpers.GetDefaultProjectID(t, defaultConfig))
	})

	t.Run("is able to create a project and use the project as default", func(t *testing.T) {
		name := testhelpers.TestName()

		stdout, _, err := defaultCmd.Exec(nil, "create", "project", "--name", name, "--use-project", "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)

		id, _, _ := parseOutput(stdout)
		assert.Equal(t, id, testhelpers.GetDefaultProjectID(t, defaultConfig))
	})

	t.Run("is able to create a project and use name from stdin", func(t *testing.T) {
		name := testhelpers.TestName()
		stdin := bytes.NewBufferString(name + "\n")
		stdout, _, err := defaultCmd.Exec(stdin, "create", "project", "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)
	})

	t.Run("is not able to create a project if no name flag and quiet flag", func(t *testing.T) {
		name := testhelpers.TestName()
		stdin := bytes.NewBufferString(name)
		_, stderr, err := defaultCmd.Exec(stdin, "create", "project", "--quiet")
		require.Error(t, err)
		assert.Contains(t, stderr, "you must specify the --name and --environment flags when using --quiet")
	})

	t.Run("is not able to create a project if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigFile(t)
		name := testhelpers.TestName()
		cmd := testhelpers.CmdWithConfig(configDir)
		_, _, err := cmd.Exec(nil, "create", "project", "--name", name, "--format", "json", "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to create a project after authenticating", func(t *testing.T) {
		configDir := testhelpers.NewConfigFile(t)
		name := testhelpers.TestName()
		password := testhelpers.FakePassword()
		email := testhelpers.FakeEmail()

		cmd := testhelpers.CmdWithConfigPassword(configDir, password)
		// Create the account
		r := testhelpers.RegistrationBuffer(name, email)
		stdout, stderr, err := cmd.Exec(r, "create", "project", "--name", name, "--format", "json")
		t.Logf("stdout: %s", stdout)
		t.Logf("stderr: %s", stderr)
		require.NoError(t, err)
		assertResult(t, configDir, stdout, name)
	})
}
