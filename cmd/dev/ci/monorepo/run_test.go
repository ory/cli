package monorepo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDryRun(t *testing.T) {
	fmt.Println("TestDryRun")
	c := Component{ID: "component1", Name: "", Path: "./"}
	cmdLine := "UnknownCommand"
	err := runCmd(&c, cmdLine, true)
	assert.NoError(t, err, "Expected no error, as dryRun only prints out the command, but does not execute it!")

	err = runCmd(&c, cmdLine, false)
	assert.Error(t, err, "Expected error as comannd specified does not exist!")
}

func TestInverseRun(t *testing.T) {
	fmt.Println("TestDryRun")
	c := Component{ID: "component1", Name: "", Path: "./"}

	//not the nicest way to test, but for now the fastest. We use the error if a command gets executed to test if
	//the conditions are handled correctly. In addition we can define a test struct with inputs and outputs to reduce
	//code.
	cmdLine := "unknowncmd"
	err := runWrapper(&c, cmdLine, ModeCurrentAffected, true, false, false, false)
	assert.Error(t, err, "Expected error, as component is affected and inverseMode set to false!")
	err = runWrapper(&c, cmdLine, ModeCurrentAffected, false, false, false, true)
	assert.Error(t, err, "Expected error, as component is not affected and inverseMode set to true!")
	err = runWrapper(&c, cmdLine, ModeCurrentAffected, true, false, false, true)
	assert.NoError(t, err, "Expected no error, as component is affected and inverseMode set to true!")
	err = runWrapper(&c, cmdLine, ModeCurrentAffected, false, false, false, false)
	assert.NoError(t, err, "Expected no error, as component is not affected and inverseMode set to false!")

	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, true, false, false)
	assert.Error(t, err, "Expected error, as component has changed and inverseMode set to false!")
	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, false, false, true)
	assert.Error(t, err, "Expected error, as component has not changed and inverseMode set to true!")
	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, true, false, true)
	assert.NoError(t, err, "Expected no error, as component has changed and inverseMode set to true!")
	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, false, false, false)
	assert.NoError(t, err, "Expected no error, as component has not changed and inverseMode set to false!")

	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, true, true, false)
	assert.Error(t, err, "Expected error, as component is involved and inverseMode set to false!")
	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, false, false, true)
	assert.Error(t, err, "Expected error, as component is not involved and inverseMode set to true!")
	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, true, true, true)
	assert.NoError(t, err, "Expected no error, as component is involved and inverseMode set to true!")
	err = runWrapper(&c, cmdLine, ModeCurrentChanged, false, false, false, false)
	assert.NoError(t, err, "Expected no error, as component is not involved and inverseMode set to false!")
}
