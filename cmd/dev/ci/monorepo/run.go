package monorepo

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var cmds string
var runMode string
var dryRun bool

var run = &cobra.Command{
	Use:   "run",
	Short: "Runs the specified commands if the current component (which is defined in the current work directory) is affected by any change in the repository.",
	Long:  `Runs the specified commands if the current component (which is defined in the current work directory) is affected by any change in the repository.`,
	Run: func(cmd *cobra.Command, args []string) {

		var graph ComponentGraph
		graph.getComponentGraph(rootDirectory)

		switch runMode {
		case "current_involved":
			c, err := getCurrentComponent()
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("Current Component '%s'!\n", c.String())
			isAffected := c.isAffected(&graph)
			isChanged := c.isChanged(&graph)
			isInvolved := isChanged || isAffected
			fmt.Printf("%s is involved: %t (isChanged: %t, isAffected: %t)\n", c.ID, isInvolved, isChanged, isAffected)
			if isInvolved {
				err := runCommands(c, cmds, true, dryRun)
				if err != nil {
					log.Fatalf("Failed to execute specified commands: '%s', Error: %v\n", cmds, err)
				}
			}
		case "current_affected":
			c, err := getCurrentComponent()
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("Current Component '%s'!\n", c.String())
			isAffected := c.isAffected(&graph)
			fmt.Printf("%s is affected: %t\n", c.ID, isAffected)
			if isAffected {
				err := runCommands(c, cmds, true, dryRun)
				if err != nil {
					log.Fatalf("Failed to execute specified commands: '%s', Error: %v\n", cmds, err)
				}
			}
		case "current_changed":
			c, err := getCurrentComponent()
			if err != nil {
				log.Fatal(err)
			}
			isChanged := c.isChanged(&graph)
			fmt.Printf("%s is changed: %t\n", c.ID, isChanged)
			if isChanged {
				err := runCommands(c, cmds, true, dryRun)
				if err != nil {
					log.Fatalln(err)
					//log.Fatalf("Failed to execute specified commands: '%s', Error: %v\n", cmds, err)
				}
			}
		default:
			log.Fatalf("Unknown runMode '%s'", runMode)
		}
	},
}

func runCommands(component *Component, cmdLine string, printOutput bool, dryRun bool) error {

	if printOutput {
		fmt.Printf("runCommands for '%s (%s)' component (dry-run: %t, cmds: '%s')\n", component.ID, component.Path, dryRun, cmdLine)
	}

	args := strings.Split(cmdLine, " ")
	cmdName := args[0]
	args = args[1:]
	cmd := exec.Command(cmdName, args...)

	if !dryRun {
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		output := string(out[:])
		if printOutput {
			fmt.Printf("Output: \n%s\n", output)
		}
	} else {
		if printOutput {
			fmt.Printf("Output: \n%s\n", "** dry-run **")
		}
	}

	return nil
}

func init() {
	Main.AddCommand(run)
	run.Flags().StringVarP(&cmds, "commands", "c", "", "Commands to be run if current component is affected.")
	run.Flags().StringVarP(&runMode, "mode", "m", "current_involved", "Defines the mode of this run command. Supported values are: current_changed, current_affected, all_changed, all_affected. Default is current_involved.")
	run.Flags().BoolVar(&dryRun, "dry-run", false, "If dry-run is used, commands are only displayed, but not executed!")
}
