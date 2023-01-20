// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience

import (
	"github.com/pkg/errors"

	"github.com/ory/cli/cmd/cloudx/client"
)

var defaultProjectNotSetError = errors.New("no project was specified")

func getSelectedProjectId(h *client.CommandHelper, args []string) (string, error) {
	if len(args) == 0 {
		if id := h.GetDefaultProjectID(); id == "" {
			return "", defaultProjectNotSetError
		} else {
			return id, nil
		}
	} else {
		return args[0], nil
	}
}
