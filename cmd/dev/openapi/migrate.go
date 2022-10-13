// Copyright © 2022 Ory Corp

package openapi

import (
	"encoding/json"
	"os"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/pkg"
	"github.com/ory/x/flagx"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [path/to/swagger2.json] [path/to/output.json]",
	Short: "Migrates Swagger 2.0 to OpenAPI 3.0",
	Long: `Migrates Swagger 2.0 to OpenAPI 3.0. Prints the OpenAPI 3.0 spec to std out.

This command can also apply a JSON Patch (https://tools.ietf.org/html/rfc7396) to the OpenAPI 3.0 output using
the --patch flag. The path can be a file:// or https:// path.

$ ory dev openapi convert -p file://path/to/patch1.(yaml|json) -p https://foo.bar/path/to/patch2.(yaml|json)

YAML-encoded patches support Go text templates with the following functions

- {{ getenv "ENV_VAR_NAME" }}: returns the environment variable "ENV_VAR_NAME"
- {{ toJson .Field }}: returns the jsonified version of the first argument.

and these fields:

- {{ .Version }}: the software version extracted from CIRCLE_TAG and CIRCLE_HASH environment variables.
- {{ .ProjectHumanName }}: the software's human readable name (e.g. Ory Kratos, Ory Hydra)
- {{ .HealthPathTags }}: the tags for health and version APIs.

Example:

	# some/patch.yaml
	- op: replace
	  path: /info
	  value:
		version: >-
		  {{ getenv "CIRCLE_TAG" }}

	- op: replace
	  path: /info
	  value:
		version: >-
		  {{ .Version }}
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var oas2 openapi2.Swagger

		in, err := os.ReadFile(args[0])
		if err != nil {
			return errors.WithStack(err)
		}

		if err := json.Unmarshal(in, &oas2); err != nil {
			return errors.WithStack(err)
		}

		oas3, err := openapi2conv.ToV3Swagger(&oas2)
		if err != nil {
			return errors.WithStack(err)
		}

		result, err := json.MarshalIndent(oas3, "", "  ")
		if err != nil {
			return errors.WithStack(err)
		}

		patches := flagx.MustGetStringSlice(cmd, "patches")
		if len(patches) == 0 {
			return renderFile(args[1], result)
		}

		for _, path := range patches {
			content, err := pkg.RenderOASPatch(cmd, path)
			if err != nil {
				return errors.WithStack(err)
			}

			patch, err := jsonpatch.DecodePatch(content)
			if err != nil {
				return errors.WithStack(err)
			}

			result, err = patch.Apply(result)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		return renderFile(args[1], result)
	},
}

func renderFile(path string, content []byte) error {
	indented, err := json.MarshalIndent(json.RawMessage(content), "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(os.WriteFile(path, indented, 0644))
}

func init() {
	migrateCmd.Flags().StringSliceP("patches", "p", []string{}, "JSON Patch file(s) to apply to the final OpenAPI v3.0 spec.")
	migrateCmd.Flags().StringSlice("health-path-tags", []string{"admin"}, "Which tags to set for the /health/* and /version endpoints.")
}
