// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"bytes"
	"testing"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
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
		testhelpers.SetDefaultProject(t, defaultConfig, defaultProject)

		name := testhelpers.TestProjectName()
		stdout, _, err := defaultCmd.Exec(nil, "create", "project", "--name", name, "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)

		assert.Equal(t, defaultProject, testhelpers.GetDefaultProject(t, defaultConfig))
	})

	t.Run("is able to create a project and use the project as default", func(t *testing.T) {
		name := testhelpers.TestProjectName()

		stdout, _, err := defaultCmd.Exec(nil, "create", "project", "--name", name, "--use-project", "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)

		id, _, _ := parseOutput(stdout)
		assert.Equal(t, id, testhelpers.GetDefaultProject(t, defaultConfig))
	})

	t.Run("is able to create a project and use name from stdin", func(t *testing.T) {
		name := testhelpers.TestProjectName()
		stdin := bytes.NewBufferString(name + "\n")
		stdout, _, err := defaultCmd.Exec(stdin, "create", "project", "--format", "json")
		require.NoError(t, err)
		assertResult(t, defaultConfig, stdout, name)
	})

	t.Run("is not able to create a project if no name flag and quiet flag", func(t *testing.T) {
		name := testhelpers.TestProjectName()
		stdin := bytes.NewBufferString(name)
		_, stderr, err := defaultCmd.Exec(stdin, "create", "project", "--quiet")
		require.Error(t, err)
		assert.Contains(t, stderr, "you must specify the --name flag when using --quiet")
	})

	t.Run("is not able to create a project if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		name := testhelpers.TestProjectName()
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "create", "project", "--name", name, "--format", "json", "--quiet")
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to create a project after authenticating", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		name := testhelpers.TestProjectName()
		password := testhelpers.FakePassword()
		email := testhelpers.FakeEmail()
		cmd := testhelpers.ConfigPasswordAwareCmd(configDir, password)
		// Create the account
		var r bytes.Buffer
		r.WriteString("n\n")        // Do you want to sign in to an existing Ory Network account? [y/n]: n
		r.WriteString(email + "\n") // Email fakeEmail()
		r.WriteString(name + "\n")  // Name: fakeName()
		r.WriteString("n\n")        // Subscribe to the Ory Security Newsletter to get platform and security updates? [y/n]: n
		r.WriteString("n\n")        // I accept the Terms of Service [y/n]: n
		r.WriteString("y\n")        // I accept the Terms of Service [y/n]: y
		stdout, stderr, err := cmd.Exec(&r, "create", "project", "--name", name, "--format", "json")
		require.NoError(t, err)
		t.Logf("stdout: %s", stdout)
		t.Logf("stderr: %s", stderr)
		assertResult(t, configDir, stdout, name)
	})
}
