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
	"gopkg.in/yaml.v2"
)

const configFile = "monorepo.yml"

//Component struct represent the configuration stored in yaml enriched by the relative path of the component in relation to the rootDirectory.
type Component struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"deps"`
	Path         string   `yaml:"path"`
	//Graph        *ComponentGraph
}

func (component Component) String() string {
	yamlOutput, err := yaml.Marshal(component)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return string(yamlOutput[:])
}

//ComponentGraph struct represent the graph of all components found in the specified rootDirectory and its subdirectories.
type ComponentGraph struct {
	components            []*Component
	componentIDs          map[string]*Component
	componentDependencies map[string]mapset.Set
	componentPaths        map[string]*Component
}

func (graph *ComponentGraph) getComponentGraph(rootDirectory string) (*ComponentGraph, error) {
	//fmt.Println("Scanning Directory Tree: ", rootDirectory)
	isDirectory, err := isDirectory(rootDirectory)
	if err != nil {
		//log.Fatalf("Error accessing '%s': %s\n", rootDirectory, err)
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
			c.getComponentFromConfig(path, rootDirectory)
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

func (component *Component) getComponentFromConfig(configFilePath, rootDir string) (*Component, error) {
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("Config file not found: '%s'", configFilePath)
	}
	err = yaml.Unmarshal(yamlFile, component)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file '%s', invalid format: %v", configFilePath, err)
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

	//validating if all declared dependencies are part of the graph
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
				//add components with no depedencies directly
				readySet.Add(id)
			}
		}

		// if no components without dependencies were added, but there
		// are still components in the componentDependencies map, there must be a circular
		// dependency
		if readySet.Cardinality() == 0 {
			var g ComponentGraph
			for id := range componentDependencies {
				g.addComponent(graph.componentIDs[id])
			}

			return g, errors.New("Circular dependency found")
		}

		// remove ready components (without dependencies) from the
		// componentDependencies map and add them to resolved graph
		for id := range readySet.Iter() {
			delete(componentDependencies, id.(string))
			resolved.addComponent(graph.componentIDs[id.(string)])
		}

		// cycle through the original componentDependencies and remove
		// components which have been removed in this run from the dependencies
		// of other components
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

// displayGraph displays the ComponentGraph. TODO: not just listed the components in the graph, also add dependency
// information (not sure if this is easily doable in shell)
func (graph *ComponentGraph) displayGraph() {
	graph.listComponents()
}

func (graph ComponentGraph) len() int {
	return len(graph.components)
}

func (component *Component) isDependent(cid string) bool {
	for _, dependencyID := range component.Dependencies {
		if cid == dependencyID {
			return true
		}
	}
	return false
}

func (component *Component) getDependentComponents(graph *ComponentGraph) []*Component {
	var dependentComponents = mapset.NewSet()
	for _, comp := range graph.components {
		if comp.isDependent(component.ID) {
			dependentComponents.Add(comp)
		}
	}
	var dependentComponentsArray []*Component
	for comp := range dependentComponents.Iter() {
		dependentComponentsArray = append(dependentComponentsArray, comp.(*Component))
	}
	return dependentComponentsArray
}

// Adding component to the dependency graph
func (graph *ComponentGraph) addComponent(component *Component) {
	//component.Graph = graph
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
