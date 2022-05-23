package testhelpers

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"github.com/ory/x/cmdx"
)

func UseStaging() {
	if _, isSet := os.LookupEnv("ORY_CLOUD_CONSOLE_URL"); isSet {
		return
	}
	if err := os.Setenv("ORY_CLOUD_CONSOLE_URL", "https://console.staging.ory.dev"); err != nil {
		panic(err)
	}
}

func CreateDefaultAssets() (defaultConfig, defaultEmail, defaultPassword, defaultProject string, defaultCmd *cmdx.CommandExecuter) {
	UseStaging()

	t := testingT{}

	defaultConfig = NewConfigDir(t)

	defaultEmail, defaultPassword = RegisterAccount(t, defaultConfig)
	defaultProject = CreateProject(t, defaultConfig)
	defaultCmd = ConfigAwareCmd(defaultConfig)
	return
}

func RunAgainstStaging(m *testing.M) {
	UseStaging()
	os.Exit(m.Run())
}

type testingT struct{}

func (testingT) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
	debug.PrintStack()
}

func (testingT) FailNow() {
	os.Exit(1)
}
