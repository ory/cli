// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience

import (
	"fmt"

	br "github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
)

const project = "project"

// func RegisterProjectFlag(f *flag.FlagSet) {
// 	f.String(project, "", "The project to use")
// }

var pages = [5]string{"login", "registration", "recovery", "verification", "settings"}

func NewAccountExperienceOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-experience ",
		Short: "Open Ory Account Experience Pages",
	}

	for _, p := range pages {
		cmd.AddCommand(NewAxCmd(p))
	}

	return cmd
}

func NewAxCmd(subcmd string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [project-id]", subcmd),
		Short: "Open " + subcmd + " page",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(subcmd *cobra.Command, args []string) error {
			// id, err := project.getSelectedProjectId(subcmd, args)
			// if err != nil {
			// 	return cmdx.PrintOpenAPIError(subcmd, err)
			// }
			return AXWrapper(subcmd, args)
		},
	}
	return cmd
}

func AXWrapper(cmd *cobra.Command, args []string) error {
	_, _, p, err := client.Client(cmd)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s.projects.oryapis.com/ui/%s", p.GetSlug(), cmd.CalledAs())
	err = br.OpenURL(url)
	if err != nil {
		return err
	}

	return nil
}
