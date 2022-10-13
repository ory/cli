// Copyright © 2022 Ory Corp

package github

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ory/cli/cmd/pkg"
)

const tagPrefix = "refs/tags/"
const branchPrefix = "refs/heads/"

var env = &cobra.Command{
	Use:   "env",
	Short: "Sets up environment variables",
	Long: `Sets up environment variables. To load the environment variables use:

$ source $(ory dev ci github env)`,
	Run: func(cmd *cobra.Command, args []string) {
		if ref := os.Getenv("GITHUB_REF"); strings.HasPrefix(ref, tagPrefix) {
			// it's a tag
			fmt.Printf("export GIT_TAG=%s",
				strings.ReplaceAll(ref, tagPrefix, ""))
		} else if strings.HasPrefix(ref, branchPrefix) {
			fmt.Printf("export GIT_BRANCH=%s",
				strings.ReplaceAll(ref, branchPrefix, ""))
		} else {
			fmt.Printf("export GIT_BRANCH=%s",
				strings.TrimSpace(pkg.CommandGetOutput("git", "rev-parse", "--abbrev-ref", "HEAD")))
		}

		repo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
		if len(repo) != 2 {
			pkg.Fatalf("Malformed repository information in GITHUB_REPOSITORY: %s", os.Getenv("GITHUB_REPOSITORY"))
		}
		fmt.Printf("export GITHUB_ORG=%s\n", repo[0])
		fmt.Printf("export GITHUB_REPO=%s\n", repo[1])

		caser := cases.Title(language.AmericanEnglish)
		fmt.Printf("export SWAGGER_APP_NAME=%s_%s\n",
			caser.String(strings.ToLower(repo[0])),
			caser.String(strings.ToLower(repo[1])),
		)

		if ignorePkgs := strings.Split(os.Getenv("SWAGGER_SPEC_IGNORE_PKGS"), ","); len(ignorePkgs) > 0 {
			for k, p := range ignorePkgs {
				ignorePkgs[k] = "-x " + p
			}
			fmt.Printf(`export SWAGGER_SPEC_IGNORE_PKGS='%s'`, strings.Join(ignorePkgs, " "))
		}
	},
}

func init() {
	Main.AddCommand(env)
}
