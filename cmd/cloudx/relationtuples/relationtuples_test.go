package relationtuples_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/relationtuples"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	project, defaultEmail, defaultPassword string
	defaultCmd                             *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	_, defaultEmail, defaultPassword, project, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
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
			configDir := testhelpers.NewConfigDir(t)
			cmd := testhelpers.ConfigAwareCmd(configDir)
			args := append([]string{tc.verb, "relation-tuples", "--quiet", "--project", project},
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
		tc := tc
		t.Run("verb="+tc.verb, func(t *testing.T) {
			t.Parallel()
			cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
			args := append([]string{tc.verb, "relation-tuples", "--project", project},
				tc.extraArgs...)
			_, stderr, err := cmd.Exec(r, args...)
			require.NoError(t, err, stderr)
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	in := strings.NewReader("nspace:obj#rel@sub")
	out, stderr, err := defaultCmd.Exec(in, "parse", "relation-tuples", "--format=json", "--project", project, "-")

	require.NoError(t, err, stderr)
	assert.JSONEq(t, `{"namespace":"nspace","object":"obj","relation":"rel","subject_id":"sub"}`, out, out)
}

// TestCRUD tests that we can create, list, and delete tuples.
// The tests must be run in a sequence.
func TestCRUD(t *testing.T) {
	t.Parallel()
	createNamespace(t, project, `{"name": "n", "id": 0}`)

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
		stdout, stderr, err := defaultCmd.Exec(in, "create", "relation-tuples", "--format", "json", "--project", project, "-")
		require.NoError(t, err, stderr)
		return stdout
	}
	list := func(t *testing.T) string {
		stdout, stderr, err := defaultCmd.Exec(nil, "list", "relation-tuples", "--format", "json", "--project", project)
		require.NoError(t, err, stderr)
		return stdout
	}

	t.Run("step=create", func(t *testing.T) {
		stdout := create(t, "o1")
		assert.True(t, gjson.Valid(stdout))
	})

	t.Run("step=list", func(t *testing.T) {
		stdout := list(t)
		assert.JSONEq(t, tuple("o1"), gjson.Get(stdout, "relation_tuples").Raw, stdout)
	})

	t.Run("step=delete", func(t *testing.T) {
		t.Run("case=no-force", func(t *testing.T) {
			stdout, stderr, err := defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", project,
				"--all")
			require.NoError(t, err, stderr)
			assert.JSONEq(t, tuple("o1"), gjson.Get(stdout, "relation_tuples").Raw, stdout)
		})

		create(t, "o2")
		t.Run("case=force without --all", func(t *testing.T) {
			_, stderr, err := defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", project,
				"--force")
			assert.ErrorIs(t, err, relationtuples.ErrDeleteMissingAllFlag, stderr)
			assert.Len(t, gjson.Get(list(t), "relation_tuples").Array(), 2, list(t))
		})

		t.Run("case=force with a query", func(t *testing.T) {
			_, stderr, err := defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", project,
				"--force", "--object=o2")
			require.NoError(t, err, stderr)
			assert.Len(t, gjson.Get(list(t), "relation_tuples").Array(), 1, list(t))
		})

		create(t, "o42")
		t.Run("case=force with --all", func(t *testing.T) {
			_, stderr, err := defaultCmd.Exec(nil, "delete", "relation-tuples", "--format", "json", "--project", project,
				"--force", "--all")
			require.NoError(t, err, stderr)
			assert.Len(t, gjson.Get(list(t), "relation_tuples").Array(), 0, list(t))
		})
	})
}

func createNamespace(t *testing.T, project, JSON string) {
	t.Helper()
	_, _, err := defaultCmd.Exec(nil, "patch", "permission-config", project,
		"--add", `/namespaces/-=`+JSON)
	if err != nil {
		t.Fatal(err)
	}
}
