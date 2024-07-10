// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

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
	if testing.Short() {
		// this test needs internet, typically not available when you're on a (german) train
		return
	}

	t.Parallel()

	content := `class Default implements Namespace {}`
	config := writeFile(t, content)
	verbs := []string{"update", "patch"}

	for _, verb := range verbs {
		t.Run(fmt.Sprintf("is able to %q the namespace config", verb), func(t *testing.T) {
			t.Parallel()

			updateNamespace := func(t *testing.T, exec execFunc) {
				stdout, stderr, err := exec(nil, verb, "opl", "--format", "json", "--file", config)
				require.NoError(t, err, stderr)

				url := gjson.Get(stdout, "namespaces.location").String()
				data, err := fetcher.NewFetcher().Fetch(url)
				require.NoError(t, err, "could not download the config")
				assert.Equal(t, content, data.String(), "the downloaded file does not match what we uploaded")
			}

			runWithProjectAsDefault(ctx, t, defaultProject.Id, updateNamespace)
			runWithProjectAsFlag(ctx, t, extraProject.Id, updateNamespace)
		})
	}
}
