// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEndpointCmd(def string) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.ErrOrStderr()
	cmd.Flags().String(ProjectFlag, def, "")
	return cmd
}

func TestGetEndpointURL(t *testing.T) {
	t.Run("should fail if no project is set", func(t *testing.T) {
		cmd := newEndpointCmd("")
		_, err := getEndpointURL(cmd, getProjectSlugId(cmd))
		require.Error(t, err)
	})

	t.Run("should return the right value from the flag", func(t *testing.T) {
		expected := "someslug"
		cmd := newEndpointCmd(expected)
		actual, err := getEndpointURL(cmd, getProjectSlugId(cmd))
		require.NoError(t, err)
		assert.Equal(t, "https://"+expected+".projects.oryapis.com/", actual.String())
	})

	t.Run("should fail when presented with multiple endpoint configs", func(t *testing.T) {
		var b bytes.Buffer
		expected := "someslug"
		t.Setenv(envVarSlug, expected)
		cmd := newEndpointCmd("not-someslug")
		cmd.SetErr(&b)
		actual, err := getEndpointURL(cmd, getProjectSlugId(cmd))
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.Contains(t, b.String(), "Attention! We found multiple sources for the project slug. Please clean up environment variables and flags to ensure that the correct value is being used. Found values:\n\n\t--project=not-someslug\n\tORY_PROJECT_SLUG=someslug")
	})

	t.Run("should fail if legacy value is not a URL", func(t *testing.T) {
		var b bytes.Buffer
		expected := "not-a-url"
		t.Setenv(envVarSDK, expected)
		cmd := newEndpointCmd("not-someslug")
		cmd.SetErr(&b)
		_, err := getEndpointURL(cmd, getProjectSlugId(cmd))
		require.Error(t, err)
	})
}
