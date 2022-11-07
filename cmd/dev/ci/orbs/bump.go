// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package orbs

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/pkg"
	"github.com/ory/x/flagx"
)

var orbs = []string{
	"ory/goreleaser",
	"ory/sdk",
	"ory/changelog",
	"ory/nancy",
	"ory/docs",
	"ory/prettier",
	"ory/golangci",
}

var orbLatestRegex = regexp.MustCompile(`(?im)^Latest:\s(.*)$`)

var bump = &cobra.Command{
	Use:   "bump <[.circleci/config.yml]>",
	Args:  cobra.RangeArgs(0, 1),
	Short: "Bump CircleCI Orb versions",
	Long: `Bumps ORY's CircleCI Orb versions to their newest version.

If no argument is supplied, this command uses the default ".circleci/config.yml" location.
`,
	Run: func(cmd *cobra.Command, args []string) {
		path := ".circleci/config.yml"
		if len(args) == 1 {
			path = args[0]
		}

		var wg sync.WaitGroup
		var lock sync.Mutex
		versions := map[string]string{}
		for _, id := range orbs {
			wg.Add(1)
			go getVersion(id, versions, &lock, &wg)
		}
		wg.Wait()

		config, err := os.ReadFile(path)
		pkg.Check(err)

		for k, r := range versions {
			replace := regexp.MustCompile(fmt.Sprintf("(?im)^(\\s\\s[^:]+:\\s)(%s@[0-9a-zA-Z\\.]+)$", k))
			config = []byte(replace.ReplaceAllString(string(config), "${1}"+r))
		}

		if flagx.MustGetBool(cmd, "write") {
			pkg.Check(os.WriteFile(path, config, 0666))
			fmt.Printf("Successfully wrote new orb versions to CircleCI config file: %s\n", path)
		} else {
			fmt.Println(string(config))
		}
	},
}

func getVersion(id string, versions map[string]string, l *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	info := pkg.CommandGetOutput("circleci", "--skip-update-check", "orb", "info", id)
	matches := orbLatestRegex.FindAllStringSubmatch(info, -1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		pkg.Fatalf(`Expected info to contain

	Latest: %s@a.b.c

but got:

%s`, id, info)
	}

	l.Lock()
	versions[id] = matches[0][1]
	l.Unlock()
}

func init() {
	Main.AddCommand(bump)
	bump.Flags().BoolP("write", "w", false, "Write output to CircleCI config file instead of stdout.")
}
