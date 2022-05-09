package project_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestMain(m *testing.M) {
	testhelpers.RunAgainstStaging(m)
}
