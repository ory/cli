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

	cmd.Use = "allowed <subject> <relation> <object_namespace>:<object_id>"
	// wrapForOryCLI sets the aliases of the relationships command, which do not
	// belong on a permission check.
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
