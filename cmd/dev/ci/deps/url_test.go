package deps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const FullWorking = `version: v1.20.2
url: https://storage.googleapis.com/kubernetes-release/release/{{.Version}}/bin/{{.Os}}/{{.Architecture}}/kubectl
architecture-mapping:
  amd64: x64
os-mapping:
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

var defaultUrl = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/darwin/amd64/kubectl`
func TestDefaultUrl(t *testing.T) {
	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/defaultUrl.yaml")
	url, err := comp1.getRenderedUrl("darwin","amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, defaultUrl, url)
}

var customArchUrl = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/darwin/x64/kubectl`
func TestCustomArchUrl(t *testing.T) {
	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/customArchUrl.yaml")
	url, err := comp1.getRenderedUrl("darwin","amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, customArchUrl, url)
}

var customOsUrl1 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/mac/amd64/kubectl`
var customOsUrl2 = `https://storage.googleapis.com/kubernetes-release/release/v1.20.2/bin/unix/amd64/kubectl`
func TestCustomOsUrl(t *testing.T) {
	comp1 := Component{}
	_ = comp1.getComponentFromConfig("test/customOsUrl.yaml")
	url, err := comp1.getRenderedUrl("darwin","amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, customOsUrl1, url)
	url, err = comp1.getRenderedUrl("linux","amd64")
	assert.Nil(t, err, "Expected no Error!")
	assert.Equal(t, customOsUrl2, url)
}
