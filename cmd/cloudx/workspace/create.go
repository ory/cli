// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const nameFlag = "name"

func NewCreateCmd() *cobra.Command {
	name := ""

	cmd := &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"workspaces", "ws"},
		Short:   "Create a new Ory Network workspace",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			if name == "" && flagx.MustGetBool(cmd, cmdx.FlagQuiet) {
				return errors.New("you must specify the --name flag when using --quiet")
			}

			for name == "" {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Enter a name for your workspace: ")
				name, err = h.Stdin.ReadString('\n')
				if err != nil {
					return errors.Wrap(err, "failed to read from stdin")
				}
			}

			ws, err := h.CreateWorkspace(cmd.Context(), name)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintf(h.VerboseErrWriter, "Workspace created successfully!")
			cmdx.PrintRow(cmd, (*outputWorkspace)(ws))
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, nameFlag, "n", "", "The name of the workspace, required when quiet mode is used")

	return cmd
}
