// Copyright © 2022 Ory Corp

package project_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/fetcher"
)

func writeFile(t *testing.T, content string) (path string) {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "keto-namespaces-*.ts")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}

	return f.Name()
}

func TestUpdateNamespaceConfig(t *testing.T) {
	project := testhelpers.CreateProject(t, defaultConfig)
	content := `class Default implements Namespace {}`
	config := writeFile(t, content)
	verbs := []string{"update", "patch"}

	for _, verb := range verbs {
		t.Run(fmt.Sprintf("is able to %q the namespace config", verb), func(t *testing.T) {
			stdout, stderr, err := defaultCmd.Exec(nil, verb, "opl", "--project", project, "--format", "json", "--file", config)
			require.NoError(t, err, stderr)

			if !testing.Short() {
				// Don't download and compare the config in short mode, might not have internet everywhere
				url := gjson.Get(stdout, "namespaces.location").String()
				data, err := fetcher.NewFetcher().Fetch(url)
				require.NoError(t, err, "could not download the config")
				assert.Equal(t, content, data.String(), "the downloaded file does not match what we uploaded")
			}
		})
	}
}
