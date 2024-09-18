// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package relationtuples_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	cloud "github.com/ory/client-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/relationtuples"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	defaultProject *cloud.Project
	defaultCmd     *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	_, _, _, _, defaultProject, defaultCmd = testhelpers.CreateDefaultAssetsBrowser()
	m.Run()
}

func TestNoUnauthenticated(t *testing.T) {
	t.Parallel()
	cases := []struct {
		verb      string
		extraArgs []string
	}{
		{verb: "create", extraArgs: []string{"-"}},
		{verb: "ls"},
		{verb: "list"},
		{verb: "delete", extraArgs: []string{"--all"}},
	}

	for _, tc := range cases {
		t.Run("verb="+tc.verb, func(t *testing.T) {
			cmd := testhelpers.Cmd(testhelpers.WithCleanConfigFile(context.Background(), t))
			args := append([]string{tc.verb, "relationships", "--quiet", "--project", defaultProject.Id},
				tc.extraArgs...)
			_, _, err := cmd.Exec(nil, args...)
			require.ErrorIsf(t, err, client.ErrNoConfigQuiet, "got error: %v", err)
		})
	}
}

func TestAfterAuthentication(t *testing.T) {
	t.Parallel()
	cases := []struct {
		verb      string
		extraArgs []string
	}{
		{verb: "ls"},
		{verb: "list"},
		{verb: "delete", extraArgs: []string{"--all"}},
	}

	for _, tc := range cases {
		t.Run("verb="+tc.verb, func(t *testing.T) {
			ctx := testhelpers.WithEmitAuthFlowTriggeredErr(context.Background(), t)
			args := append([]string{tc.verb, "relation-tuples", "--project", defaultProject.Id},
				tc.extraArgs...)
			_, stderr, err := testhelpers.Cmd(ctx).Exec(nil, args...)
			require.ErrorIs(t, err, testhelpers.ErrAuthFlowTriggered, stderr)
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	in := strings.NewReader("nspace:obj#rel@sub")
	out, stderr, err := defaultCmd.Exec(in, "parse", "relation-tuples", "--format=json", "--project", defaultProject.Id, "-")

	require.NoError(t, err, stderr)
	assert.JSONEq(t, `{"namespace":"nspace","object":"obj","relation":"rel","subject_id":"sub"}`, out, out)
}

// TestCRUD tests that we can create, list, and delete tuples.
// The tests must be run in a sequence.
func TestCRUD(t *testing.T) {
	t.Parallel()
	createLegacyNamespace(t, defaultProject.Id, `{"name": "n", "id": 0}`)

	tuple := func(object string) string {
		return fmt.Sprintf(`[{
	"namespace": "n",
	"object": %q,
	"relation": "r",
	"subject_id": "s"
}]`, object)
	}
	create := func(t *testing.T, object string) string {
		in := strings.NewReader(tuple(object))
		stdout, stderr, err := defaultCmd.Exec(in, "create", "relation-tuples", "--format", "json", "--project", defaultProject.Id, "-")
		require.NoError(t, err, stderr)
		return stdout
	}
	list := func(t *testing.T) string {
		stdout, stderr, err := defaultCmd.Exec(nil, "list", "relation-tuples", "--format", "json", "--project", defaultProject.Id)
		require.NoError(t, err, stderr)
		return stdout
	}
	isAllowed := func(t *testing.T, subject, relation, namespace, object string) string {
		stdout, stderr, err := defaultCmd.Exec(nil,
			"is", "allowed", subject, relation, namespace, object,
			"--project", defaultProject.Id, "--format", "json")
		require.NoError(t, err, stderr)
		return stdout
	}

	// 1. create a tuple
	stdout := create(t, "o1")
	require.JSONEq(t, tuple("o1"), stdout)

	// 2. check that it is listed
	stdout = list(t)
	require.JSONEq(t, tuple("o1"), gjson.Get(stdout, "relation_tuples").Raw, stdout)

	// check that it is allowed
	stdout = isAllowed(t, "s", "r", "n", "o1")
	require.JSONEq(t, `{"allowed":true}`, stdout, stdout)

	// 3. delete with --all but without --force
	stdout, stderr, err := defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", defaultProject.Id,
		"--all")
	require.NoError(t, err, stderr)
	require.JSONEq(t, tuple("o1"), gjson.Get(stdout, "relation_tuples").Raw, stdout)

	// 4. create a second tuple
	create(t, "o2")

	// 5. delete without --all but with --force
	_, stderr, err = defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", defaultProject.Id,
		"--force")
	require.ErrorIs(t, err, relationtuples.ErrDeleteMissingAllFlag, stderr)
	require.Len(t, gjson.Get(list(t), "relation_tuples").Array(), 2, list(t))

	// 6. delete one tuple with query and with --force
	_, stderr, err = defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", defaultProject.Id,
		"--force", "--object=o2")
	require.NoError(t, err, stderr)
	require.JSONEq(t, tuple("o1"), gjson.Get(list(t), "relation_tuples").Raw, list(t))

	// 7. create another tuple (now two on the server)
	create(t, "o42")

	// 8. delete with --all and with --force
	_, stderr, err = defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", defaultProject.Id,
		"--force", "--all")
	require.NoError(t, err, stderr)
	assert.Len(t, gjson.Get(list(t), "relation_tuples").Array(), 0, list(t))
}

func createLegacyNamespace(t *testing.T, project, rawNamespace string) {
	t.Helper()
	_, _, err := defaultCmd.Exec(nil, "patch", "permission-config", "--project", project,
		"--add", `/namespaces/-=`+rawNamespace)
	if err != nil {
		t.Fatal(err)
	}
}
