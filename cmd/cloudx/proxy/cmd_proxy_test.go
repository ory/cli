// Copyright © 2022 Ory Corp
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
		_, err := getEndpointURL(newEndpointCmd(""))
		require.Error(t, err)
	})

	t.Run("should return the right value from the flag", func(t *testing.T) {
		expected := "someslug"
		cmd := newEndpointCmd(expected)
		actual, err := getEndpointURL(cmd)
		require.NoError(t, err)
		assert.Equal(t, "https://"+expected+".projects.oryapis.com/", actual.String())
	})

	t.Run("should return the right value from the OS", func(t *testing.T) {
		var b bytes.Buffer
		expected := "someslug"
		t.Setenv(envVarSlug, expected)
		cmd := newEndpointCmd("not-someslug")
		cmd.SetErr(&b)
		actual, err := getEndpointURL(cmd)
		require.NoError(t, err)
		assert.Equal(t, "https://"+expected+".projects.oryapis.com/", actual.String())
		assert.Equal(t, "Attention! We found multiple sources for the project slug. Please clean up environment variables and flags to ensure that the correct value is being used. Found values:\n\n\t--project=not-someslug\n\tORY_PROJECT_SLUG=someslug\n\nOrder of precedence is: ORY_SDK_URL > ORY_KRATOS_URL > ORY_PROJECT_SLUG > --project\nDecided to use value: https://someslug.projects.oryapis.com/\n\n", b.String())
	})

	t.Run("should return the right value from the OS using a legacy value", func(t *testing.T) {
		var b bytes.Buffer
		expected := "https://someslug.projects.oryapis.com/"
		t.Setenv(envVarSDK, expected)
		cmd := newEndpointCmd("not-someslug")
		cmd.SetErr(&b)
		actual, err := getEndpointURL(cmd)
		require.NoError(t, err)
		assert.Equal(t, expected, actual.String())
		assert.Equal(t, "It is recommended to use the --project flag or the ORY_PROJECT_SLUG environment variable for better developer experience. Environment variables ORY_SDK_URL and ORY_KRATOS_URL will continue to work!\nAttention! We found multiple sources for the project slug. Please clean up environment variables and flags to ensure that the correct value is being used. Found values:\n\n\t--project=not-someslug\n\tORY_SDK_URL=https://someslug.projects.oryapis.com/\n\nOrder of precedence is: ORY_PROJECT_SLUG > ORY_SDK_URL > ORY_KRATOS_URL > --project\nDecided to use value: https://someslug.projects.oryapis.com/\n\n", b.String())
	})

	t.Run("should fail if legacy value is not a URL", func(t *testing.T) {
		var b bytes.Buffer
		expected := "not-a-url"
		t.Setenv(envVarSDK, expected)
		cmd := newEndpointCmd("not-someslug")
		cmd.SetErr(&b)
		_, err := getEndpointURL(cmd)
		require.Error(t, err)
	})
}
