package monorepo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonorepoInvalidFile(t *testing.T) {
	var graph ComponentGraph
	_, err := graph.getComponentGraph("test/invalidConfigFile")
	//todo check for specific error
	assert.Error(t, err)
}

func TestMonorepoCirclular(t *testing.T) {
	var graph ComponentGraph
	g, err := graph.getComponentGraph("test/circular")

	//successfully read configurations
	assert.Nil(t, err, "Successfully read configuration")
	_, err = g.resolveGraph()
	assert.NotNil(t, err, "Failed resolving graph because of circular dependency")
}

func TestMonorepoWorking(t *testing.T) {
	var graph ComponentGraph

	graph.getComponentGraph("test/working")
	graph.displayGraph()

	resolved, err := graph.resolveGraph()
	assert.Nil(t, err)
	assert.Equal(t, 6, resolved.len())
	/*
		1. Invalid Root Directory
		2. No Config Files
		3. Invalid Config Files
		4. Valid Graph
			a) Parse
			b) Resolved Graph
		5. Trigger
		5. Circular Dependencies

	*/
}
