package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/markbates/pkger"
	"github.com/ory/cli/cmd/pkg"
	"github.com/ory/jsonschema/v3"
	_ "github.com/ory/jsonschema/v3/httploader"
	"github.com/ory/x/viperx"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
)

var _ = pkger.Include("../../../.schema")

var RenderVersion = &cobra.Command{
	Use:   "render-version",
	Short: "Release infrastructure for ORY and related components",
	Run:   addVersionToSchema,
}

func addVersionToSchema(_ *cobra.Command, args []string) {
	const destFile = ".schema/version.schema.json"
	newVersion := args[0]
	wd, err := os.Getwd()
	pkg.Check(err)

	project := pkg.ProjectFromDir(wd)
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
