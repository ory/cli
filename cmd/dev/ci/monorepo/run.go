package monorepo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var cmds string
var runMode string
var dryRun bool

var run = &cobra.Command{
	Use:   "run",
	Short: "Runs the specified commands on changes",
	Long:  `Runs the specified commands if the current component (which is defined in the current work directory) is affected by any change in the repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var graph ComponentGraph
		if _, err := graph.getComponentGraph(rootDirectory); err != nil {
			return err
		}

		switch runMode {
		case "current_involved":
			c, err := getCurrentComponent()
			if err != nil {
				return err
			}
			//fmt.Printf("Current Component '%s'!\n", c.String())
			isAffected := c.isAffected(&graph)
			isChanged := c.isChanged(&graph)
			isInvolved := isChanged || isAffected
			fmt.Printf("%s is involved: %t (isChanged: %t, isAffected: %t)\n", c.ID, isInvolved, isChanged, isAffected)
			if isInvolved {
				if err := runCommands(c, cmds, dryRun); err != nil {
					return fmt.Errorf("failed to execute command '%s': %s", cmds, err)
				}
			}
		case "current_affected":
			c, err := getCurrentComponent()
			if err != nil {
				return err
			}

			isAffected := c.isAffected(&graph)
			fmt.Printf("%s is affected: %t\n", c.ID, isAffected)
			if isAffected {
				if err := runCommands(c, cmds, dryRun); err != nil {
					return fmt.Errorf("failed to execute command '%s': %s", cmds, err)
				}
			}
		case "current_changed":
			c, err := getCurrentComponent()
			if err != nil {
				return err
			}

			isChanged := c.isChanged(&graph)
			fmt.Printf("%s is changed: %t\n", c.ID, isChanged)
			if isChanged {
				if err := runCommands(c, cmds, dryRun); err != nil {
					return fmt.Errorf("failed to execute command '%s': %s", cmds, err)
				}
			}
		default:
			return fmt.Errorf("unknown runMode: %s", runMode)
		}

		return nil
	},
}

func runCommands(component *Component, cmdLine string, dryRun bool) error {
	fmt.Printf("runCommands for '%s (%s)' component (dry-run: %t, cmds: '%s')\n", component.ID, component.Path, dryRun, cmdLine)
	if dryRun {
		fmt.Print("Skipping execution because --dry-run was set.")
		return nil
	}

	args := strings.Split(cmdLine, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return nil
}

func init() {
	Main.AddCommand(run)
	run.Flags().StringVarP(&cmds, "commands", "c", "", "Commands to be run if current component is affected.")
	run.Flags().StringVarP(&runMode, "mode", "m", "current_involved", "Defines the mode of this run command. Supported values are: current_changed, current_affected, all_changed, all_affected. Default is current_involved.")
	run.Flags().BoolVar(&dryRun, "dry-run", false, "If dry-run is used, commands are only displayed, but not executed!")
}
