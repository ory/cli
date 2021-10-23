package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/ory/cli/spec"

	"github.com/ory/x/jsonschemax"

	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"

	"github.com/ory/cli/cmd/pkg"
	"github.com/ory/jsonschema/v3"
	_ "github.com/ory/jsonschema/v3/httploader"
)

var RenderVersion = &cobra.Command{
	Use:   "render-version <project-name> <new-version> <schema-path>",
	Args:  cobra.ExactArgs(3),
	Short: "Renders the version schema for <project-name> and <new-version> in the current directory. The `$ref` is pointing to the <schema-path> relative to the repo root.",
	Run:   addVersionToSchema,
}

var preReleaseVersion = regexp.MustCompile(".*[-.]pre\\.")

func addVersionToSchema(cmd *cobra.Command, args []string) {
	const destFile = ".schema/version.schema.json"
	project := args[0]
	newVersion := args[1]

	if preReleaseVersion.MatchString(newVersion) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Going to silently skip version schema rendering because '%s' is a pre release", newVersion)
		return
	}

	ref := "https://raw.githubusercontent.com/ory/" + path.Join(project, newVersion, args[2])
	newVersionEntry := fmt.Sprintf(`
{
	"allOf": [
		{
			"properties": {
				"version": {
					"const": "%s"
				}
			}
		},
		{
			"$ref": "%s"
		}
	]
}`, newVersion, ref)

	versionSchema, err := ioutil.ReadFile(destFile)
	pkg.Check(err)

	renderedVersionSchema, err := sjson.SetBytes(versionSchema, "oneOf.-1", json.RawMessage(newVersionEntry))
	pkg.Check(err)

	var prettyVersionSchema bytes.Buffer
	pkg.Check(json.Indent(&prettyVersionSchema, renderedVersionSchema, "", strings.Repeat(" ", 4)))

	schema, err := jsonschema.CompileString("version_meta.schema.json", string(spec.VersionSchema))
	pkg.Check(err)

	err = schema.Validate(bytes.NewBuffer(prettyVersionSchema.Bytes()))
	if err != nil {
		jsonschemax.FormatValidationErrorForCLI(os.Stderr, prettyVersionSchema.Bytes(), err)
		os.Exit(1)
	}

	pkg.Check(ioutil.WriteFile(destFile, prettyVersionSchema.Bytes(), 0600))
}
