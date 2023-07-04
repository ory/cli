// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy_test

import (
	"net/url"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/proxy"
	cloud "github.com/ory/client-go"
)

func TestProxyUseProjectID(t *testing.T) {
	projectId, err := uuid.NewV4()
	p := &cloud.Project{
		Id:   projectId.String(),
		Slug: "test",
		Name: "test",
	}

	var noopCommand client.Command = &client.MockCommandHelper{
		Project: p,
	}

	projectURL := func(t *testing.T, slugOrId string) *url.URL {
		t.Helper()
		projectURL, err := url.Parse("https://" + slugOrId + ".projects.oryapis.com/")
		require.NoError(t, err)
		return projectURL
	}

	t.Run("case=should be able to use project id instead of a slug", func(t *testing.T) {
		conf := &proxy.ProxyConfig{
			Upstream:      projectURL(t, projectId.String()).String(),
			OryURL:        projectURL(t, projectId.String()),
			ProjectSlugId: projectId.String(),
		}

		conf, err = proxy.UseProjectIdOrSlug(noopCommand, conf, "test")
		require.NoError(t, err)
		require.Equal(t, projectId.String(), conf.ProjectSlugId)
		require.Equal(t, projectURL(t, p.Slug).String(), conf.Upstream)
		require.Equal(t, projectURL(t, p.Slug).String(), conf.OryURL.String())
	})

	t.Run("case=should be able to still use project slug", func(t *testing.T) {
		conf := &proxy.ProxyConfig{
			Upstream:      projectURL(t, projectId.String()).String(),
			OryURL:        projectURL(t, projectId.String()),
			ProjectSlugId: projectId.String(),
		}

		conf, err = proxy.UseProjectIdOrSlug(noopCommand, conf, "test")
		require.NoError(t, err)
		require.Equal(t, projectId.String(), conf.ProjectSlugId)
		require.Equal(t, "https://"+p.Slug+".projects.oryapis.com/", conf.Upstream)
		require.Equal(t, "https://"+p.Slug+".projects.oryapis.com/", conf.OryURL.String())
	})

	t.Run("case=should get an error if no api key is set while using project id", func(t *testing.T) {
		conf := &proxy.ProxyConfig{
			Upstream:      projectURL(t, projectId.String()).String(),
			OryURL:        projectURL(t, projectId.String()),
			ProjectSlugId: projectId.String(),
		}

		conf, err = proxy.UseProjectIdOrSlug(noopCommand, conf, "")
		require.ErrorContains(t, err, "A project ID was provided instead of a project slug, but no API key was found.")
	})

	t.Run("case=should be able to still use project slug with no api key", func(t *testing.T) {
		slug := "test-project"
		conf := &proxy.ProxyConfig{
			Upstream:      projectURL(t, slug).String(),
			OryURL:        projectURL(t, slug),
			ProjectSlugId: slug,
		}
		conf, err = proxy.UseProjectIdOrSlug(noopCommand, conf, "")
		require.NoError(t, err)
		require.Equal(t, slug, conf.ProjectSlugId)
		require.Equal(t, projectURL(t, slug).String(), conf.Upstream)
		require.Equal(t, projectURL(t, slug).String(), conf.OryURL.String())
	})
}
