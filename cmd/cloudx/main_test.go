// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cloudx_test

import (
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

func TestMain(m *testing.M) {
	testhelpers.UseStaging()
	m.Run()
}
