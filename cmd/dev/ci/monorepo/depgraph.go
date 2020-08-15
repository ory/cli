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

type ComponentGraph struct {
	components            []*Component
	componentIDs          map[string]*Component
	componentDependencies map[string]mapset.Set
}

var rootDirectory string
var changedComponentIds string
var graph ComponentGraph

var dep = &cobra.Command{
	Use:   "dep",
	Short: "Ready dependency configs and displays dependency graph.",
	Long:  `Ready dependency configs and displays dependency graph.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root Directory: " + rootDirectory)
		isDirectory, err := isDirectory(rootDirectory)
		if err != nil {
			fmt.Printf("Error acceesing '%s': %s\n", rootDirectory, err)
			return
		}
		if !isDirectory {
			fmt.Printf("Provided path '%s' is not a directory!\n", rootDirectory)
			return
		}
		filepath.Walk(rootDirectory, visitFile)

		graph.displayGraph()
		resolved, err := resolveGraph(&graph)
		if err != nil {
			fmt.Printf("Failed to resolve dependency graph: %s\n", err)
		} else {
			fmt.Println("The dependency graph resolved successfully")
		}

		for _, component := range resolved.components {
			fmt.Println(component.ID)
		}

		if len(changedComponentIds) > 0 {
			ids := strings.Split(changedComponentIds, ",")
			fmt.Println("IDS: %s ", ids)
			triggers := resolved.triggerComponents(ids)
			for id := range triggers.Iter() {
				fmt.Printf(" -> %s \n", id)
			}
		}

	},
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func (c *Component) getMonoRepoConfigFile(path string) *Component {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf("Error reading monorepo config file '%s': %v \n", path, err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func visitFile(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if fi.IsDir() {
		return nil // not a file.  ignore.
	}
	matched, err := filepath.Match("monorepo.yml", fi.Name())
	if err != nil {
		fmt.Println(err) // malformed pattern
		return err       // this is fatal.
	}
	if matched {
		fmt.Println(fp)
		var c Component
		c.getMonoRepoConfigFile(fp)
		graph.addComponent(&c)
		fmt.Printf("Components in Graph: %d\n", graph.len())
	}
	return nil
}

// Resolves the dependency graph
func resolveGraph(graph *ComponentGraph) (ComponentGraph, error) {

	componentIDs := graph.componentIDs
	componentDependencies := graph.componentDependencies

	//validate only defined ids are used in dependencies
	for id, deps := range componentDependencies {
		for dep := range deps.Iter() {
			_, found := componentIDs[dep.(string)]
			if !found {
				log.Fatalf("Component '%s': dependency '%s' unknown!\n", id, dep.(string))
			}
		}
	}

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved ComponentGraph
	for len(componentDependencies) != 0 {
		// Get all nodes from the graph which have no dependencies
		readySet := mapset.NewSet()
		for id, deps := range componentDependencies {
			if deps.Cardinality() == 0 {
				fmt.Printf(" - adding component without deps: '%s'\n", id)
				readySet.Add(id)
			}
		}

		// If there aren't any ready nodes, then we have a cicular dependency
		if readySet.Cardinality() == 0 {
			var g ComponentGraph
			for id := range componentDependencies {
				g.addComponent(g.componentIDs[id])
			}

			return g, errors.New("Circular dependency found")
		}

		// Remove the ready nodes and add them to the resolved graph
		for id := range readySet.Iter() {
			fmt.Printf(" - removing ready compoents and adding them to resolved graph: '%s'\n", id)
			delete(componentDependencies, id.(string))
			resolved.addComponent(graph.componentIDs[id.(string)])
		}

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for id, deps := range graph.componentDependencies {
			diff := deps.Difference(readySet)
			graph.componentDependencies[id] = diff
		}
	}

	return resolved, nil
}

// Displays the dependency graph
func (graph *ComponentGraph) displayGraph() {
	fmt.Println("-Display Graph-----")
	for _, component := range graph.components {
		if len(component.Dependencies) == 0 {
			fmt.Printf("%s -|\n", component.ID)
		} else {
			for _, dep := range component.Dependencies {
				fmt.Printf("%s -> %s\n", component.ID, dep)
			}
		}
	}
	fmt.Println("-End-Display Graph-----")
}

func (graph ComponentGraph) len() int {
	return len(graph.components)
}

func (graph *ComponentGraph) triggerComponents(ids []string) mapset.Set {
	triggerSet := mapset.NewSet()
	for _, item := range ids {
		triggerSet.Add(item)
	}

	fmt.Printf("trigger? %d\n", graph.len())
	for _, component := range graph.components {
		for _, dependencyID := range component.Dependencies {
			if triggerSet.Contains(dependencyID) {
				fmt.Printf(" * Triggering Component: '%s'\n", component.ID)
				triggerSet.Add(component.ID)
				break
			}
		}
	}

	return triggerSet
}

// Adding component to the dependency graph
func (graph *ComponentGraph) addComponent(component *Component) {
	graph.components = append(graph.components, component)
	if graph.componentIDs == nil {
		graph.componentIDs = make(map[string]*Component)
	}
	graph.componentIDs[component.ID] = component

	fmt.Printf(" + Adding Component: '%s'\n", component.ID)
	if graph.componentDependencies == nil {
		graph.componentDependencies = make(map[string]mapset.Set)
	}
	dependencySet := mapset.NewSet()
	for _, dep := range component.Dependencies {
		dependencySet.Add(dep)
		fmt.Printf(" +- dependency: '%s'\n", dep)
	}
	graph.componentDependencies[component.ID] = dependencySet
}

func init() {
	Main.AddCommand(dep)
	//func StringVarP(p *string, name, shorthand string, value string, usage string) {
	dep.Flags().StringVarP(&rootDirectory, "root", "r", ".", "Root directory to be used to traverse and search for dependency configurations.")
	dep.Flags().StringVarP(&changedComponentIds, "changed", "c", "", "Changed Components IDs.")
}
