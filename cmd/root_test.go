// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/ory/x/cmdx"
)

func TestUsageTemplating(t *testing.T) {
	cmdx.AssertUsageTemplates(t, NewRootCmd())
}
