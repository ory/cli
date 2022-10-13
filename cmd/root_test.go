// Copyright © 2022 Ory Corp

package cmd

import (
	"testing"

	"github.com/ory/x/cmdx"
)

func TestUsageTemplating(t *testing.T) {
	cmdx.AssertUsageTemplates(t, NewRootCmd())
}
