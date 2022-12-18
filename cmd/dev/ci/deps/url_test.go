// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package deps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const FullWorking = `version: v1.20.2
url: https://storage.googleapis.com/kubernetes-release/release/{{.Version}}/bin/{{.Os}}/{{.Architecture}}/kubectl
mappings:
  architecture:
    amd64: x64
    arm64: aarch64
  os:
    darwin: mac
    linux: unix
`

func TestValidFile(t *testing.T) {
	comp1 := Component{}
	path := "test/full-working.yaml"
	err := comp1.getComponentFromConfig(path)
	require.NoError(t, err)
	assert.Equal(t, FullWorking, comp1.String())
}

func TestInvalidFile(t *testing.T) {
	comp1 := Component{}
	err := comp1.getComponentFromConfig("test/invalidFile.yaml")
	var ifError InvalidFileError
	assert.ErrorAs(t, err, &ifError, "Wrong Error Type")
}

func TestFileNotFound(t *testing.T) {
	path := "test/this-does-not-exist.yaml"
	comp1 := Component{}
	err := comp1.getComponentFromConfig(path)
	var fnfError FileNotFoundError
	assert.ErrorAs(t, err, &fnfError, "Wrong Error Type")
}

func TestDefaultURL(t *testing.T) {
	var defaultURL = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/darwin/amd64/kubectl`

	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/defaultURL.yaml")
	url, err := comp1.getRenderedURL("darwin", "amd64")
	require.NoError(t, err)
	assert.Equal(t, defaultURL, url)
}

func TestCustomArchURL(t *testing.T) {
	var customArchURL1 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/darwin/x64/kubectl`
	var customArchURL2 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/darwin/aarch64/kubectl`

	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/customArchURL.yaml")

	url, err := comp1.getRenderedURL("darwin", "amd64")
	require.NoError(t, err)
	assert.Equal(t, customArchURL1, url)

	url, err = comp1.getRenderedURL("darwin", "arm64")
	require.NoError(t, err)
	assert.Equal(t, customArchURL2, url)
}

func TestCustomOSURL(t *testing.T) {
	var customOSURL1 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/mac/amd64/kubectl`
	var customOSURL2 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/unix/amd64/kubectl`

	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/customOSURL.yaml")
	url, err := comp1.getRenderedURL("darwin", "amd64")
	require.NoError(t, err)
	assert.Equal(t, customOSURL1, url)
	url, err = comp1.getRenderedURL("linux", "amd64")
	require.NoError(t, err)
	assert.Equal(t, customOSURL2, url)
}
