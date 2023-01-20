// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience

import (
	"fmt"
	"os"

	"os/exec"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	client "github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/project"
	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
)

func NewAccountExperienceOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-experience [project-id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Open Ory Account Experience Pages",
	}
	var pages = [5]string{"login", "registration", "recovery", "verification", "settings"}
	for _, p := range pages {
		cmd.AddCommand(NewAxCmd(p))
	}

	return cmd
}

func NewAxCmd(cmd string) *cobra.Command {
	return &cobra.Command{
		Use:   cmd,
		Short: "Open " + cmd + " page",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}
			id, err := project.GetSelectedProjectId(h, args)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			project, err := h.GetProject(id)
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			return AxWrapper(cmd, project)

		}}
}

func AxWrapper(cmd *cobra.Command, p *cloud.Project) error {
	url := fmt.Sprintf("https://%s.projects.oryapis.com/ui/%s", p.GetSlug(), cmd.CalledAs())

	err := browser.OpenURL(url)
	if err != nil {

		// #nosec G204 - this is ok
		if err := exec.Command("open", url); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to automatically open the %s page in your browser. Please open it manually!", cmd.CalledAs())
		}
	}

	return nil
}
