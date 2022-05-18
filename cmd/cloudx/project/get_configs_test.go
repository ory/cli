package project_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetServiceConfig(t *testing.T) {
	t.Run("service=kratos", func(t *testing.T) {
		stdout, _, err := defaultCmd.Exec(nil, "get", "kratos-config", defaultProject, "--format", "json")
		require.NoError(t, err)
		assert.True(t, gjson.Get(stdout, "selfservice.flows.error.ui_url").Exists())
	})

	t.Run("service=keto", func(t *testing.T) {
		stdout, _, err := defaultCmd.Exec(nil, "get", "keto-config", defaultProject, "--format", "json")
		require.NoError(t, err)
		assert.True(t, gjson.Get(stdout, "namespaces").Exists(), stdout)
	})
}
