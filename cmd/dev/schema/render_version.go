package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"

	"github.com/ory/jsonschema/v3"
	_ "github.com/ory/jsonschema/v3/httploader"
	"github.com/ory/x/viperx"

	"github.com/ory/cli/cmd/pkg"
)

var _ = pkger.Include("../../../.schema")

var RenderVersion = &cobra.Command{
	Use:   "render-version <project-name> <new-version>",
	Args:  cobra.ExactArgs(2),
	Short: "Renders the version schema for <project-name> and <new-version> in the current directory.",
	Run:   addVersionToSchema,
}

func addVersionToSchema(_ *cobra.Command, args []string) {
	const destFile = ".schema/version.schema.json"
	project := args[0]
	newVersion := args[1]

	if strings.Contains(newVersion, ".pre.") {
		fmt.Printf("Going to silently skip version schema rendering because '%s' contains '.pre.'", newVersion)
		return
	}

	ref := fmt.Sprintf("https://raw.githubusercontent.com/ory/%s/%s/.schema/config.schema.json", project, newVersion)
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

	f, err := pkger.Open("/.schema/version_meta.schema.json")
	pkg.Check(err)
	metaSchema, err := ioutil.ReadAll(f)
	pkg.Check(err)
	pkg.Check(f.Close())
	schema, err := jsonschema.CompileString("version_meta.schema.json", string(metaSchema))
	pkg.Check(err)

	err = schema.Validate(bytes.NewBuffer(prettyVersionSchema.Bytes()))
	if err != nil {
		viperx.PrintHumanReadableValidationErrors(os.Stderr, err)
		os.Exit(1)
	}

	pkg.Check(ioutil.WriteFile(destFile, prettyVersionSchema.Bytes(), 0600))
}
