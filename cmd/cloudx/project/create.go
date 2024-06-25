// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	"github.com/ory/cli/cmd/cloudx/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/stringsx"
)

const (
	nameFlag        = "name"
	environmentFlag = "environment"
	useProjectFlag  = "use-project"
)

func NewCreateProjectCmd() *cobra.Command {
	name := ""
	environment := environmentValue("dev")
	useProject := false

	cmd := &cobra.Command{
		Use:   "project",
		Short: "Create a new Ory Network project",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			if (len(name) == 0 || len(environment) == 0) && flagx.MustGetBool(cmd, cmdx.FlagQuiet) {
				return errors.New("you must specify the --name and --environment flags when using --quiet")
			}

			for name == "" {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Enter a name for your project: ")
				name, err = h.Stdin.ReadString('\n')
				if err != nil {
					return errors.Wrap(err, "failed to read from stdin")
				}
			}

			p, err := h.CreateProject(ctx, name, string(environment), h.WorkspaceID(), useProject)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "Project created successfully!")
			cmdx.PrintRow(cmd, (*outputProject)(p))
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, nameFlag, "n", "", "The name of the project, required when quiet mode is used")
	cmd.Flags().VarP(&environment, environmentFlag, "e", "The environment of the project. Valid values are: prod, stage, dev")
	cmd.Flags().BoolVar(&useProject, useProjectFlag, false, "Set the created project as the default.")
	client.RegisterWorkspaceFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}

type environmentValue string

const (
	EnvironmentProduction  environmentValue = "prod"
	EnvironmentStaging     environmentValue = "stage"
	EnvironmentDevelopment environmentValue = "dev"
)

var _ pflag.Value = (*environmentValue)(nil)

func (e *environmentValue) normalize() {
	if e == nil {
		return
	}
	switch *e {
	case "production", "p":
		*e = EnvironmentProduction
	case "staging", "s":
		*e = EnvironmentStaging
	case "development", "d":
		*e = EnvironmentDevelopment
	}
}

func (e *environmentValue) valid() error {
	if e == nil {
		return errors.Errorf("environment value is nil")
	}
	switch c := stringsx.SwitchExact(string(*e)); {
	case c.AddCase(string(EnvironmentProduction)),
		c.AddCase(string(EnvironmentStaging)),
		c.AddCase(string(EnvironmentDevelopment)):
		return nil
	default:
		return c.ToUnknownCaseErr()
	}
}

func (e *environmentValue) String() string {
	if e == nil {
		return ""
	}
	return string(*e)
}

func (e *environmentValue) Set(s string) error {
	se := environmentValue(s)
	se.normalize()
	if err := se.valid(); err != nil {
		return err
	}

	*e = se
	return nil
}

func (e *environmentValue) Type() string {
	return "environment"
}
