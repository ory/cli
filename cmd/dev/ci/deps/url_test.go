package deps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const FullWorking = `version: v1.20.2
url: https://storage.googleapis.com/kubernetes-release/release/{{.Version}}/bin/{{.Os}}/{{.Architecture}}/kubectl
mappings:
  architecture:
    amd64: x64
  os:
    darwin: mac
    linux: unix
`

func TestValidFile(t *testing.T) {
	comp1 := Component{}
	path := "test/full-working.yaml"
	err := comp1.getComponentFromConfig(path)
	assert.Nil(t, err)
	assert.Equal(t, FullWorking, comp1.String())
}

func TestInvalidFile(t *testing.T) {
	comp1 := Component{}
	err := comp1.getComponentFromConfig("test/invalidFile.yaml")
	var ifError InvalidFileError
	assert.NotNil(t, err, "Excepted Error!")
	assert.ErrorAs(t, err, &ifError, "Wrong Error Type")
}

func TestFileNotFound(t *testing.T) {
	path := "test/this-does-not-exist.yaml"
	comp1 := Component{}
	err := comp1.getComponentFromConfig(path)
	var fnfError FileNotFoundError
	assert.NotNil(t, err, "Excepted Error!")
	assert.ErrorAs(t, err, &fnfError, "Wrong Error Type")
}

func TestDefaultURL(t *testing.T) {
	var defaultURL = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/darwin/amd64/kubectl`

	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/defaultURL.yaml")
	url, err := comp1.getRenderedURL("darwin", "amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, defaultURL, url)
}

func TestCustomArchURL(t *testing.T) {
	var customArchURL = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/darwin/x64/kubectl`

	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/customArchURL.yaml")
	url, err := comp1.getRenderedURL("darwin", "amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, customArchURL, url)
}

func TestCustomOSURL(t *testing.T) {
	var customOSURL1 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/mac/amd64/kubectl`
	var customOSURL2 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/unix/amd64/kubectl`

	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/customOSURL.yaml")
	url, err := comp1.getRenderedURL("darwin", "amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, customOSURL1, url)
	url, err = comp1.getRenderedURL("linux", "amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, customOSURL2, url)
}
