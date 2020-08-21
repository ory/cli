package monorepo

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var componentMode string

var components = &cobra.Command{
	Use:   "components",
	Short: "List components based on mode.",
	Long:  `Read dependency configs and displays dependency graph.`,
	Run: func(cmd *cobra.Command, args []string) {

		var graph ComponentGraph
		graph.readConfiguration(rootDirectory)

		switch componentMode {
		case "affected":
			fmt.Println("Not implemented yet!")
		case "all":
			graph.listComponents()
		case "changed":

			//get changes from Git CLI
			changedDirectories, _ := getChangedDirectories(rootDirectory, "origin/master")

			changedComponents := detectChangedComponents(graph, changedDirectories)
			for _, component := range changedComponents {
				fmt.Println(component.String())
			}

		case "involved":
			fmt.Println("Not implemented yet!")
		default:
			log.Fatalf("Unknown ListMode '%s'", componentMode)
		}

		/*
			resolved, err := graph.resolveGraph()
			if err != nil {
				fmt.Printf("Failed to resolve dependency graph: %s\n", err)
			} else {
				fmt.Println("The dependency graph resolved successfully")
			}
		*/
	},
}

func detectChangedComponents(graph ComponentGraph, changedDirectories string) []*Component {
	changeDirectoryArray := strings.Split(changedDirectories, "\n")
	var changedComponents []*Component
	componentPaths := graph.componentPaths
	for _, changedPath := range changeDirectoryArray {
		for path, component := range componentPaths {
			//fmt.Printf("'%s' // '%s'\n", changedPath, path)
			if strings.HasPrefix(changedPath, path) {
				changedComponents = append(changedComponents, component)
				//	fmt.Printf("Adding changed component: %s\n", component.String())
			}
		}
	}
	return changedComponents
}

func init() {
	Main.AddCommand(components)
	components.Flags().StringVarP(&componentMode, "mode", "m", "all", "Define which components you want to get listed (affected, all, changed, involved). Default is all.")
}
