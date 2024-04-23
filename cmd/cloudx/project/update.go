// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

func NewProjectsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Update Ory Network project service configuration",
		Example: `$ ory update project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--name \"my updated name\" \
	--file /path/to/config.json \
	--file /path/to/config.yml \
	--file https://example.org/config.yaml \
	--file base64://<json>

ID		ecaaa3cb-0730-4ee8-a6df-9553cdfeef89
SLUG	good-wright-t7kzy3vugf
STATE	running
NAME	Example Project

$ ory update project ecaaa3cb-0730-4ee8-a6df-9553cdfeef89 \
	--name \"my updated name\" \
	--file /path/to/config.json \
	--format json-pretty

{
  "name": "my updated name",
  "identity": {
	"services": {
	  "config": {
		"courier": {
		  "smtp": {
			"from_name": "..."
		  }
		  // ...
		}
	  }
	}
  }
}`,
		Long: `Updates your Ory Network project's service configuration. All values
will be overwritten. To update individual settings use the ` + "`patch`" + ` command instead.

If the ` + "`--name`" + ` flag is not set, the project's name will not be changed.

The full configuration payload can be found at

	https://www.ory.sh/docs/reference/api#operation/updateProject

As an example an input could look like:

	{
      "name": "my updated name",
	  "identity": {
		"services": {
		  "config": {
			"courier": {
			  "smtp": {
				"from_name": "..."
			  }
			  // ...
			}
		  }
		}
	  }
	}
`,
		RunE: runUpdate(prefixFileNop, outputFullProject),
	}

	cmd.Flags().StringP("name", "n", "", "The new name of the project.")
	cmd.Flags().StringSliceP("file", "f", nil, "Configuration file(s) (file://config.json, https://example.org/config.yaml, ...) to update the project")
	client.RegisterYesFlag(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}

func runUpdate(filePrefixer func([]json.RawMessage) ([]json.RawMessage, error), outputter func(*cobra.Command, *cloud.SuccessfulProjectUpdate)) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		h, err := client.NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		var configs []json.RawMessage
		files := flagx.MustGetStringSlice(cmd, "file")
		if len(files) == 0 {
			content, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return errors.New("error reading from STDIN: use --file flag to read from a file instead: " + err.Error())
			}
			configs = []json.RawMessage{content}
		} else {
			configs, err = client.ReadConfigFiles(files)
			if err != nil {
				return err
			}
		}

		configs, err = filePrefixer(configs)
		if err != nil {
			return err
		}

		name := ""
		if n := cmd.Flags().Lookup("name"); n != nil {
			name = n.Value.String()
		}
		id, err := selectedProjectID(h, args)
		if err != nil {
			return cmdx.PrintOpenAPIError(cmd, err)
		}
		p, err := h.UpdateProject(id, name, configs)
		if err != nil {
			return cmdx.PrintOpenAPIError(cmd, err)
		}

		outputter(cmd, p)
		return h.PrintUpdateProjectWarnings(p)
	}
}
