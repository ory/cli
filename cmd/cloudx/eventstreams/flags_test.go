// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package eventstreams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr(s string) *string { return &s }

func TestStreamConfigValidate(t *testing.T) {
	t.Parallel()

	base := func() streamConfig {
		return streamConfig{
			Type:          "https",
			HttpsEndpoint: ptr("https://example.com/webhook"),
		}
	}

	t.Run("accepts a valid active status", func(t *testing.T) {
		c := base()
		c.Status = ptr(StatusActive)
		require.NoError(t, c.Validate())
		assert.Equal(t, StatusActive, *c.Status)
	})

	t.Run("accepts a valid paused status", func(t *testing.T) {
		c := base()
		c.Status = ptr(StatusPaused)
		require.NoError(t, c.Validate())
		assert.Equal(t, StatusPaused, *c.Status)
	})

	t.Run("normalizes an empty status to nil so the server default applies", func(t *testing.T) {
		c := base()
		c.Status = ptr("")
		require.NoError(t, c.Validate())
		assert.Nil(t, c.Status)
	})

	t.Run("rejects an unknown status", func(t *testing.T) {
		c := base()
		c.Status = ptr("frozen")
		assert.ErrorContains(t, c.Validate(), "--status")
	})

	t.Run("status is optional when unset", func(t *testing.T) {
		c := base()
		require.NoError(t, c.Validate())
		assert.Nil(t, c.Status)
	})
}

func TestStreamConfigToSetBody(t *testing.T) {
	t.Parallel()

	t.Run("maps all fields including the required type as a pointer", func(t *testing.T) {
		c := streamConfig{
			Type:          "https",
			HttpsEndpoint: ptr("https://example.com/webhook"),
			Status:        ptr(StatusPaused),
		}
		body := c.toSetBody()
		require.NotNil(t, body.Type)
		assert.Equal(t, "https", *body.Type)
		require.NotNil(t, body.HttpsEndpoint)
		assert.Equal(t, "https://example.com/webhook", *body.HttpsEndpoint)
		require.NotNil(t, body.Status)
		assert.Equal(t, StatusPaused, *body.Status)
	})

	t.Run("leaves type nil when unset so the current type is kept", func(t *testing.T) {
		c := streamConfig{Status: ptr(StatusActive)}
		body := c.toSetBody()
		assert.Nil(t, body.Type)
		require.NotNil(t, body.Status)
		assert.Equal(t, StatusActive, *body.Status)
	})
}
