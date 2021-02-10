package monorepo

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/cobra"
)

var componentMode string

var components = &cobra.Command{
	Use:   "components",
	Short: "List components based on mode.",
	Long:  `Read dependency configs and displays dependency graph.`,
	Run: func(cmd *cobra.Command, args []string) {

		var graph ComponentGraph
		graph.getComponentGraph(rootDirectory)

		switch componentMode {
		case "affected":
			affectedComponents := getAffectedComponents(&graph)
			displayComponents(affectedComponents)
		case "all":
			allComponents := graph.components
			displayComponents(allComponents)
		case "changed":
			changedComponents := getChangedComponents(&graph)
			displayComponents(changedComponents)
		case "involved":
			involvedComponents := getInvolvedComponents(&graph)
			displayComponents(involvedComponents)
		default:
			log.Fatalf("Unknown ListMode '%s'", componentMode)
		}
	},
}

func displayComponents(components []*Component) {
	for _, component := range components {
		if verbose {
			fmt.Println(component.String())
		} else {
			fmt.Println(component.ID)
		}
	}
}

func getInvolvedComponents(graph *ComponentGraph) []*Component {
	components := append(getChangedComponents(graph), getAffectedComponents(graph)...)
	componentSet := mapset.NewSet()
	for _, comp := range components {
		componentSet.Add(comp.ID)
	}

	var involvedComponents []*Component
	for componentID := range componentSet.Iter() {
		involvedComponents = append(involvedComponents, graph.componentIDs[componentID.(string)])
	}
	return involvedComponents
}

func getChangedComponents(graph *ComponentGraph) []*Component {
	changedDirectories, _ := getChangedDirectories(rootDirectory, revisionRange, gitOpts)
	changeDirectoryArray := strings.Split(changedDirectories, "\n")
	var changedComponents []*Component
	componentPaths := graph.componentPaths

	changedComponentIds := mapset.NewSet()

	for _, changedPath := range changeDirectoryArray {
		for path, component := range componentPaths {
			if !changedComponentIds.Contains(component.ID) && strings.HasPrefix(changedPath, path) {
				changedComponents = append(changedComponents, component)
				changedComponentIds.Add(component.ID)
				if debug {
					fmt.Printf("Adding changed component: %s\n", component.String())
				}
			}
		}
	}
	return changedComponents
}

func getChangedComponentIDs(graph *ComponentGraph) []string {
	var changedComponentIds []string
	for _, changedComponent := range getChangedComponents(graph) {
		changedComponentIds = append(changedComponentIds, changedComponent.ID)
	}
	return changedComponentIds
}

func getAffectedComponentIDs(graph *ComponentGraph) []string {
	var affectedComponentIds []string
	for _, affectedComponent := range getAffectedComponents(graph) {
		affectedComponentIds = append(affectedComponentIds, affectedComponent.ID)
	}
	return affectedComponentIds
}

func getAffectedComponents(graph *ComponentGraph) []*Component {
	changedComponents := getChangedComponents(graph)

	affectedComponentIds := mapset.NewSet()
	for _, changedComponent := range changedComponents {
		dependentComponents := changedComponent.getDependentComponents(graph)
		for _, dependentComponent := range dependentComponents {
			affectedComponentIds.Add(dependentComponent.ID)
		}
	}
	acs := 0
	i := 0
	for affectedComponentIds.Cardinality() > acs {
		if verbose {
			fmt.Printf("1) getAffectedComponents: i=%d, acs=%d, new acs=%d, comps: %s\n", i, acs, affectedComponentIds.Cardinality(), affectedComponentIds)
		}
		for id := range affectedComponentIds.Iter() {
			affectedComponent := graph.componentIDs[id.(string)]
			dependentComponents := affectedComponent.getDependentComponents(graph)
			if verbose {
				fmt.Printf("2) getAffectedComponents: affectedComponent id=%s, comps=%v\n", affectedComponent.ID, dependentComponents)
			}

			for _, dependentComponent := range dependentComponents {
				if verbose {
					fmt.Printf("3) getAffectedComponents: adding dependent Component id=%s ....", dependentComponent.ID)
				}
				if !affectedComponentIds.Contains(dependentComponent.ID) {
					affectedComponentIds.Add(dependentComponent.ID)
				}
				if verbose {
					fmt.Println("done!")
				}
			}
		}
		acs = affectedComponentIds.Cardinality()
	}
	var affectedComponents []*Component
	for componentID := range affectedComponentIds.Iter() {
		affectedComponents = append(affectedComponents, graph.componentIDs[componentID.(string)])
	}
	return affectedComponents
}

func (component *Component) isChanged(graph *ComponentGraph) bool {
	changedComponents := getChangedComponents(graph)
	for _, changedComponent := range changedComponents {
		if component.ID == changedComponent.ID {
			return true
		}
	}
	return false
}

func (component *Component) isAffected(graph *ComponentGraph) bool {
	affectedComponents := getAffectedComponents(graph)
	for _, affectedComponent := range affectedComponents {
		if component.ID == affectedComponent.ID {
			return true
		}
	}
	return false
}

func (component *Component) isInvolved(graph *ComponentGraph) bool {
	return component.isChanged(graph) || component.isAffected(graph)
}

func getCurrentComponent() (*Component, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return getComponent(pwd)
}

func getComponent(wd string) (*Component, error) {
	var c Component
	_, err := c.getComponentFromConfig(path.Join(wd, configFile), rootDirectory)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func init() {
	Main.AddCommand(components)
	components.Flags().StringVarP(&componentMode, "mode", "m", "involved", "Define which components you want to get listed (affected, all, changed, involved). Default is 'involved'.")
}
