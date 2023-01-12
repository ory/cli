// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience

import (
	br "github.com/pkg/browser"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ory/cli/cmd/cloudx/client"
)

const project = "project"

func RegisterProjectFlag(f *flag.FlagSet) {
	f.String(project, "", "The project to use")
}

var axmap = map[string]string{
	"signin":       "login",
	"signup":       "registration",
	"recovery":     "recovery",
	"verification": "verification",
	"settings":     "settings",
}

func NewAccountExperienceOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-experience",
		Short: "Open Ory Account Experience Pages",
	}

	for _, v := range axmap {
		cmd.AddCommand(NewAxCmd(v))
	}

	return cmd
}

func NewAxCmd(subcmd string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   subcmd,
		Short: "Open " + subcmd + " page",
		RunE: func(cmd *cobra.Command, args []string) error {
			return AXWrapper(cmd, args)
		},
	}
	return cmd
}

func AXWrapper(cmd *cobra.Command, args []string) error {
	_, _, p, err := client.Client(cmd)
	if err != nil {
		return err
	}

	url := "https://" + p.GetSlug() + ".projects.oryapis.com/ui/" + axmap[cmd.CalledAs()]

	err = br.OpenURL(url)
	if err != nil {
		return err
	}

	return nil
}
