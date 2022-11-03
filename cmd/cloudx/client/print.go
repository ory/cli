// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"

	"github.com/ory/client-go"
)

func (h *CommandHelper) PrintUpdateProjectWarnings(p *client.SuccessfulProjectUpdate) error {
	if len(p.Warnings) > 0 {
		_, _ = fmt.Fprintln(h.VerboseErrWriter)
		_, _ = fmt.Fprintln(h.VerboseErrWriter, "Warnings were found.")
		for _, warning := range p.Warnings {
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "- %s\n", *warning.Message)
		}
		_, _ = fmt.Fprintln(h.VerboseErrWriter, "It is save to ignore these warnings unless your intention was to set these keys.")
	}

	_, _ = fmt.Fprintf(h.VerboseErrWriter, "\nProject updated successfully!\n")
	return nil
}
