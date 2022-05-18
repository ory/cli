package project_test

import (
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
	"testing"
)

var (
	defaultProject, defaultConfig, defaultEmail, defaultPassword string
	defaultCmd                                                   *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	defaultConfig, defaultEmail, defaultPassword, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
}
