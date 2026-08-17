// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package relationtuples

import (
	"github.com/spf13/cobra"

	"github.com/ory/keto/cmd/check"
)

func NewAllowedCmd() *cobra.Command {
	cmd := check.NewCheckCmd()
	wrapForOryCLI(cmd)

	// The previous usage line, `allowed <subject> <relation> <namespace> <object>`,
	// documented only the deprecated four-argument form and gave no hint that
	// the subject may be a subject set — so it steered callers into passing a
	// bare subject ID, which Ory Network rejects with "subject_id is not
	// supported; please migrate to subject sets".
	cmd.Use = "allowed <subject> <relation> <object_namespace>:<object_id>"
	// wrapForOryCLI names the command it is normally applied to, so without this
	// `ory is allowed` also answers to `relationships`, `relation-tuples` and
	// friends — the aliases of a different command entirely.
	cmd.Aliases = nil
	cmd.Long = `Check whether a subject has a relation on an object.

The subject is either a subject set ` + "`<namespace>:<object>#<relation>`" + ` or a
plain subject ID. Ory Network no longer accepts relationships written with a
plain subject ID, so a subject set is what a check against it will match.

Passing the object as two separate arguments still works but is deprecated;
use ` + "`<object_namespace>:<object_id>`" + ` instead.`
	cmd.Example = `$ {{ .CommandPath }} 'groups:engineering#member' view documents:readme

{
  "allowed": true
}`

	return cmd
}
