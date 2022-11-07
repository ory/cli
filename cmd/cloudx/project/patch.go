// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

func NewProjectsPatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Patch the Ory Network project configuration.",
		Example: `ory patch project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--replace '/name="My new project name"' \
	--add '/services/identity/config/courier/smtp={"from_name":"My new email name"}' \
	--replace '/services/identity/config/selfservice/methods/password/enabled=false' \
	--delete '/services/identity/config/selfservice/methods/totp/enabled'

ory patch project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--replace '/name="My new project name"' \
	--delete '/services/identity/config/selfservice/methods/totp/enabled'
	--format kratos-config > my-config.yaml`,
		Long: `Patch the Ory Network project configuration. Only values
specified in the patch will be overwritten. To replace the config use the ` + "`update`" + ` command instead.

The format of the patch is a JSON-Patch document. For more details please check:

	https://www.ory.sh/docs/reference/api#operation/patchProject
	https://jsonpatch.com`,
		RunE: runPatch(
			func(s []string) []string {
				return s
			},
			prefixFileNop,
			outputFullProject,
		),
	}

	cmd.Flags().StringSliceP("file", "f", nil, "Configuration file(s) (file://config.json, https://example.org/config.yaml, ...) to update the project")
	cmd.Flags().StringArray("replace", nil, "Replace a specific key in the configuration")
	cmd.Flags().StringArray("add", nil, "Add a specific key to the configuration")
	cmd.Flags().StringArray("remove", nil, "Remove a specific key from the configuration")
	client.RegisterYesFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}

func runPatch(patchPrefixer func([]string) []string, filePrefixer func([]json.RawMessage) ([]json.RawMessage, error), outputter func(*cobra.Command, *cloud.SuccessfulProjectUpdate)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		h, err := client.NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		files := patchPrefixer(flagx.MustGetStringSlice(cmd, "file"))
		add := patchPrefixer(flagx.MustGetStringArray(cmd, "add"))
		replace := patchPrefixer(flagx.MustGetStringArray(cmd, "replace"))
		remove := patchPrefixer(flagx.MustGetStringArray(cmd, "remove"))

		if len(files)+len(add)+len(replace)+len(remove) == 0 {
			return errors.New("at least one of --file, --add, --replace, or --remove must be set")
		}

		configs, err := client.ReadConfigFiles(files)
		if err != nil {
			return err
		}

		configs, err = filePrefixer(configs)
		if err != nil {
			return err
		}

		p, err := h.PatchProject(args[0], configs, add, replace, remove)
		if err != nil {
			return cmdx.PrintOpenAPIError(cmd, err)
		}

		outputter(cmd, p)
		return h.PrintUpdateProjectWarnings(p)
	}
}
