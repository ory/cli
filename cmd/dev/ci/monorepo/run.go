// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

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
var inverseMode bool

// ModeCurrentAffected will execute the specified commands if any of the current components dependecies have changed
const ModeCurrentAffected = "current_affected"

// ModeCurrentChanged will execute the specified commands if the current component has changed
const ModeCurrentChanged = "current_changed"

// ModeCurrentInvolved will execute the specified commands if the current component or any of its dependecies have changed
const ModeCurrentInvolved = "current_involved"

var run = &cobra.Command{
	Use:   "run",
	Short: "Runs the specified commands on changes",
	Long:  `Runs the specified commands if the current component (which is defined in the current work directory) is affected by any change in the repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var graph ComponentGraph
		if _, err := graph.getComponentGraph(rootDirectory); err != nil {
			return err
		}

		c, err := getCurrentComponent()
		if err != nil {
			return err
		}
		isAffected := c.isAffected(&graph)
		isChanged := c.isChanged(&graph)
		isInvolved := isChanged || isAffected

		return runWrapper(c, cmds, runMode, isAffected, isChanged, isInvolved, inverseMode)
	},
}

func runWrapper(c *Component, cmdLine string, mode string, affected bool, changed bool, involved bool, inverse bool) error {
	switch mode {
	case ModeCurrentInvolved:
		fmt.Printf("%s runCmd: %t (affected: %t, changed: %t, involved: %t, inverse: %t)\n", c.ID, involved != inverse, affected, changed, involved, inverse)
		if involved != inverse {
			if err := runCmd(c, cmdLine, dryRun); err != nil {
				return fmt.Errorf("failed to execute command '%s': %s", cmdLine, err)
			}
		}

	case ModeCurrentAffected:
		fmt.Printf("%s runCmd: %t (affected: %t, inverse: %t)\n", c.ID, affected != inverse, affected, inverse)
		if affected != inverse {
			if err := runCmd(c, cmdLine, dryRun); err != nil {
				return fmt.Errorf("failed to execute command '%s': %s", cmdLine, err)
			}
		}

	case ModeCurrentChanged:
		fmt.Printf("%s runCmd: %t (changed: %t, inverse: %t)\n", c.ID, changed != inverse, changed, inverse)
		if changed != inverse {
			if err := runCmd(c, cmdLine, dryRun); err != nil {
				return fmt.Errorf("failed to execute command '%s': %s", cmdLine, err)
			}
		}

	default:
		return fmt.Errorf("unknown runMode: %s", runMode)
	}
	return nil
}

func runCmd(component *Component, cmdLine string, dryRun bool) error {
	if debug {
		fmt.Printf("runCmd for '%s (%s)' component (dry-run: %t, cmds: '%s')\n", component.ID, component.Path, dryRun, cmdLine)
	}
	if dryRun {
		fmt.Print("Skipping execution because --dry-run was set.")
		return nil
	}

	args := strings.Split(cmdLine, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Main.AddCommand(run)
	run.Flags().StringVarP(&cmds, "commands", "c", "", "Commands to be run if current component is affected.")
	run.Flags().StringVarP(&runMode, "mode", "m", "current_involved", "Defines the mode of this run command. Supported values are: current_changed, current_affected, all_changed, all_affected. Default is current_involved.")
	run.Flags().BoolVar(&dryRun, "dry-run", false, "If dry-run is used, commands are only displayed, but not executed!")
	run.Flags().BoolVar(&inverseMode, "inverse", false, "If inverse is used, the specified commands will be executes if the current component is not affected/involved/changed (depending on mode)!")
}
