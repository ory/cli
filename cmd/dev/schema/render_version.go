package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/pkg"
	"github.com/ory/cli/spec"
	"github.com/ory/jsonschema/v3"
	"github.com/ory/jsonschema/v3/httploader"
	"github.com/ory/x/httpx"
	"github.com/ory/x/jsonschemax"
)

var RenderVersion = &cobra.Command{
	Use:   "render-version <project-name> <new-version> <schema-path>",
	Args:  cobra.ExactArgs(3),
	Short: "Renders the version schema for <project-name> and <new-version> in the current directory. The `$ref` is pointing to the <schema-path> relative to the repo root.",
	Run:   addVersionToSchema,
}

var preReleaseVersion = regexp.MustCompile(`.*[-.]pre\.`)

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
			},
			"required": [
				"version"
			]
		},
		{
			"$ref": "%s"
		}
	]
}`, newVersion, ref)

	f, err := os.Open(destFile)
	pkg.Check(err)
	defer f.Close()

	var versionSchema map[string]json.RawMessage
	pkg.Check(json.NewDecoder(f).Decode(&versionSchema))

	var oneOf []json.RawMessage
	if err := json.Unmarshal(versionSchema["oneOf"], &oneOf); err != nil {
		pkg.Check(err)
	}

	// prepend the newest entry
	oneOf = append([]json.RawMessage{json.RawMessage(newVersionEntry)}, oneOf...)

	versionSchema["oneOf"], err = json.Marshal(oneOf)
	pkg.Check(err)

	prettyVersionSchema, err := json.MarshalIndent(versionSchema, "", strings.Repeat(" ", 4))
	pkg.Check(err)

	ctx := context.WithValue(cmd.Context(), httploader.ContextKey, httpx.NewResilientClient())
	schema, err := jsonschema.CompileString(ctx, "version_meta.schema.json", string(spec.VersionSchema))
	pkg.Check(err)

	err = schema.Validate(bytes.NewReader(prettyVersionSchema))
	if err != nil {
		jsonschemax.FormatValidationErrorForCLI(os.Stderr, prettyVersionSchema, err)
		os.Exit(1)
	}

	pkg.Check(ioutil.WriteFile(destFile, prettyVersionSchema, 0600))
}
