package monorepo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Component struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"deps"`
	Path         string   `yaml:"path"`
}

func (component Component) String() string {
	yamlOutput, err := yaml.Marshal(component)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return string(yamlOutput[:])
}

type ComponentGraph struct {
	components            []*Component
	componentIDs          map[string]*Component
	componentDependencies map[string]mapset.Set
	componentPaths        map[string]*Component
}

const configFile = "monorepo.yml"

var dep = &cobra.Command{
	Use:   "dep",
	Short: "Read dependency configs and displays dependency graph.",
	Long:  `Read dependency configs and displays dependency graph.`,
	Run: func(cmd *cobra.Command, args []string) {

		var graph ComponentGraph
		graph.readConfiguration(rootDirectory)
		graph.displayGraph()

		resolved, err := graph.resolveGraph()
		if err != nil {
			fmt.Printf("Failed to resolve dependency graph: %s\n", err)
		} else {
			fmt.Println("The dependency graph resolved successfully")
		}

		if len(changedComponentIds) > 0 {
			ids := strings.Split(changedComponentIds, ",")
			triggers := resolved.triggerComponents(ids)
			for id, component := range triggers {
				fmt.Printf(" %s -> %v \n", id, component)
			}
		}
	},
}

func (graph *ComponentGraph) readConfiguration(rootDirectory string) (*ComponentGraph, error) {
	//fmt.Println("Scanning Directory Tree: ", rootDirectory)
	isDirectory, err := isDirectory(rootDirectory)
	if err != nil {
		//log.Fatalf("Error acceesing '%s': %s\n", rootDirectory, err)
		return nil, err
	}
	if !isDirectory {
		return nil, fmt.Errorf("Provided path '%s' is not a directory", rootDirectory)
	}
	filepath.Walk(rootDirectory, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // not a file.  ignore.
		}
		matched, err := filepath.Match(configFile, fi.Name())
		if err != nil {
			fmt.Println(err) // malformed pattern
			return err       // this is fatal.
		}
		if matched {
			if debug {
				fmt.Printf("Debug: reading config file '%s'\n", path)
			}
			var c Component
			c.getMonoRepoConfigFile(path, rootDirectory)
			graph.addComponent(&c)
		}
		return nil
	})
	return graph, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func (component *Component) getMonoRepoConfigFile(configFilePath, rootDir string) (*Component, error) {
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading monorepo config file '%s': %v", configFilePath, err)
	}
	err = yaml.Unmarshal(yamlFile, component)
	if err != nil {
		return nil, fmt.Errorf("Error reading monorepo config file '%s', invalid format", configFilePath)
	}

	configFilePath, err = filepath.Abs(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("Error determining absolute config file path '%s': %v", configFilePath, err)
	}
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("Error determining absolute root directory path '%s': %v", rootDir, err)
	}
	rootDir = rootDir + "/"
	component.Path = strings.TrimSuffix(strings.TrimPrefix(configFilePath, rootDir), "/"+configFile)
	return component, nil
}

// Resolves the dependency graph
func (graph *ComponentGraph) resolveGraph() (ComponentGraph, error) {
	componentIDs := graph.componentIDs
	componentDependencies := graph.componentDependencies

	for id, deps := range componentDependencies {
		for dep := range deps.Iter() {
			_, found := componentIDs[dep.(string)]
			if !found {
				log.Fatalf("Component '%s': dependency '%s' unknown!\n", id, dep.(string))
			}
		}
	}

	var resolved ComponentGraph
	for len(componentDependencies) != 0 {
		readySet := mapset.NewSet()
		for id, deps := range componentDependencies {
			if deps.Cardinality() == 0 {
				//fmt.Printf(" - adding component without deps: '%s'\n", id)
				readySet.Add(id)
			}
		}

		if readySet.Cardinality() == 0 {
			var g ComponentGraph
			for id := range componentDependencies {
				g.addComponent(graph.componentIDs[id])
			}

			return g, errors.New("Circular dependency found")
		}

		for id := range readySet.Iter() {
			//fmt.Printf(" - removing ready compoents and adding them to resolved graph: '%s'\n", id)
			delete(componentDependencies, id.(string))
			resolved.addComponent(graph.componentIDs[id.(string)])
		}

		for id, deps := range graph.componentDependencies {
			diff := deps.Difference(readySet)
			graph.componentDependencies[id] = diff
		}
	}

	return resolved, nil
}

func (graph *ComponentGraph) listComponents() {
	for _, component := range graph.components {
		if verbose {
			fmt.Printf("%v\n", component)
		} else {
			fmt.Printf("%s\n", component.ID)
		}
	}
}

// Displays the dependency graph
func (graph *ComponentGraph) displayGraph() {
	for _, component := range graph.components {
		if verbose {
			fmt.Printf("%v\n", component)
		} else {
			fmt.Printf("%s\n", component.ID)
		}
	}
}

func (graph ComponentGraph) len() int {
	return len(graph.components)
}

func (graph *ComponentGraph) triggerComponents(paths []string) map[string]*Component {
	triggerMap := make(map[string]*Component)

	for _, path := range paths {
		//1:1 comparisson not correct, we need to check if there is a change for the component directory or below
		changedComponent, found := graph.componentPaths[path]
		if found {
			triggerMap[changedComponent.Path] = changedComponent
		} else {
			fmt.Printf("No component defined for path '%s'", path)
		}
	}

	for _, component := range graph.components {
		for _, dependencyID := range component.Dependencies {
			dependencyComponent := graph.componentIDs[dependencyID]
			_, found := graph.componentPaths[dependencyComponent.Path]
			if found {
				triggerMap[component.Path] = component
				break
			}
		}
	}

	return triggerMap
}

// Adding component to the dependency graph
func (graph *ComponentGraph) addComponent(component *Component) {
	graph.components = append(graph.components, component)
	if graph.componentIDs == nil {
		graph.componentIDs = make(map[string]*Component)
	}
	graph.componentIDs[component.ID] = component

	if graph.componentPaths == nil {
		graph.componentPaths = make(map[string]*Component)
	}
	graph.componentPaths[component.Path] = component

	if graph.componentDependencies == nil {
		graph.componentDependencies = make(map[string]mapset.Set)
	}
	dependencySet := mapset.NewSet()
	for _, dep := range component.Dependencies {
		dependencySet.Add(dep)
	}
	graph.componentDependencies[component.ID] = dependencySet
}

func init() {
	Main.AddCommand(dep)
	//func StringVarP(p *string, name, shorthand string, value string, usage string) {

}
