package project_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
	"github.com/ory/x/cmdx"
)

var (
	defaultProject, defaultConfig, defaultEmail, defaultPassword string
	defaultCmd                                                   *cmdx.CommandExecuter
)

func TestMain(m *testing.M) {
	defaultConfig, defaultEmail, defaultPassword, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()
	testhelpers.RunAgainstStaging(m)
}
