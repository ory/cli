package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/fetcher"
	"github.com/ory/x/flagx"
	"github.com/ory/x/stringsx"
)

func RenderOASPatch(cmd *cobra.Command, uri string) ([]byte, error) {
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
