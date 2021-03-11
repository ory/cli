package deps

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/ory/x/flagx"
)

const ExampleConfigFile = `version: v1.20.2
url: https://storage.googleapis.com/kubernetes-release/release/{{.Version}}/bin/{{.Os}}/{{.Architecture}}/kubectl
mappings:
  architecture:
    amd64: x64
  os:
    darwin: mac
    linux: unix
`

type FileNotFoundError struct {
	Path string
	Err  error
}

func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("Config file not found: '%s'\n", e.Path)
}
func (e FileNotFoundError) Unwrap() error { return e.Err }

type InvalidFileError struct {
	Path string
	Err  error
}

func (e InvalidFileError) Error() string {
	return fmt.Sprintf("Config file '%s' is invalid:\n%s\n\nExample of a valid config file:\n%s\n", e.Path, e.Err.Error(), ExampleConfigFile)
}
func (e InvalidFileError) Unwrap() error { return e.Err }

type Component struct {
	Version      string   `yaml:"version"`
	Url          string   `yaml:"url"`
	Mappings     Mappings `yaml:"mappings"`
	Os           string   `yaml:"os,omitempty"`
	Architecture string   `yaml:"architecture,omitempty"`
}

type Mappings struct {
	ArchitectureMapping ArchitectureMapping `yaml:"architecture"`
	OsMapping           OsMapping           `yaml:"os"`
}

type ArchitectureMapping struct {
	AMD64 string `yaml:"amd64"`
}

type OsMapping struct {
	Darwin string `yaml:"darwin"`
	Linux  string `yaml:"linux"`
	ready  bool
}

func (c *Component) String() string {
	d, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return fmt.Sprintf("%s", string(d))
}

func (c *Component) getComponentFromConfig(configFilePath string) error {
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return FileNotFoundError{configFilePath, err}
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return InvalidFileError{
			Path: configFilePath,
			Err:  err,
		}
	}
	return nil
}

func (c *Component) getRenderedURL(osString string, archString string) (string, error) {
	c.Os = osString
	if osString == "darwin" && c.Mappings.OsMapping.Darwin != "" && osString != c.Mappings.OsMapping.Darwin {
		c.Os = c.Mappings.OsMapping.Darwin
	}
	if osString == "linux" && c.Mappings.OsMapping.Linux != "" && osString != c.Mappings.OsMapping.Linux {
		c.Os = c.Mappings.OsMapping.Linux
	}
	c.Architecture = archString
	if archString == "amd64" && c.Mappings.ArchitectureMapping.AMD64 != "" && archString != c.Mappings.ArchitectureMapping.AMD64 {
		c.Architecture = c.Mappings.ArchitectureMapping.AMD64
	}
	t := template.Must(template.New("url").Parse(c.Url))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var url = &cobra.Command{
	Use:   "url",
	Short: "Returns the download url based on the provided config file.",
	Long:  `This cmd will help to simplify our Makefile logic to download binary dependencies. As the values used for os and arch as well as the structure of the download url for different binary tools are not standardized it makes it quite cumbersome to handle this efficiently in Makefiles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		component := Component{}
		var pConfig = flagx.MustGetString(cmd, "config")
		err := component.getComponentFromConfig(pConfig)
		if err != nil {
			return err
		}
		var pOS = flagx.MustGetString(cmd, "os")
		var pArch = flagx.MustGetString(cmd, "architecture")
		url, err := component.getRenderedURL(pOS, pArch)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, url)
		return nil
	},
}

func init() {
	Main.AddCommand(url)
}
