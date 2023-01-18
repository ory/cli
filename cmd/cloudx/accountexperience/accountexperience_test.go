// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

var _, _, _, _, defaultProject, defaultCmd = testhelpers.CreateDefaultAssets()

func TestOpenAXPages(t *testing.T) {

	t.Run("is able to open login page", func(t *testing.T) {
		var pages = [5]string{"login", "registration", "recovery", "verification", "settings"}

		for _, p := range pages {
			e := errors.New(fmt.Sprintf("xdg-open: no method available for opening 'https://cool-shamir-px8pubwbfq.projects.oryapis.com/ui/%s'", p))

			_, _, err := defaultCmd.Exec(nil, "open", "account-experience", p, "--project", defaultProject)
			if err != nil || err != e {
				t.Fail()
			}
		}
	})

}
