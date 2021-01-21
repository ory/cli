package deps

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

const ExampleConfigFile = `version: v1.20.2
url: https://storage.googleapis.com/kubernetes-release/release/{{.Version}}/bin/{{.Os}}/{{.Architecture}}/kubectl
architecture-mapping:
  amd64: x64
os-mapping:
  darwin: mac
  linux: unix
`

type FileNotFoundError struct {
	Path string
	Err  error
}

func (e FileNotFoundError) Error() string { return fmt.Sprintf("Error! Config file not found: '%s'", e.Path) }
func (e FileNotFoundError) Unwrap() error { return e.Err }

type InvalidFileError struct {
	Path string
	Err  error
}

func (e InvalidFileError) Error() string { return fmt.Sprintf("Error! Config file '%s' is invalid:\n%s\n\nExample of a valid config file:\n%s", e.Path, e.Err.Error(), ExampleConfigFile) }
func (e InvalidFileError) Unwrap() error { return e.Err }

type Component struct {
	Version string `yaml:"version"`
	Url string `yaml:"url"`
	ArchitectureMapping ArchitectureMapping `yaml:"architecture-mapping"`
	OsMapping OsMapping `yaml:"os-mapping"`
	Os string `yaml:"os,omitempty"`
	Architecture string `yaml:"architecture,omitempty"`
}
type ArchitectureMapping struct {
	AMD64 string `yaml:"amd64"`
	ready bool
}

type OsMapping struct {
	Darwin string `yaml:"darwin"`
	Linux string `yaml:"linux"`
	ready bool
}

func (component Component) String() string {
	d, err := yaml.Marshal(component)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return fmt.Sprintf("%s", string(d))
}

func (component *Component) getComponent(data []byte) *Component{
	component1 := Component{}
	err := yaml.Unmarshal([]byte(data), &component1)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return &component1
}

func (component *Component) getComponentFromConfig(configFilePath string) (error) {
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return FileNotFoundError{configFilePath, err}
	}
	err = yaml.Unmarshal(yamlFile, component)
	if err != nil {
		return InvalidFileError{
			Path: configFilePath,
			Err:  err,
		}
	}
	return nil
}

func (component *Component) getRenderedUrl(osString string, archString string) (string, error){
	component.Os = 	osString
	if osString == "darwin" && component.OsMapping.Darwin != "" && osString != component.OsMapping.Darwin {
		component.Os = 	component.OsMapping.Darwin
	}
	if osString == "linux" && component.OsMapping.Linux != "" && osString != component.OsMapping.Linux {
		component.Os = 	component.OsMapping.Linux
	}
	component.Architecture = archString
	if archString == "amd64" && component.ArchitectureMapping.AMD64 != "" && archString != component.ArchitectureMapping.AMD64 {
		component.Architecture = component.ArchitectureMapping.AMD64
	}
	t := template.Must(template.New("url").Parse(component.Url))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, component)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var url = &cobra.Command{
	Use:   "url",
	Short: "Returns the download url based on the provided config file.",
	Long: `This cmd will help to simplify our Makefile logic to download binary dependencies. As the values used for os and arch as well as the structure of the download url for different binary tools are not standardized it makes it quite cumbersome to handle this efficiently in Makefiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		component := Component{}
		err:=component.getComponentFromConfig(pConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		url, err2 := component.getRenderedUrl(pOs,pArch)
		if err2 != nil {
			fmt.Fprintln(os.Stderr, err2.Error())
		}
		fmt.Fprintln(os.Stdout, url)
	},
}

func init() {
	Main.AddCommand(url)
}
