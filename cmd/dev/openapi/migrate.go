package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/fetcher"
	"github.com/ory/x/flagx"
	"github.com/ory/x/stringsx"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [path/to/swagger2.json] [path/to/output.json]",
	Short: "Migrates Swagger 2.0 to OpenAPI 3.0",
	Long: `Prints the OpenAPI 3.0 spec to std out.

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

		in, err := ioutil.ReadFile(args[0])
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
			content, err := renderPatch(cmd, path)
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

	return errors.WithStack(ioutil.WriteFile(path, indented, 0644))
}

func renderPatch(cmd *cobra.Command, uri string) ([]byte, error) {
	f := fetcher.NewFetcher()
	buf, err := f.Fetch(uri)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	content := buf.Bytes()
	switch ext := filepath.Ext(uri); ext {
	case ".yaml", "yml":
		// Do nothing
	case ".json":
		// We convert this to JSON in order to be able to run text/template on it as text/tempalte
		// gives parse errors for JSON.
		content, err = yaml.JSONToYAML(content)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	default:
		return nil, fmt.Errorf("unexpected patch extension \"%s\", only \"json\", \"yaml\", \"yml\" are supported", ext)
	}

	t, err := template.New(uri).Funcs(template.FuncMap{
		"getenv": func(key string) string {
			return os.Getenv(key)
		},
		"toJson": func(raw interface{}) (string, error) {
			out, err := json.Marshal(raw)
			return string(out), err
		},
	}).Parse(string(content))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var data struct {
		Version          string
		ProjectHumanName string
		HealthPathTags   []string
	}

	data.Version = stringsx.Coalesce(os.Getenv("CIRCLE_TAG"), os.Getenv("CIRCLE_HASH"))
	data.ProjectHumanName = fmt.Sprintf("%s %s",
		stringsx.ToUpperInitial(strings.ToLower(os.Getenv("CIRCLE_PROJECT_USERNAME"))),
		stringsx.ToUpperInitial(strings.ToLower(os.Getenv("CIRCLE_PROJECT_REPONAME"))),
	)
	data.HealthPathTags = flagx.MustGetStringSlice(cmd, "health-path-tags")

	var rendered bytes.Buffer
	if err := t.Execute(&rendered, data); err != nil {
		return nil, errors.WithStack(err)
	}

	content, err = yaml.YAMLToJSON(rendered.Bytes())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `Something went wrong when rendering the template:

%s`, rendered.Bytes())
		return nil, errors.WithStack(err)
	}

	return content, nil
}

func init() {
	migrateCmd.Flags().StringSliceP("patches", "p", []string{}, "JSON Patch file(s) to apply to the final OpenAPI v3.0 spec.")
	migrateCmd.Flags().StringSlice("health-path-tags", []string{"admin"}, "Which tags to set for the /health/* and /version endpoints.")
}
