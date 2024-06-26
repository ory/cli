// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package accountexperience

import (
	"fmt"
	"path"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/flagx"
	"github.com/ory/x/stringsx"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
)

func NewAccountExperienceOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-experience <login|registration|recovery|verification|settings>",
		Aliases: []string{"ax", "ui"},
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			switch f := stringsx.SwitchExact(args[0]); {
			case f.AddCase("login", "registration", "recovery", "verification", "settings"):
				return nil
			default:
				return errors.Wrap(f.ToUnknownCaseErr(), "unknown flow type")
			}
		},
		Short: "Open Ory Account Experience Pages",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			project, err := h.GetSelectedProject(cmd.Context())
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			url := client.CloudAPIsURL(project.Slug)
			url.Path = path.Join(url.Path, "ui", args[0])
			if flagx.MustGetBool(cmd, cmdx.FlagQuiet) {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", url)
				return nil
			}
			if err := browser.OpenURL(url.String()); err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\n\nUnable to automatically open %s in your browser. Please open it manually!\n", err, url)
				return cmdx.FailSilently(cmd)
			}
			return nil
		},
	}
	cmdx.RegisterNoiseFlags(cmd.Flags())
	return cmd
}
