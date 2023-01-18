// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

const useProjectFlag = "use-project"

func NewCreateProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Create a new Ory Network project",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			name := flagx.MustGetString(cmd, "name")
			if len(name) == 0 && flagx.MustGetBool(cmd, cmdx.FlagQuiet) {
				return errors.New("you must specify the --name flag when using --quiet")
			}

			stdin := h.Stdin
			for name == "" {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Enter a name for your project: ")
				name, err = stdin.ReadString('\n')
				if err != nil {
					return errors.Wrap(err, "failed to read from stdin")
				}
			}

			use := flagx.MustGetBool(cmd, useProjectFlag)
			p, err := h.CreateProject(name, use)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Project created successfully!")
			cmdx.PrintRow(cmd, (*outputProject)(p))
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "The name of the project, required when quiet mode is used")
	cmd.Flags().Bool(useProjectFlag, false, "Set the created project as the default.")
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
